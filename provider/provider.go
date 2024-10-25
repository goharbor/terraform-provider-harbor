package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_URL", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_PASSWORD", ""),
			},
			"bearer_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_BEARER_TOKEN", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_IGNORE_CERT", ""),
			},
			"api_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"harbor_config_auth":            resourceConfigAuth(),
			"harbor_config_system":          resourceConfigSystem(),
			"harbor_config_security":        resourceConfigSecurity(),
			"harbor_project":                resourceProject(),
			"harbor_project_member_group":   resourceMembersGroup(),
			"harbor_project_member_user":    resourceMembersUser(),
			"harbor_project_webhook":        resourceProjectWebhook(),
			"harbor_tasks":                  resourceTasks(),
			"harbor_interrogation_services": resourceVuln(),
			"harbor_robot_account":          resourceRobotAccount(),
			"harbor_user":                   resourceUser(),
			"harbor_group":                  resourceGroup(),
			"harbor_registry":               resourceRegistry(),
			"harbor_replication":            resourceReplication(),
			"harbor_retention_policy":       resourceRetention(),
			"harbor_garbage_collection":     resourceGC(),
			"harbor_purge_audit_log":        resourcePurgeAudit(),
			"harbor_label":                  resourceLabel(),
			"harbor_preheat_instance":       resourcePreheatInstance(),
			"harbor_immutable_tag_rule":     resourceImmutableTagRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"harbor_project":               dataProject(),
			"harbor_projects":              dataProjects(),
			"harbor_registry":              dataRegistry(),
			"harbor_groups":                dataGroups(),
			"harbor_robot_accounts":        dataRobotAccounts(),
			"harbor_project_member_groups": dataProjectMemberGroups(),
			"harbor_project_member_users":  dataProjectMemberUsers(),
			"harbor_users":                 dataUsers(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var apiPath string

	//url := d.Get("url").(string)
	url := os.Getenv("HARBOR_URL")
	if d.Get("url").(string) != "" {
		url = d.Get("url").(string)
	}
	if url == "" {
		return nil, fmt.Errorf("url is required and must be provided in the provider config or the HARBOR_URL environment variable")
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	bearerToken := d.Get("bearer_token").(string)
	insecure := d.Get("insecure").(bool)
	apiVersion := d.Get("api_version").(int)

	if strings.HasSuffix(url, "/") {
		url = strings.Trim(url, "/")
	}

	if apiVersion == 1 {
		apiPath = "/api"
	} else if apiVersion == 2 {
		apiPath = "/api/v2.0"
	}

	return client.NewClient(url+apiPath, username, password, bearerToken, insecure), nil
}
