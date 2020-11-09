package provider

import (
	"encoding/json"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGarbageCollect() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_untagged": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourceGarbageCollectCreate,
		Read:   resourceGarbageCollectRead,
		Update: resourceGarbageCollectCreate,
		Delete: resourceGarbageCollectDelete,
	}
}

func resourceGarbageCollectCreate(d *schema.ResourceData, m interface{}) error {
	var method string = "POST"
	var code int = 201
	if resourceGarbageCollectExist(m) {
		method = "PUT"
		code = 200
	}
	apiClient := m.(*client.Client)

	body := client.GCBodyPost(d)

	_, _, err := apiClient.SendRequest(method, models.PathGC, body, code)
	if err != nil {
		return err
	}

	return resourceGarbageCollectRead(d, m)
}

func resourceGarbageCollectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", models.PathGC, nil, 200)
	if err != nil {
		fmt.Println(err)
	}

	var jsonData models.GCBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}
	d.Set("type", jsonData.Schedule.Type)
	if jsonData.Schedule.Type == "Custom" {
		d.Set("cron", jsonData.Schedule.Cron)
	} else {
		d.Set("cron", "")
	}
	d.Set("delete_untagged", jsonData.JobParameters)
	d.SetId("system/gc/schedule")
	return nil
}

func resourceGarbageCollectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	body := client.GCBodyDelete()
	_, _, err := apiClient.SendRequest("POST", models.PathGC, body, 201)
	if err != nil {
		fmt.Println(err)
	}

	return resourceGarbageCollectRead(d, m)
}

func resourceGarbageCollectExist(m interface{}) bool {
	apiClient := m.(*client.Client)
	resp, _, err := apiClient.SendRequest("GET", models.PathGC, nil, 200)
	if err != nil {
		fmt.Println(err)
	}

	var jsonData models.GCBodyResponses
	err = json.Unmarshal([]byte(resp), &jsonData)
	if jsonData.JobParameters == "" {
		return false
	} else {
		return true
	}
}