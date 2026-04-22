package client

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProjectBody return a json body
func ProjectBody(d *schema.ResourceData) models.ProjectsBodyPost {
	quota := d.Get("storage_quota").(int)

	body := models.ProjectsBodyPost{
		ProjectName:  d.Get("name").(string),
		RegistryID:   d.Get("registry_id").(int),
		StorageLimit: quota,
	}

	if quota != -1 {
		body.StorageLimit = quota * 1073741824
	}

	body.Metadata.AutoScan = strconv.FormatBool(d.Get("vulnerability_scanning").(bool))
	body.Metadata.Public = strconv.FormatBool(d.Get("public").(bool))

	security := d.Get("deployment_security").(string)
	if security != "" {
		body.Metadata.Severity = security
		body.Metadata.PreventVul = "true"
	} else {
		body.Metadata.Severity = "low"
		body.Metadata.PreventVul = "false"
	}

	body.Metadata.EnableContentTrust = strconv.FormatBool(d.Get("enable_content_trust").(bool))
	body.Metadata.EnableContentTrustCosign = strconv.FormatBool(d.Get("enable_content_trust_cosign").(bool))
	body.Metadata.AutoSbomGeneration = strconv.FormatBool(d.Get("auto_sbom_generation").(bool))
	body.Metadata.ProxySpeedKb = strconv.Itoa(d.Get("proxy_speed_kb").(int))

	cveAllowList := d.Get("cve_allowlist").([]interface{})
	log.Printf("[DEBUG] %v ", cveAllowList)
	if len(cveAllowList) > 0 {
		log.Printf("[DEBUG] %v ", expandCveAllowList(cveAllowList))
		body.CveAllowlist.Items = expandCveAllowList(cveAllowList)
		body.Metadata.ReuseSysCveAllowlist = "false"
	} else {
		body.Metadata.ReuseSysCveAllowlist = "true"
	}

	return body
}

// GetScannerByName lists all scanners and returns the one matching the given name.
func (client *Client) GetScannerByName(scanner string) (models.ScannerBody, error) {
	resp, _, _, err := client.SendRequest("GET", models.PathScanners, nil, 0)
	if err != nil {
		return models.ScannerBody{}, err
	}

	var scanners []models.ScannerBody
	err = json.Unmarshal([]byte(resp), &scanners)
	if err != nil {
		return models.ScannerBody{}, err
	}

	for _, v := range scanners {
		if strings.EqualFold(v.Name, scanner) {
			return v, nil
		}
	}

	return models.ScannerBody{}, fmt.Errorf("scanner %q not found", scanner)
}

// SetProjectScanner sets the vulnerability scanner for a specific project.
func (client *Client) SetProjectScanner(d *schema.ResourceData) error {
	scanner := d.Get("vulnerability_scanner").(string)
	if scanner == "" {
		return nil
	}

	scannerData, err := client.GetScannerByName(scanner)
	if err != nil {
		return err
	}

	body := models.ProjectScannerBody{
		UUID: scannerData.UUID,
	}

	_, _, _, err = client.SendRequest("PUT", d.Id()+"/scanner", body, 200)
	return err
}

// GetProjectScanner returns the name of the scanner assigned to a project.
func (client *Client) GetProjectScanner(projectPath string) (string, error) {
	resp, _, respCode, err := client.SendRequest("GET", projectPath+"/scanner", nil, 200)
	if err != nil {
		if respCode == 404 {
			return "", nil
		}
		return "", err
	}

	var scannerData models.ScannerBody
	err = json.Unmarshal([]byte(resp), &scannerData)
	if err != nil {
		return "", err
	}

	return scannerData.Name, nil
}

func expandCveAllowList(cveAllowlist []interface{}) models.CveAllowlistItems {
	allowlist := models.CveAllowlistItems{}

	for _, data := range cveAllowlist {
		item := models.CveAllowlistItem{
			CveID: data.(string),
		}
		allowlist = append(allowlist, item)
	}

	return allowlist
}

func (client *Client) UpdateStorageQuota(d *schema.ResourceData) (err error) {
	resp, _, _, err := client.SendRequest("GET", models.PathConfig, nil, 200)
	if err != nil {
		return err
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}

	if jsonData.QuotaPerProjectEnable.Value == true {
		if d.HasChange("storage_quota") {
			projectID := strings.Replace(d.Id(), "/projects/", "", -1)
			page := 1

			for {
				quotasPath := fmt.Sprintf("/quotas/?page=%d&page_size=100", page)

				resp, _, _, err := client.SendRequest("GET", quotasPath, nil, 200)
				if err != nil {
					return err
				}

				var quotaResponse []models.QuotaResponse
				if err := json.Unmarshal([]byte(resp), &quotaResponse); err != nil {
					return err
				}

				if len(quotaResponse) == 0 {
					return nil
				}

				for _, q := range quotaResponse {
					pid := strconv.Itoa(q.Ref.ID)
					if pid == projectID {
						quotaID := "/quotas/" + strconv.Itoa(q.ID)
						storage := d.Get("storage_quota").(int)
						if storage > 0 {
							storage *= 1073741824 // GB
						}
						quota := models.Hard{
							Storage: int64(storage),
						}
						body := models.StorageQuota{quota}

						_, _, _, err = client.SendRequest("PUT", quotaID, body, 200)
						if err != nil {
							return err
						}
						break
					}
				}
				page++
			}
		}

	}

	return nil
}
