package client

import (
	"fmt"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRegistryBody(d *schema.ResourceData) models.RegistryBody {
	regType, _ := GetRegistryType(d.Get("provider_name").(string))

	body := models.RegistryBody{
		Description:   d.Get("description").(string),
		Insecure:      d.Get("insecure").(bool),
		Name:          d.Get("name").(string),
		Type:          regType,
		URL:           d.Get("endpoint_url").(string),
		CACertificate: d.Get("ca_certificate").(string),
	}

	body.Credential.AccessKey = d.Get("access_id").(string)
	body.Credential.AccessSecret = d.Get("access_secret").(string)
	body.Credential.Type = "basic"

	return body
}

func GetRegistryUpdateBody(d *schema.ResourceData) models.RegistryUpdateBody {

	body := models.RegistryUpdateBody{
		AccessKey:     d.Get("access_id").(string),
		AccessSecret:  d.Get("access_secret").(string),
		Description:   d.Get("description").(string),
		Insecure:      d.Get("insecure").(bool),
		Name:          d.Get("name").(string),
		URL:           d.Get("endpoint_url").(string),
		CACertificate: d.Get("ca_certificate").(string),
	}

	return body
}

func GetRegistryType(regType string) (regName string, err error) {
    // Map of Terraform provider names to Harbor API types
    providerToAPI := map[string]string{
        "alibaba":         "ali-acr",
        "artifact-hub":    "artifact-hub",
        "aws":             "aws-ecr",
        "azure":           "azure-acr",
        "docker-hub":      "docker-hub",
        "docker-registry": "docker-registry",
        "gitlab":          "gitlab",
        "github":          "github-ghcr",
        "google":          "google-gcr",
        "harbor":          "harbor",
        "helm":            "helm-hub",
        "huawei":          "huawei-SWR",
        "jfrog":           "jfrog-artifactory",
        "quay":            "quay",
    }

    // If the user passed a provider name (e.g., "azure"), return the API type (e.g., "azure-acr")
    if apiType, ok := providerToAPI[regType]; ok {
        return apiType, nil
    }

    // If the user passed an API type directly (e.g., "azure-acr"), validate it and return as-is
    for _, apiType := range providerToAPI {
        if regType == apiType {
            return regType, nil
        }
    }

    // If it's neither a provider name nor a valid API type, return an error
    return "", fmt.Errorf("unknown registry type: %s. Please use a valid provider name (e.g., 'azure') or API type (e.g., 'azure-acr')", regType)
}
