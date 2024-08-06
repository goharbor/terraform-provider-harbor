package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigSystem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_creation_restriction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "adminonly",
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"robot_token_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"robot_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "robot$",
			},
			"scanner_skip_update_pulltime": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"storage_per_project": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"audit_log_forward_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"skip_audit_log_database": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"banner_message": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"closable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"message": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "info",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								allowedValues := []string{"success", "info", "warning", "danger"}
								for _, av := range allowedValues {
									if v == av {
										return
									}
								}
								errs = append(errs, fmt.Errorf("%q must be one of [%s], got %s", key, strings.Join(allowedValues, ", "), v))
								return
							},
						},
						"from_date": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"to_date": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
		Create: resourceConfigSystemCreate,
		Read:   resourceConfigSystemRead,
		Update: resourceConfigSystemCreate,
		Delete: resourceConfigSystemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				resourceConfigSystemRead(d, m)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceConfigSystemCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.GetConfigSystem(d)

	_, _, _, err := apiClient.SendRequest("PUT", models.PathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathConfig, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("error getting system configuration %s", err)
	}

	var jsonData models.ConfigBodyResponse
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("error getting system configuration %s", err)
	}
	storage := jsonData.StoragePerProject.Value
	if storage > 0 {
		storage /= 1073741824 // Byte to GB
	}

	// set the data for banner message
	if jsonData.BannerMessage.Value != "" {
		bannerMessage := jsonData.BannerMessage.Value
		// bannerMessage is a json string and we need to parse it to get the values
		var bannerMessageData models.BannerMessage
		err = json.Unmarshal([]byte(bannerMessage), &bannerMessageData)
		if err != nil {
			return fmt.Errorf("error getting system configuration %s", err)
		}
		bannerMessageList := make([]interface{}, 1)
		bannerMessageList[0] = map[string]interface{}{
			"closable":  bannerMessageData.Closable,
			"message":   bannerMessageData.Message,
			"type":      bannerMessageData.Type,
			"from_date": bannerMessageData.FromDate,
			"to_date":   bannerMessageData.ToDate,
		}
		d.Set("banner_message", bannerMessageList)
	}

	d.Set("project_creation_restriction", jsonData.ProjectCreationRestriction.Value)
	d.Set("read_only", jsonData.ReadOnly.Value)
	d.Set("robot_token_expiration", jsonData.RobotTokenDuration.Value)
	d.Set("robot_name_prefix", jsonData.RobotNamePrefix.Value)
	d.Set("scanner_skip_update_pulltime", jsonData.ScannerSkipUpdatePulltime.Value)
	d.Set("storage_per_project", storage)
	d.Set("audit_log_forward_endpoint", jsonData.AuditLogForwardEndpoint.Value)
	d.Set("skip_audit_log_database", jsonData.SkipAuditLogDatabase.Value)

	d.SetId("configuration/system")
	return nil
}

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
