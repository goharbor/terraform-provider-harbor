package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePurgeAudit() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audit_retention_hour": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"include_operations": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIncludeOperations,
			},
		},
		Create: resourcePurgeAuditCreate,
		Read:   resourcePurgeAuditRead,
		Update: resourcePurgeAuditUpdate,
		Delete: resourcePurgeAuditDelete,
	}
}

func resourcePurgeAuditCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	err := apiClient.SetSchedule(d, "purgeaudit")
	if err != nil {
		return err
	}
	d.SetId(models.PathPurgeAudit)
	return resourcePurgeAuditRead(d, m)
}

func resourcePurgeAuditRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", models.PathPurgeAudit, nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return fmt.Errorf("resource not found %s", d.Id())
	}
	if len(resp) == 0 {
		d.SetId("")
		return nil
	}

	var jsonData models.SystemBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return err
	}
	jobParameters := jsonData.JobParameters

	var jsonJobParameters models.JobParameters
	err = json.Unmarshal([]byte(jobParameters), &jsonJobParameters)
	if err != nil {
		fmt.Println(err)
	}

	if jsonData.Schedule.Type == "Custom" {
		d.Set("schedule", jsonData.Schedule.Cron)
	} else {
		d.Set("schedule", jsonData.Schedule.Type)
	}
	d.Set("audit_retention_hour", jsonJobParameters.AuditRetentionHour)
	d.Set("include_operations", jsonJobParameters.IncludeOperations)
	return nil
}

func resourcePurgeAuditUpdate(d *schema.ResourceData, m interface{}) error {
	return resourcePurgeAuditCreate(d, m)
}

func resourcePurgeAuditDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	d.Set("schedule", "")
	err := apiClient.SetSchedule(d, "purgeaudit")
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func validateIncludeOperations(v interface{}, k string) (warns []string, errs []error) {
	includeOperations := v.(string)
	validValues := []string{"create", "pull", "delete"}

	ops := strings.Split(includeOperations, ",")
	for _, op := range ops {
		op = strings.TrimSpace(op)
		if !containsString(validValues, op) {
			errs = append(errs, fmt.Errorf("Invalid value %q in %q. Valid values are: create, pull, delete", op, k))
		}
	}

	return warns, errs
}

func containsString(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
