package provider

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getWriteOnlyString(d *schema.ResourceData, key string) (string, error) {
	rawValue, diags := d.GetRawConfigAt(cty.GetAttrPath(key))

	return parseWriteOnlyString(rawValue, diags, key)
}

func parseWriteOnlyString(rawValue cty.Value, diags diag.Diagnostics, key string) (string, error) {
	if diags.HasError() {
		return "", fmt.Errorf("error retrieving write-only argument %q: %v", key, diags)
	}

	if rawValue.IsNull() {
		return "", nil
	}

	if !rawValue.Type().Equals(cty.String) {
		return "", fmt.Errorf("error retrieving write-only argument %q: value must be a string", key)
	}

	return rawValue.AsString(), nil
}
