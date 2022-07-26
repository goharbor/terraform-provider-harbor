package client

import (
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRegistryBody(d *schema.ResourceData) models.RegistryBody {
	regType, _ := GetRegistryType(d.Get("provider_name").(string))

	body := models.RegistryBody{
		Description: d.Get("description").(string),
		Insecure:    d.Get("insecure").(bool),
		Name:        d.Get("name").(string),
		Type:        regType,
		URL:         d.Get("endpoint_url").(string),
	}

	body.Credential.AccessKey = d.Get("access_id").(string)
	body.Credential.AccessSecret = d.Get("access_secret").(string)
	body.Credential.Type = "basic"

	return body
}

func GetRegistryType(regType string) (regName string, err error) {

	registryType := map[string]string{
		"alibaba":         "ali-acr",
		"aws":             "aws-ecr",
		"azure":           "azure-acr",
		"docker-hub":      "docker-hub",
		"docker-registry": "docker-registry",
		"gitlab":          "gitlab",
		"google":          "google-gcr",
		"harbor":          "harbor",
		"helm":            "helm-hub",
		"huawei":          "huawei-SWR",
		"jfrog":           "jfrog-artifactory",
		"quay":            "quay-io",
		// for reverse lookup
		"ali-acr":           "alibaba",
		"aws-ecr":           "aws",
		"azure-acr":         "azure",
		"google-gcr":        "google",
		"helm-hub":          "helm",
		"huawei-SWR":        "huawei",
		"jfrog-artifactory": "jfrog",
		"quay-io":           "quay",
	}

	return registryType[regType], nil
}
