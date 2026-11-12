package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestMembersGroupCreateOIDCCreatesMissingGroup covers the regression from
// #592: granting a project to an OIDC group nobody from has logged in yet
// must create the usergroup instead of failing with "not found".
func TestMembersGroupCreateOIDCCreatesMissingGroup(t *testing.T) {
	var created *models.GroupBody
	var requests []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/usergroups":
			groups := []models.GroupBody{}
			if r.URL.Query().Get("page") == "1" && created != nil {
				groups = append(groups, *created)
			}
			_ = json.NewEncoder(w).Encode(groups)
		case r.Method == http.MethodPost && r.URL.Path == "/usergroups":
			var g models.GroupBody
			if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
				t.Error(err)
			}
			if g.Groupname != "team.dev" || g.GroupType != client.GroupTypeOIDC {
				t.Errorf("unexpected usergroup body: %+v", g)
			}
			g.ID = 7
			created = &g
			w.Header().Set("Location", "/api/v2.0/usergroups/7")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPost && r.URL.Path == "/projects/1/members":
			var body struct {
				Role  int                    `json:"role_id"`
				Group map[string]interface{} `json:"member_group"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Role != 2 || len(body.Group) != 1 || body.Group["id"] != float64(7) {
				t.Errorf("unexpected member body: %+v", body)
			}
			w.Header().Set("Location", "/api/v2.0/projects/1/members/9")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && r.URL.Path == "/projects/1/members/9":
			_ = json.NewEncoder(w).Encode(models.ProjectMembersBodyResponses{ID: 9, ProjectID: 1, RoleID: 2, EntityID: 7, EntityType: "g", EntityName: "team.dev"})
		case r.Method == http.MethodGet && r.URL.Path == "/usergroups/7":
			_ = json.NewEncoder(w).Encode(created)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			http.Error(w, "unexpected request", http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	api := client.NewClient(srv.URL, "u", "p", "", "", false, "", nil)
	d := schema.TestResourceDataRaw(t, resourceMembersGroup().Schema, map[string]interface{}{
		"project_id": "/projects/1", "role": "developer", "type": "oidc", "group_name": "team.dev",
	})
	if err := resourceMembersGroupCreate(d, api); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "/projects/1/members/9" {
		t.Errorf("membership ID=%q", d.Id())
	}
	for key, want := range map[string]interface{}{"group_id": 7, "group_name": "team.dev", "type": "oidc", "role": "developer", "member_id": 9} {
		if got := d.Get(key); got != want {
			t.Errorf("%s=%v, want %v", key, got, want)
		}
	}
	want := "[GET /usergroups?page=1 POST /usergroups POST /projects/1/members GET /projects/1/members/9 GET /usergroups/7]"
	if fmt.Sprint(requests) != want {
		t.Errorf("requests=%v, want %s", requests, want)
	}
}
