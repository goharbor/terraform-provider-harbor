package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResolveMemberGroupID(t *testing.T) {
	const dn = "cn=team,ou=groups,dc=example,dc=com"
	groups := []models.GroupBody{
		{ID: 4, Groupname: "team", GroupType: client.GroupTypeOIDC},
		{ID: 5, Groupname: "team", GroupType: client.GroupTypeInternal},
		{ID: 7, Groupname: "team", GroupType: client.GroupTypeLDAP, LdapGroupDn: dn},
	}
	cases := []struct {
		name         string
		config       map[string]interface{}
		status       int
		want         int
		wantErr      string
		wantRequests int
	}{
		{name: "LDAP short name", config: map[string]interface{}{"type": "ldap", "group_name": "team"}, want: 7, wantRequests: 3},
		{name: "missing LDAP group", config: map[string]interface{}{"type": "ldap", "group_name": "missing"}, wantErr: `usergroup "missing" of type "ldap" not found`, wantRequests: 3},
		{name: "no LDAP identifier", config: map[string]interface{}{"type": "ldap"}, wantErr: "requires either group_id, group_name or ldap_group_dn"},
		{name: "name is case sensitive", config: map[string]interface{}{"type": "ldap", "group_name": "TEAM"}, wantErr: `usergroup "TEAM" of type "ldap" not found`, wantRequests: 3},
		{name: "name must match exactly", config: map[string]interface{}{"type": "ldap", "group_name": "tea"}, wantErr: `usergroup "tea" of type "ldap" not found`, wantRequests: 3},
		{name: "explicit ID takes precedence", config: map[string]interface{}{"type": "ldap", "group_id": 42, "group_name": "missing", "ldap_group_dn": dn}, want: 42},
		{name: "explicit DN takes precedence", config: map[string]interface{}{"type": "ldap", "group_name": "missing", "ldap_group_dn": dn}, want: 7, wantRequests: 3},
		{name: "legacy DN in name", config: map[string]interface{}{"type": "ldap", "group_name": dn}, want: 7, wantRequests: 3},
		{name: "OIDC name", config: map[string]interface{}{"type": "oidc", "group_name": "team"}, want: 4, wantRequests: 3},
		{name: "internal name", config: map[string]interface{}{"type": "internal", "group_name": "team"}, want: 5, wantRequests: 3},
		{name: "API failure", config: map[string]interface{}{"type": "ldap", "group_name": "team"}, status: 500, wantErr: "list usergroups page 1", wantRequests: 1},
		{name: "permission denied", config: map[string]interface{}{"type": "ldap", "group_name": "team"}, status: 403, wantErr: "forbidden", wantRequests: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requests := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if r.Method != http.MethodGet || r.URL.Path != "/usergroups" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					http.Error(w, "unexpected request", 500)
					return
				}
				if tc.status != 0 {
					w.WriteHeader(tc.status)
					return
				}
				// The LDAP match is on page two, after same-name groups of other types.
				var page []models.GroupBody
				switch r.URL.Query().Get("page") {
				case "1":
					page = groups[:2]
				case "2":
					page = groups[2:]
				}
				_ = json.NewEncoder(w).Encode(page)
			}))
			defer srv.Close()
			api := client.NewClient(srv.URL, "u", "p", "", "", false, "", nil)
			d := schema.TestResourceDataRaw(t, resourceMembersGroup().Schema, tc.config)
			got, err := resolveMemberGroupID(api, d)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("got id=%d err=%v; want error containing %q", got, err, tc.wantErr)
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("got id=%d err=%v; want id=%d", got, err, tc.want)
			}
			if requests != tc.wantRequests {
				t.Errorf("requests=%d, want %d", requests, tc.wantRequests)
			}
		})
	}
}

func TestMembersGroupCreateLDAPByName(t *testing.T) {
	const dn = "cn=team,ou=groups,dc=example,dc=com"
	group := models.GroupBody{ID: 7, Groupname: "team", GroupType: client.GroupTypeLDAP, LdapGroupDn: dn}
	var requests []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/usergroups":
			groups := []models.GroupBody{}
			if r.URL.Query().Get("page") == "1" {
				groups = append(groups, group)
			}
			_ = json.NewEncoder(w).Encode(groups)
		case r.Method == http.MethodPost && r.URL.Path == "/projects/1/members":
			var body struct {
				Role  int                    `json:"role_id"`
				Group map[string]interface{} `json:"member_group"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Role != 3 || len(body.Group) != 1 || body.Group["id"] != float64(7) {
				t.Errorf("unexpected member body: %+v", body)
			}
			w.Header().Set("Location", "/api/v2.0/projects/1/members/9")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && r.URL.Path == "/projects/1/members/9":
			_ = json.NewEncoder(w).Encode(models.ProjectMembersBodyResponses{ID: 9, ProjectID: 1, RoleID: 3, EntityID: 7, EntityType: "g", EntityName: "team"})
		case r.Method == http.MethodGet && r.URL.Path == "/usergroups/7":
			_ = json.NewEncoder(w).Encode(group)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			http.Error(w, "unexpected request", 500)
		}
	}))
	defer srv.Close()
	api := client.NewClient(srv.URL, "u", "p", "", "", false, "", nil)
	d := schema.TestResourceDataRaw(t, resourceMembersGroup().Schema, map[string]interface{}{
		"project_id": "/projects/1", "role": "guest", "type": "ldap", "group_name": "team",
	})
	if err := resourceMembersGroupCreate(d, api); err != nil {
		t.Fatal(err)
	}
	if d.Id() != "/projects/1/members/9" {
		t.Errorf("membership ID=%q", d.Id())
	}
	for key, want := range map[string]interface{}{"group_id": 7, "group_name": "team", "ldap_group_dn": dn, "type": "ldap", "role": "guest", "member_id": 9} {
		if got := d.Get(key); got != want {
			t.Errorf("%s=%v, want %v", key, got, want)
		}
	}
	want := "[GET /usergroups?page=1 GET /usergroups?page=2 POST /projects/1/members GET /projects/1/members/9 GET /usergroups/7]"
	if fmt.Sprint(requests) != want {
		t.Errorf("requests=%v, want %s", requests, want)
	}
}
