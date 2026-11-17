package provider

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestGCScheduleDiff(t *testing.T) {
	cases := []struct {
		name     string
		state    string
		config   string
		wantDiff bool
	}{
		{"same value", "Weekly", "Weekly", false},
		{"lowercase preset", "Weekly", "weekly", false},
		{"uppercase preset", "Hourly", "HOURLY", false},
		{"cron matching preset", "Daily", "0 0 0 * * *", false},
		{"same custom cron", "0 5 4 * * *", "0 5 4 * * *", false},
		{"different preset", "Weekly", "daily", true},
		{"preset to custom cron", "Daily", "0 5 4 * * *", true},
		{"different custom cron", "0 5 4 * * *", "0 5 5 * * *", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID: "/system/gc/schedule",
				Attributes: map[string]string{
					"id":              "/system/gc/schedule",
					"schedule":        tc.state,
					"delete_untagged": "false",
					"workers":         "1",
				},
			}
			cfg := terraform.NewResourceConfigRaw(map[string]interface{}{"schedule": tc.config})
			diff, err := resourceGC().SimpleDiff(context.Background(), state, cfg, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := false
			if diff != nil {
				_, got = diff.Attributes["schedule"]
			}
			if got != tc.wantDiff {
				t.Errorf("schedule %q -> %q: diff = %v, want %v", tc.state, tc.config, got, tc.wantDiff)
			}
		})
	}
}

// Harbor stores delete_tag=false when the request does not carry it, so the
// request body must hold it whenever the configuration enables it.
func TestGCDeleteTagRequestBody(t *testing.T) {
	cases := []struct {
		name      string
		deleteTag bool
		want      bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceGC().Schema, map[string]interface{}{
				"schedule":        "Weekly",
				"delete_untagged": true,
				"delete_tag":      tc.deleteTag,
				"workers":         2,
			})
			body, err := json.Marshal(client.GetSystemBoby(d, "gc"))
			if err != nil {
				t.Fatal(err)
			}
			got := strings.Contains(string(body), `"delete_tag":true`)
			if got != tc.want {
				t.Errorf("delete_tag=%v: request body %s, delete_tag sent = %v, want %v", tc.deleteTag, body, got, tc.want)
			}
		})
	}
}
