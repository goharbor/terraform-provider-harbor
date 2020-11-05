package provider

import (
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
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
			"harbor_config_email":           resourceConfigEmail(),
			"harbor_config_system":          resourceConfigSystem(),
			"harbor_project":                resourceProject(),
			"harbor_project_member":         resourceMembers(),
			"harbor_project_member_group":   resourceMembersGroup(),
			"harbor_project_member_user":    resourceMembersUser(),
			"harbor_tasks":                  resourceTasks(),
			"harbor_interrogation_services": resourceVuln(),
			"harbor_robot_account":          resourceRobotAccount(),
			"harbor_user":                   resourceUser(),
			"harbor_registry":               resourceRegistry(),
			"harbor_replication":            resourceReplication(),
			"resource_garbage_collection":   resourceGC(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"harbor_project":  dataProject(),
			"harbor_registry": dataRegistry(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var apiPath string

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
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

	return client.NewClient(url+apiPath, username, password, insecure), nil
}
