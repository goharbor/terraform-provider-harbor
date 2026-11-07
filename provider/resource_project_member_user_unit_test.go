package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func projectMemberUserAPI(t *testing.T, putStatus int) *client.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/projects/1/members/5" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodPut:
			w.WriteHeader(putStatus)
			if putStatus != http.StatusOK {
				_, _ = w.Write([]byte(`{"errors":[{"code":"BAD_REQUEST","message":"invalid role"}]}`))
			}
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"id":5,"project_id":1,"role_id":3,"entity_id":12,"entity_name":"alice","entity_type":"u"}`))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	return client.NewClient(server.URL, "admin", "Harbor12345", "", "", true, "", nil)
}

func newProjectMemberUserResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, resourceMembersUser().Schema, map[string]interface{}{
		"project_id": "/projects/1",
		"user_name":  "alice",
		"role":       "developer",
	})
	d.SetId("/projects/1/members/5")
	return d
}

func TestResourceMembersUserUpdateReturnsAPIError(t *testing.T) {
	apiClient := projectMemberUserAPI(t, http.StatusBadRequest)
	d := newProjectMemberUserResourceData(t)

	if err := resourceMembersUserUpdate(d, apiClient); err == nil {
		t.Fatal("expected the failed PUT to be returned as an error")
	}
}

func TestResourceMembersUserReadSetsMemberID(t *testing.T) {
	apiClient := projectMemberUserAPI(t, http.StatusOK)
	d := newProjectMemberUserResourceData(t)

	if err := resourceMembersUserRead(d, apiClient); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got := d.Get("member_id").(int); got != 5 {
		t.Fatalf("expected member_id 5, got %d", got)
	}
	if got := d.Get("user_name").(string); got != "alice" {
		t.Fatalf("expected user_name alice, got %q", got)
	}
}
