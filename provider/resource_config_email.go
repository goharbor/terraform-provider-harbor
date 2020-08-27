package provider

import (
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type email struct {
	EmailHost     string `json:"email_host"`
	EmailPort     string `json:"email_port"`
	EmailUsername string `json:"email_username,omitempty"`
	EmailPassword string `json:"email_password,omitempty"`
	EmailFrom     string `json:"email_from"`
	EmailSsl      string `json:"email_ssl,omitempty"`
	// EmailVerifyCert string `json:"email_verify_cert,omitempty"`
}

func resourceConfigEmail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email_from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_ssl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
			// "email_verify_cert": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
		},
		Create: resourceConfigEmailCreate,
		Read:   resourceConfigEmailRead,
		Update: resourceConfigEmailUpdate,
		Delete: resourceConfigEmailDelete,
	}
}

func resourceConfigEmailCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := email{
		EmailHost:     d.Get("email_host").(string),
		EmailPort:     d.Get("email_port").(string),
		EmailUsername: d.Get("email_username").(string),
		EmailPassword: d.Get("email_password").(string),
		EmailFrom:     d.Get("email_from").(string),
		EmailSsl:      d.Get("email_ssl").(string),
		// EmailVerifyCert: d.Get("email_verify_cert").(string),
	}

	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	d.SetId(randomString(15))
	// return resourceConfigEmailRead(d, m)
	return nil
}

func resourceConfigEmailRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConfigEmailUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := email{
		EmailHost:     d.Get("email_host").(string),
		EmailPort:     d.Get("email_port").(string),
		EmailUsername: d.Get("email_username").(string),
		EmailPassword: d.Get("email_password").(string),
		EmailFrom:     d.Get("email_from").(string),
		EmailSsl:      d.Get("email_ssl").(string),
		// EmailVerifyCert: d.Get("email_verify_cert").(string),
	}

	_, err := apiClient.SendRequest("PUT", pathConfig, body, 200)
	if err != nil {
		return err
	}

	return resourceConfigEmailRead(d, m)
}

func resourceConfigEmailDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
