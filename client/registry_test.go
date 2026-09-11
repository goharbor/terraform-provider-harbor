package client

import (
    "testing"
)

func TestGetRegistryAPIType(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"provider_name", "azure", "azure-acr", false},
        {"api_type", "azure-acr", "azure-acr", false},
        {"quay-io alias", "quay-io", "quay", false},
        {"unknown", "foo", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GetRegistryAPIType(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetRegistryAPIType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("GetRegistryAPIType() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestGetRegistryType(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"api_type_to_provider", "azure-acr", "azure"},
        {"quay_api_type", "quay", "quay"},
        {"unmapped_type_passes_through", "custom-registry", "custom-registry"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, _ := GetRegistryType(tt.input)
            if got != tt.want {
                t.Errorf("GetRegistryType() = %v, want %v", got, tt.want)
            }
        })
    }
}
