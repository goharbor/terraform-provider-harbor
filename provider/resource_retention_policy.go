package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRentention() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 15,
				MinItems: 1,
				Elem: &schema.Resource{
					// Schema: retentionPolicyRuleFields(),
					Schema: map[string]*schema.Schema{
						"n_days_since_last_pull": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  nil,
							// ConflictsWith: []string{"n_days_since_last_push"},
						},
						"n_days_since_last_push": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  nil,
							// ConflictsWith: []string{"n_days_since_last_pull"},
						},
						"most_recently_pulled": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"most_recently_pushed": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"always": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"repo_matching": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{".repo_excluding"},
						},
						"repo_excluding": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{".repo_matching"},
						},
						"tag_matching": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{".tag_excluding"},
						},
						"tag_excluding": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{".tag_matching"},
						},
						"untagged_artifacts": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
		Create: resourceRententionCreate,
		Read:   resourceRententionRead,
		Update: resourceRententionUpdate,
		Delete: resourceRententionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRententionCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GetRententionBody(d)

	_, headers, err := apiClient.SendRequest("POST", models.PathRetentions, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	d.SetId(id)
	return resourceRententionRead(d, m)
}

func resourceRententionRead(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	// if err != nil {
	// 	return err
	// }

	// var jsonData models.Labels
	// err = json.Unmarshal([]byte(resp), &jsonData)
	// if err != nil {
	// 	return fmt.Errorf("Resource not found %s", d.Id())
	// }

	// d.Set("name", jsonData.Name)
	// d.Set("description", jsonData.Description)
	// d.Set("color", jsonData.Color)
	// d.Set("scope", jsonData.Scope)

	return nil
}

func resourceRententionUpdate(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)
	// body := client.LabelsBody(d)

	// _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	// if err != nil {
	// 	return err
	// }

	return resourceRententionRead(d, m)
}

func resourceRententionDelete(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// _, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	// if err != nil {
	// 	return err
	// }
	return nil
}
