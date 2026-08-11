package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Harbor user group types, as returned/accepted by the /usergroups API.
const (
	GroupTypeLDAP     = 1
	GroupTypeInternal = 2
	GroupTypeOIDC     = 3
)

// GroupBody returns a JSON body built from resource data.
func GroupBody(d *schema.ResourceData) models.GroupBody {
	return models.GroupBody{
		Groupname:   d.Get("group_name").(string),
		GroupType:   d.Get("group_type").(int),
		LdapGroupDn: d.Get("ldap_group_dn").(string),
	}
}

// GroupTypeName returns the string form of a Harbor group type, inverse of
// GroupType. Unknown values return an empty string.
func GroupTypeName(t int) string {
	switch t {
	case GroupTypeLDAP:
		return "ldap"
	case GroupTypeInternal:
		return "internal"
	case GroupTypeOIDC:
		return "oidc"
	}
	return ""
}

// ListGroups fetches all Harbor user groups, paginating through the API.
// Harbor returns an empty page once the list is exhausted.
func (c *Client) ListGroups() ([]models.GroupBody, error) {
	var groups []models.GroupBody
	for page := 1; ; page++ {
		resp, _, _, err := c.SendRequest("GET", models.PathGroups+"?page="+strconv.Itoa(page), nil, 200)
		if err != nil {
			return nil, fmt.Errorf("list usergroups page %d: %w", page, err)
		}

		var pageData []models.GroupBody
		if err := json.Unmarshal([]byte(resp), &pageData); err != nil {
			return nil, fmt.Errorf("decode usergroups page %d: %w", page, err)
		}
		if len(pageData) == 0 {
			break
		}
		groups = append(groups, pageData...)
	}
	return groups, nil
}

// GetGroupByID returns a single usergroup by its numeric id.
func (c *Client) GetGroupByID(id int) (*models.GroupBody, error) {
	resp, _, code, err := c.SendRequest("GET", fmt.Sprintf("%s/%d", models.PathGroups, id), nil, 200)
	if err != nil {
		return nil, fmt.Errorf("get usergroup %d (status %d): %w", id, code, err)
	}
	var g models.GroupBody
	if err := json.Unmarshal([]byte(resp), &g); err != nil {
		return nil, fmt.Errorf("decode usergroup %d: %w", id, err)
	}
	// The /usergroups/{id} endpoint does not echo the id in the body, so
	// propagate the one we already know.
	g.ID = id
	return &g, nil
}

// CreateGroup creates a new usergroup and returns its numeric id, parsed from
// the Location response header.
func (c *Client) CreateGroup(g models.GroupBody) (int, error) {
	_, headers, _, err := c.SendRequest("POST", models.PathGroups, g, 201)
	if err != nil {
		return 0, fmt.Errorf("create usergroup %q: %w", g.Groupname, err)
	}
	loc, err := GetID(headers)
	if err != nil {
		return 0, fmt.Errorf("parse Location of created usergroup: %w", err)
	}
	// loc looks like "/usergroups/25".
	idx := strings.LastIndex(loc, "/")
	if idx < 0 || idx == len(loc)-1 {
		return 0, fmt.Errorf("unexpected usergroup Location %q", loc)
	}
	id, err := strconv.Atoi(loc[idx+1:])
	if err != nil {
		return 0, fmt.Errorf("parse usergroup id from Location %q: %w", loc, err)
	}
	return id, nil
}

// ResolveOrCreateLdapGroup looks up an LDAP-backed usergroup by its distinguished
// name, creating it if absent. It returns the resolved numeric id and whether
// an existing group was adopted (false means the group was just created).
//
// Harbor's POST /projects/{pid}/members endpoint has a known failure mode when
// called with member_group.ldap_group_dn: it creates the backing usergroup as
// a side effect but then returns HTTP 500 on the member attachment, leaving an
// orphan. Resolving the DN here and passing a numeric member_group.id to the
// members endpoint avoids that path entirely and also adopts any orphan left
// over from a prior failure.
func (c *Client) ResolveOrCreateLdapGroup(dn, name string) (id int, adopted bool, err error) {
	if dn == "" {
		return 0, false, fmt.Errorf("ldap group DN must not be empty")
	}
	if existing, err := c.lookupLdapGroupByDN(dn); err != nil {
		return 0, false, err
	} else if existing != 0 {
		return existing, true, nil
	}
	if name == "" {
		name = ShortNameFromDN(dn)
	}
	id, createErr := c.CreateGroup(models.GroupBody{
		Groupname:   name,
		GroupType:   GroupTypeLDAP,
		LdapGroupDn: dn,
	})
	if createErr == nil {
		return id, false, nil
	}
	// If another concurrent apply created the same group between our List
	// and our Create, Harbor rejects the second POST with an "already exist"
	// error. Re-list and try to adopt before surfacing the error.
	if existing, lookupErr := c.lookupLdapGroupByDN(dn); lookupErr == nil && existing != 0 {
		return existing, true, nil
	}
	return 0, false, createErr
}

// lookupLdapGroupByDN returns the id of the LDAP-backed usergroup matching dn,
// or 0 if no match is found. Comparison is case-insensitive.
func (c *Client) lookupLdapGroupByDN(dn string) (int, error) {
	groups, err := c.ListGroups()
	if err != nil {
		return 0, err
	}
	for _, g := range groups {
		if g.GroupType == GroupTypeLDAP && strings.EqualFold(g.LdapGroupDn, dn) {
			return g.ID, nil
		}
	}
	return 0, nil
}

// ShortNameFromDN extracts the value of the first RDN in a distinguished name,
// e.g. "cn=harbor_users,cn=groups,dc=example" -> "harbor_users". If the input
// is not a DN, it is returned unchanged.
func ShortNameFromDN(dn string) string {
	first := strings.SplitN(dn, ",", 2)[0]
	if eq := strings.Index(first, "="); eq >= 0 {
		return strings.TrimSpace(first[eq+1:])
	}
	return dn
}

// LooksLikeDN reports whether s looks like an LDAP distinguished name. This is
// a loose heuristic used only for back-compat: older configurations allowed a
// DN in `group_name`, and we want to keep accepting that transparently.
func LooksLikeDN(s string) bool {
	if !strings.Contains(s, "=") {
		return false
	}
	head := strings.ToLower(strings.SplitN(s, "=", 2)[0])
	head = strings.TrimSpace(head)
	switch head {
	case "cn", "ou", "dc", "uid", "o":
		return true
	}
	return false
}
