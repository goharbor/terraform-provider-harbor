package client

import (
	"fmt"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRegistryBody(d *schema.ResourceData) models.RegistryBody {
	regType, _ := GetReigstryType(d.Get("provider_name").(string))

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

func GetReigstryType(regType string) (regName string, err error) {
	switch regType {
	case "alibaba":
		return "ali-acr", nil
	case "aws":
		return "aws-ecr", nil
	case "azure":
		return "azure-acr", nil
	case "docker-hub":
		return "docker-hub", nil
	case "docker-registry":
		return "docker-registry", nil
	case "gitlab":
		return "gitlab", nil
	case "google":
		return "google-gcr", nil
	case "harbor":
		return "harbor", nil
	case "helm":
		return "helm-hub", nil
	case "huawei":
		return "huawei-SWR", nil
	case "jfrog":
		return "jfrog-artifactory", nil
	case "quay":
		return "quay-io", nil

	// for reverse lookup
	case "ali-acr":
		return "alibaba", nil
	case "aws-ecr":
		return "aws", nil
	case "azure-acr":
		return "azure", nil
	case "google-gcr":
		return "google", nil
	case "helm-hub":
		return "helm", nil
	case "huawei-SWR":
		return "huawei", nil
	case "jfrog-artifactory":
		return "jfrog", nil
	case "quay-io":
		return "quay", nil

	default:
		return "", fmt.Errorf("Unable to find type for %s", regType)
	}

}
