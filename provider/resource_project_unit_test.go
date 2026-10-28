package provider

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
)

func TestProjectBodyCveAllowlistExpiresAt(t *testing.T) {
	cases := []struct {
		name        string
		raw         map[string]interface{}
		wantExpires int
		wantReuse   string
	}{
		{
			name: "expiry is sent with the allowlist",
			raw: map[string]interface{}{
				"name":                     "p",
				"cve_allowlist":            []interface{}{"CVE-2024-0001"},
				"cve_allowlist_expires_at": 1701167767,
			},
			wantExpires: 1701167767,
			wantReuse:   "false",
		},
		{
			name: "no expiry by default",
			raw: map[string]interface{}{
				"name":          "p",
				"cve_allowlist": []interface{}{"CVE-2024-0001"},
			},
			wantExpires: 0,
			wantReuse:   "false",
		},
		{
			name: "expiry ignored without allowlist",
			raw: map[string]interface{}{
				"name":                     "p",
				"cve_allowlist_expires_at": 1701167767,
			},
			wantExpires: 0,
			wantReuse:   "true",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := resourceProject().TestResourceData()
			for k, v := range tc.raw {
				if err := d.Set(k, v); err != nil {
					t.Fatalf("set %s: %v", k, err)
				}
			}
			body := client.ProjectBody(d)
			if body.CveAllowlist.ExpiresAt != tc.wantExpires {
				t.Errorf("ExpiresAt = %d, want %d", body.CveAllowlist.ExpiresAt, tc.wantExpires)
			}
			if body.Metadata.ReuseSysCveAllowlist != tc.wantReuse {
				t.Errorf("ReuseSysCveAllowlist = %q, want %q", body.Metadata.ReuseSysCveAllowlist, tc.wantReuse)
			}
			out, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(out), `"expires_at"`) {
				t.Errorf("expires_at missing from request body: %s", out)
			}
		})
	}
}
