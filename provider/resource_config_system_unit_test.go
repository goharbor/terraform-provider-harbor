package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestConfigSystemBannerMessageSurvivesQuotes applies a banner message holding
// a double quote and a backslash and reads back what was PUT to /configurations.
func TestConfigSystemBannerMessageSurvivesQuotes(t *testing.T) {
	const message = `Upgrade on "Friday", see \\fileserver\notes`

	var sent models.ConfigBodySystemPost
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(body, &sent)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	d := schema.TestResourceDataRaw(t, resourceConfigSystem().Schema, map[string]interface{}{
		"banner_message": []interface{}{
			map[string]interface{}{
				"closable":  true,
				"message":   message,
				"type":      "warning",
				"from_date": "2026-11-01T00:00:00Z",
				"to_date":   "2026-11-08T00:00:00Z",
			},
		},
	})

	apiClient := client.NewClient(server.URL, "admin", "Harbor12345", "", "", true, "", nil)
	if err := resourceConfigSystemCreate(d, apiClient); err != nil {
		t.Fatal(err)
	}

	var got models.BannerMessage
	if err := json.Unmarshal([]byte(sent.BannerMessage), &got); err != nil {
		t.Fatalf("banner_message is not valid JSON: %v\nsent: %s", err, sent.BannerMessage)
	}
	if got.Message != message {
		t.Fatalf("message was %q; want %q", got.Message, message)
	}
	if !got.Closable || got.Type != "warning" ||
		got.FromDate != "2026-11-01T00:00:00Z" || got.ToDate != "2026-11-08T00:00:00Z" {
		t.Fatalf("other banner fields were lost: %+v", got)
	}
}
