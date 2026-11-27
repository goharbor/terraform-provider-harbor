package provider

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("access_secret"), cty.GetAttrPath("access_secret_wo")),
		},
		Schema: map[string]*schema.Schema{
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Normalize both old (state) and new (config) values to their API type
					// to prevent perpetual diffs when equivalent types are used (e.g., "azure" vs "azure-acr").
					oldNorm, _ := client.GetRegistryAPIType(old)
					newNorm, _ := client.GetRegistryAPIType(new)
					return oldNorm == newNorm
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok || v == "" {
						return
					}
					if _, err := client.GetRegistryAPIType(v); err != nil {
						errs = append(errs, err)
					}
					return
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"access_secret": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				Default:   "",
				ConflictsWith: []string{
					"access_secret_wo",
					"access_secret_wo_version",
				},
			},
			"access_secret_wo": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				WriteOnly: true,
				RequiredWith: []string{
					"access_secret_wo_version",
				},
				ConflictsWith: []string{
					"access_secret",
				},
			},
			"access_secret_wo_version": {
				Type:     schema.TypeInt,
				Optional: true,
				RequiredWith: []string{
					"access_secret_wo",
				},
				ConflictsWith: []string{
					"access_secret",
				},
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRegistryCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body, err := client.GetRegistryBody(d)
	if err != nil {
		return err
	}

	accessSecretWriteOnly, err := getWriteOnlyString(d, "access_secret_wo")
	if err != nil {
		return err
	}
	if accessSecretWriteOnly != "" {
		body.Credential.AccessSecret = accessSecretWriteOnly
	}

	_, headers, _, err := apiClient.SendRequest("POST", models.PathRegistries, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceRegistryRead(d, m)
}

func resourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	var jsonData models.RegistryBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("resource not found %s", d.Id())
	}

	registryName, _ := client.GetRegistryType(jsonData.Type)

	d.Set("name", jsonData.Name)
	d.Set("description", jsonData.Description)
	d.Set("endpoint_url", jsonData.URL)
	d.Set("access_id", jsonData.Credential.AccessKey)
	d.Set("insecure", jsonData.Insecure)
	d.Set("ca_certificate", jsonData.CACertificate)
	d.Set("status", jsonData.Status)
	d.Set("registry_id", jsonData.ID)
	d.Set("provider_name", registryName)

	return nil
}

func resourceRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	oldAccessSecret, _ := d.GetChange("access_secret")
	oldAccessSecretWOVersion, newAccessSecretWOVersion := d.GetChange("access_secret_wo_version")
	accessSecretWriteOnly, err := getWriteOnlyString(d, "access_secret_wo")
	if err != nil {
		return err
	}

	if d.HasChange("access_secret_wo_version") && newAccessSecretWOVersion.(int) != 0 &&
		accessSecretWriteOnly == "" && !d.HasChange("access_secret") {
		_ = d.Set("access_secret_wo_version", oldAccessSecretWOVersion)
		return fmt.Errorf("access_secret_wo must be configured when access_secret_wo_version changes")
	}

	body := client.GetRegistryUpdateBody(d)
	// Harbor overwrites the stored secret with whatever is sent, so always send the write-only value when configured.
	if accessSecretWriteOnly != "" {
		body.AccessSecret = accessSecretWriteOnly
	}

	_, _, _, err = apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		if d.HasChange("access_secret") {
			_ = d.Set("access_secret", oldAccessSecret)
		}
		if d.HasChange("access_secret_wo_version") {
			_ = d.Set("access_secret_wo_version", oldAccessSecretWOVersion)
		}
		return err
	}

	return resourceRegistryRead(d, m)
}

func resourceRegistryDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if respCode != 404 && err != nil { // We can't delete something that doesn't exist. Hence the 404-check
		return err
	} else if err != nil {
		return err
	}
	return nil
}
