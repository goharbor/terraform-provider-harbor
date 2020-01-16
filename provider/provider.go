package provider

import (
	"bitbucket.org/bestsellerit/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"harbor_config_auth":    resourceConfigAuth(),
			"harbor_config_email":   resourceConfigEmail(),
			"harbor_config_system":  resourceConfigSystem(),
			"harbor_project":        resourceProject(),
			"harbor_project_member": resourceMembers(),
			"harbor_tasks":          resourceTasks(),
			"harbor_robot_account":  resourceRobotAccount(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)

	return client.NewClient(url, username, password, insecure), nil
}
