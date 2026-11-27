package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func newRetentionResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()

	return schema.TestResourceDataRaw(t, resourceRetention().Schema, map[string]interface{}{
		"scope": "/projects/7",
		"rule": []interface{}{
			map[string]interface{}{
				"n_days_since_last_pull": 5,
				"repo_matching":          "**",
				"tag_matching":           "latest",
			},
		},
	})
}

// retentionAPI is a fake Harbor that rejects every retention POST and reports the retention id of project 7.
func retentionAPI(t *testing.T, retentionID string) (*client.Client, func() []string) {
	t.Helper()

	var mu sync.Mutex
	var requests []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/retentions"):
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors":[{"code":"BAD_REQUEST","message":"invalid cron string for scheduled tag retention"}]}`))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/projects/7"):
			_, _ = w.Write([]byte(`{"metadata":{"retention_id":"` + retentionID + `"}}`))
		case r.Method == http.MethodPut:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	apiClient := client.NewClient(server.URL, "admin", "Harbor12345", "", "", true, "", nil)

	return apiClient, func() []string {
		mu.Lock()
		defer mu.Unlock()

		return append([]string(nil), requests...)
	}
}

func TestRetentionCreateReportsTheCreateErrorWhenTheProjectHasNoPolicy(t *testing.T) {
	apiClient, requests := retentionAPI(t, "")

	err := resourceRetentionCreate(newRetentionResourceData(t), apiClient)

	if err == nil || !strings.Contains(err.Error(), "invalid cron string") {
		t.Fatalf("expected the create error to be reported, got %v", err)
	}
	for _, request := range requests() {
		if strings.HasPrefix(request, "PUT") {
			t.Fatalf("expected no update without an existing policy, got %v", requests())
		}
	}
}

func TestRetentionCreateUpdatesTheExistingPolicy(t *testing.T) {
	apiClient, requests := retentionAPI(t, "42")
	d := newRetentionResourceData(t)

	// The read after the update is not under test; only the update request is.
	_ = resourceRetentionCreate(d, apiClient)

	var updated bool
	for _, request := range requests() {
		if strings.HasPrefix(request, "PUT") && strings.HasSuffix(request, "/retentions/42") {
			updated = true
		}
	}
	if !updated {
		t.Fatalf("expected an update of the existing policy, got %v", requests())
	}
	if d.Id() == "" {
		t.Fatal("expected the existing policy id to be stored")
	}
}

func TestRetentionScopeDiff(t *testing.T) {
	cases := []struct {
		name     string
		state    string
		config   string
		wantDiff bool
	}{
		{"same prefixed scope", "/projects/7", "/projects/7", false},
		{"plain project id", "/projects/7", "7", false},
		{"other prefixed scope", "/projects/7", "/projects/8", true},
		{"other plain project id", "/projects/7", "8", true},
		{"project id sharing a prefix", "/projects/7", "70", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID: "/retentions/1",
				Attributes: map[string]string{
					"id":                            "/retentions/1",
					"scope":                         tc.state,
					"schedule":                      "",
					"rule.#":                        "1",
					"rule.0.n_days_since_last_pull": "5",
					"rule.0.repo_matching":          "**",
					"rule.0.tag_matching":           "latest",
					"rule.0.disabled":               "false",
					"rule.0.untagged_artifacts":     "true",
				},
			}
			cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
				"scope": tc.config,
				"rule": []interface{}{
					map[string]interface{}{
						"n_days_since_last_pull": 5,
						"repo_matching":          "**",
						"tag_matching":           "latest",
					},
				},
			})
			diff, err := resourceRetention().SimpleDiff(context.Background(), state, cfg, nil)
			if err != nil {
				t.Fatal(err)
			}
			got := false
			if diff != nil {
				_, got = diff.Attributes["scope"]
			}
			if got != tc.wantDiff {
				t.Errorf("scope %q -> %q: diff = %v, want %v", tc.state, tc.config, got, tc.wantDiff)
			}
		})
	}
}
