package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/models"
)

func TestGroupTypeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		num  int
	}{
		{"ldap", GroupTypeLDAP},
		{"internal", GroupTypeInternal},
		{"oidc", GroupTypeOIDC},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := GroupType(c.name); got != c.num {
				t.Errorf("GroupType(%q) = %d, want %d", c.name, got, c.num)
			}
			if got := GroupTypeName(c.num); got != c.name {
				t.Errorf("GroupTypeName(%d) = %q, want %q", c.num, got, c.name)
			}
		})
	}
	if GroupTypeName(0) != "" || GroupTypeName(99) != "" {
		t.Errorf("GroupTypeName should return empty for unknown types")
	}
}

func TestLooksLikeDN(t *testing.T) {
	cases := map[string]bool{
		"cn=harbor_users,cn=groups,dc=example,dc=com": true,
		"CN=Admins,OU=Groups,DC=Example,DC=com":       true,
		"uid=alice,ou=people,dc=example,dc=com":       true,
		"harbor_users":                                false,
		"":                                            false,
		"name=value":                                  false, // unknown RDN head
	}
	for in, want := range cases {
		if got := LooksLikeDN(in); got != want {
			t.Errorf("LooksLikeDN(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestShortNameFromDN(t *testing.T) {
	cases := map[string]string{
		"cn=harbor_users,cn=groups,dc=example,dc=com": "harbor_users",
		"CN=Admins ,OU=x": "Admins",
		"plain":           "plain",
	}
	for in, want := range cases {
		if got := ShortNameFromDN(in); got != want {
			t.Errorf("ShortNameFromDN(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestProjectMembersGroupBodyByID(t *testing.T) {
	body := ProjectMembersGroupBodyByID("developer", 42)
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	// Critical: the DN-bearing fields must never leak into the members POST
	// body; that is the exact shape that triggers the Harbor 500 bug.
	for _, forbidden := range []string{"ldap_group_dn", "group_name", "group_type"} {
		if strings.Contains(s, forbidden) {
			t.Errorf("member body must not contain %q, got %s", forbidden, s)
		}
	}
	if !strings.Contains(s, `"id":42`) {
		t.Errorf("member body missing resolved group id, got %s", s)
	}
	if body.RoleID != 2 {
		t.Errorf("RoleID = %d, want 2 (developer)", body.RoleID)
	}
}

// fakeHarbor is a minimal stand-in for Harbor's /usergroups endpoints used by
// the resolver tests. It supports pagination, lookup and create.
type fakeHarbor struct {
	groups   []models.GroupBody
	nextID   int32
	creates  int32
	lastPost models.GroupBody
}

func (f *fakeHarbor) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/usergroups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Single page of results; return empty for any page != 1.
			page := r.URL.Query().Get("page")
			if page != "" && page != "1" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("[]"))
				return
			}
			_ = json.NewEncoder(w).Encode(f.groups)
		case http.MethodPost:
			var g models.GroupBody
			if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id := int(atomic.AddInt32(&f.nextID, 1))
			g.ID = id
			f.groups = append(f.groups, g)
			f.lastPost = g
			atomic.AddInt32(&f.creates, 1)
			w.Header().Set("Location", fmt.Sprintf("/api/v2.0/usergroups/%d", id))
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func newTestClient(t *testing.T, f *fakeHarbor) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(f.handler())
	c := NewClient(srv.URL, "user", "pass", "", "", false, "")
	return c, srv.Close
}

func TestResolveOrCreateLdapGroup_Adopts(t *testing.T) {
	f := &fakeHarbor{
		nextID: 10,
		groups: []models.GroupBody{
			{ID: 7, Groupname: "harbor_users", GroupType: GroupTypeLDAP, LdapGroupDn: "cn=harbor_users,dc=example,dc=com"},
		},
	}
	c, stop := newTestClient(t, f)
	defer stop()

	id, adopted, err := c.ResolveOrCreateLdapGroup("CN=harbor_users,DC=example,DC=com", "")
	if err != nil {
		t.Fatal(err)
	}
	if !adopted {
		t.Errorf("expected to adopt existing group")
	}
	if id != 7 {
		t.Errorf("id = %d, want 7", id)
	}
	if f.creates != 0 {
		t.Errorf("should not have created a group, creates=%d", f.creates)
	}
}

func TestResolveOrCreateLdapGroup_Creates(t *testing.T) {
	f := &fakeHarbor{nextID: 100}
	c, stop := newTestClient(t, f)
	defer stop()

	id, adopted, err := c.ResolveOrCreateLdapGroup("cn=new_group,dc=example,dc=com", "")
	if err != nil {
		t.Fatal(err)
	}
	if adopted {
		t.Errorf("expected to create, not adopt")
	}
	if id != 101 {
		t.Errorf("id = %d, want 101", id)
	}
	if f.lastPost.Groupname != "new_group" {
		t.Errorf("derived name = %q, want new_group", f.lastPost.Groupname)
	}
	if f.lastPost.GroupType != GroupTypeLDAP {
		t.Errorf("group_type = %d, want %d", f.lastPost.GroupType, GroupTypeLDAP)
	}
}

func TestResolveOrCreateLdapGroup_EmptyDN(t *testing.T) {
	f := &fakeHarbor{}
	c, stop := newTestClient(t, f)
	defer stop()
	if _, _, err := c.ResolveOrCreateLdapGroup("", ""); err == nil {
		t.Errorf("expected error for empty DN")
	}
}

func TestGetGroupByID(t *testing.T) {
	want := models.GroupBody{Groupname: "harbor_users", GroupType: GroupTypeLDAP, LdapGroupDn: "cn=harbor_users,dc=example,dc=com"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/usergroups/25" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "u", "p", "", "", false, "")
	got, err := c.GetGroupByID(25)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 25 || got.Groupname != want.Groupname || got.LdapGroupDn != want.LdapGroupDn {
		t.Errorf("got %+v, want id=25 name=%s dn=%s", got, want.Groupname, want.LdapGroupDn)
	}
}
