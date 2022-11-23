package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembersGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"group_id", "group_name"},
			},
			"ldap_group_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validRoles := []string{"projectadmin", "developer", "guest", "maintainer", "limitedguest"}
					for _, r := range validRoles {
						if v == r {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validRoles, v))
					return
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "ldap" && v != "internal" && v != "oidc" {
						errs = append(errs, fmt.Errorf("%q must be either ldap, internal or oidc, got: %s", key, v))
					}
					return
				},
			},
		},
		CreateContext: resourceMembersGroupCreate,
		ReadContext:   resourceMembersGroupRead,
		UpdateContext: resourceMembersGroupUpdate,
		DeleteContext: resourceMembersGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMembersGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))
	path := projectid + "/members"

	body := client.ProjectMembersGroupBody(d)

	_, headers, _, err := apiClient.SendRequest(ctx, "POST", path, body, 201)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.GetID(headers)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceMembersGroupRead(ctx, d, m)
}

func resourceMembersGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest(ctx, "GET", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Project member group %q was not found - removing from state!", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("making read request on project member group %s : %+v", d.Id(), err)
	}

	var jsonData models.ProjectMembersBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("role", client.RoleTypeNumber(jsonData.RoleID))
	d.Set("project_id", checkProjectid(strconv.Itoa(jsonData.ProjectID)))
	d.Set("group_name", jsonData.EntityName)
	return nil
}

func resourceMembersGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	body := client.ProjectMembersGroupBody(d)
	_, _, _, err := apiClient.SendRequest(ctx, "PUT", d.Id(), body, 200)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMembersGroupRead(ctx, d, m)
}

func resourceMembersGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest(ctx, "DELETE", d.Id(), nil, 200)
	if err != nil {
		if respCode == 404 {
			log.Printf("[DEBUG] Project member group %q was not found - already deleted!", d.Id())
			return nil
		}
		return diag.Errorf("making delete request on project member group %s : %+v", d.Id(), err)
	}
	return nil
}
