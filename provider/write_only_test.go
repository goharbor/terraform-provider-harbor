package provider

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func TestParseWriteOnlyString(t *testing.T) {
	t.Run("returns configured write-only value", func(t *testing.T) {
		value, err := parseWriteOnlyString(cty.StringVal("my-secret"), nil, "secret_wo")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if value != "my-secret" {
			t.Fatalf("unexpected value: got %q", value)
		}
	})

	t.Run("returns empty string when value is not configured", func(t *testing.T) {
		value, err := parseWriteOnlyString(cty.NullVal(cty.String), nil, "secret_wo")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if value != "" {
			t.Fatalf("unexpected value: got %q", value)
		}
	})

	t.Run("returns error for non-string write-only value", func(t *testing.T) {
		_, err := parseWriteOnlyString(cty.NumberIntVal(1), nil, "secret_wo")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("returns error when diagnostics contain error", func(t *testing.T) {
		diags := diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid config path",
				Detail:   "Cannot find config value for given path.",
			},
		}

		_, err := parseWriteOnlyString(cty.DynamicVal, diags, "password_wo")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
