package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceRegistrySchemaWriteOnlyAccessSecret(t *testing.T) {
	r := resourceRegistry()

	if err := r.InternalValidate(nil, true); err != nil {
		t.Fatalf("unexpected schema validation error: %s", err)
	}

	accessSecretWO := r.Schema["access_secret_wo"]
	if !accessSecretWO.WriteOnly {
		t.Fatal("access_secret_wo must be write-only")
	}
	if !accessSecretWO.Sensitive {
		t.Fatal("access_secret_wo must be sensitive")
	}
	if r.Schema["access_secret_wo_version"].WriteOnly {
		t.Fatal("access_secret_wo_version must be stored in state to trigger updates")
	}
}

func TestResourceRegistryValidateAccessSecretArguments(t *testing.T) {
	base := func(extra map[string]interface{}) map[string]interface{} {
		cfg := map[string]interface{}{
			"provider_name": "docker-hub",
			"name":          "docker-hub",
			"endpoint_url":  "https://hub.docker.com",
			"access_id":     "user",
		}
		for k, v := range extra {
			cfg[k] = v
		}
		return cfg
	}

	cases := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{"no secret", base(nil), false},
		{"access_secret only", base(map[string]interface{}{"access_secret": "s3cr3t"}), false},
		{"write-only secret with version", base(map[string]interface{}{"access_secret_wo": "s3cr3t", "access_secret_wo_version": 1}), false},
		{"write-only secret without version", base(map[string]interface{}{"access_secret_wo": "s3cr3t"}), true},
		{"version without write-only secret", base(map[string]interface{}{"access_secret_wo_version": 1}), true},
		{"access_secret with write-only secret", base(map[string]interface{}{"access_secret": "s3cr3t", "access_secret_wo": "s3cr3t", "access_secret_wo_version": 1}), true},
		{"access_secret with version", base(map[string]interface{}{"access_secret": "s3cr3t", "access_secret_wo_version": 1}), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diags := resourceRegistry().Validate(terraform.NewResourceConfigRaw(tc.config))
			if diags.HasError() != tc.wantErr {
				t.Fatalf("Validate() errors = %v, wantErr %v", diags, tc.wantErr)
			}
		})
	}
}
