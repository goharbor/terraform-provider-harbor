package client

import (
    "fmt"

    "github.com/goharbor/terraform-provider-harbor/models"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerToAPI is the single source of truth for mapping Terraform provider names to Harbor API types.
var providerToAPI = map[string]string{
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
    "quay-io":         "quay",
}

// apiToProvider is the reverse map of providerToAPI, generated at init.
var apiToProvider map[string]string

func init() {
    apiToProvider = make(map[string]string)
    for provider, api := range providerToAPI {
        apiToProvider[api] = provider
    }
}

// GetRegistryBody populates the registry body from schema.
func GetRegistryBody(d *schema.ResourceData) (models.RegistryBody, error) {
    regType, err := GetRegistryAPIType(d.Get("provider_name").(string))
    if err != nil {
        return models.RegistryBody{}, err
    }

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

    return body, nil
}

// GetRegistryUpdateBody populates the update body from schema.
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

// GetRegistryType is used during Read to translate the Harbor API type back to the provider name for state.
func GetRegistryType(regType string) (regName string, err error) {
    if v, ok := apiToProvider[regType]; ok {
        return v, nil
    }
    // If it's not in the map, it might be a custom or unmapped type. 
    // Returning the original string prevents state drift/diffs.
    return regType, nil
}

// GetRegistryAPIType normalizes user input to the Harbor API type.
// It accepts both provider names (e.g., "azure") and API types (e.g., "azure-acr").
func GetRegistryAPIType(regType string) (string, error) {
    if apiType, ok := providerToAPI[regType]; ok {
        return apiType, nil
    }

    // Check if the user passed a valid API type directly
    if _, ok := apiToProvider[regType]; ok {
        return regType, nil
    }

    return "", fmt.Errorf("unknown registry type: %q. Please use a valid provider name (e.g., 'azure') or API type (e.g., 'azure-acr')", regType)
}
