package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/goharbor/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembersGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// Back-compat: older configurations put an LDAP DN directly
				// into group_name. The current Read populates group_name with
				// the short form returned by the /usergroups endpoint, which
				// would otherwise produce a permanent diff against those old
				// configs. Suppress the diff when old and new both point at
				// the same backing LDAP usergroup.
				DiffSuppressFunc: suppressLdapGroupNameDiff,
			},
			"group_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"group_id", "group_name", "ldap_group_dn"},
			},
			"ldap_group_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validRoles := []string{"projectadmin", "developer", "guest", "maintainer", "limitedguest"}
					for _, r := range validRoles {
						if v == r {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validRoles, v))
					return
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "ldap" && v != "internal" && v != "oidc" {
						errs = append(errs, fmt.Errorf("%q must be either ldap, internal or oidc, got: %s", key, v))
					}
					return
				},
			},
		},
		Create: resourceMembersGroupCreate,
		Read:   resourceMembersGroupRead,
		Update: resourceMembersGroupUpdate,
		Delete: resourceMembersGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// ldapGroupNameEquivalent reports whether two group_name values refer to the
// same backing LDAP usergroup — i.e. one is the distinguished name and the
// other is its first-RDN short form, or they are identical. It is pure logic
// so it can be unit-tested without constructing a schema.ResourceData.
func ldapGroupNameEquivalent(groupType, ldapDN, old, new string) bool {
	if groupType != "ldap" || ldapDN == "" {
		return false
	}
	short := client.ShortNameFromDN(ldapDN)
	matches := func(v string) bool { return v == ldapDN || v == short }
	return matches(old) && matches(new)
}

// suppressLdapGroupNameDiff is the schema-level wrapper around
// ldapGroupNameEquivalent. It keeps existing state backwards-compatible for
// users who historically put the DN directly into group_name before
// ldap_group_dn existed or was enforced.
func suppressLdapGroupNameDiff(_, old, new string, d *schema.ResourceData) bool {
	return ldapGroupNameEquivalent(
		d.Get("type").(string),
		d.Get("ldap_group_dn").(string),
		old, new,
	)
}

// resolveMemberGroupID resolves the configured group_id / group_name /
// ldap_group_dn inputs to a concrete Harbor usergroup id, creating the
// usergroup on the fly for LDAP-backed groups when it does not exist yet.
//
// For type=ldap we never pass ldap_group_dn through to the members endpoint —
// that path triggers a server-side 500 on known Harbor versions. See
// Client.ResolveOrCreateLdapGroup for details.
func resolveMemberGroupID(apiClient *client.Client, d *schema.ResourceData) (int, error) {
	if id, ok := d.GetOk("group_id"); ok {
		return id.(int), nil
	}
	groupType := d.Get("type").(string)
	groupName := d.Get("group_name").(string)
	ldapDN := d.Get("ldap_group_dn").(string)
	if groupType == "ldap" {
		return resolveLdapMemberGroupID(apiClient, ldapDN, groupName)
	}
	return resolveNamedMemberGroupID(apiClient, groupType, groupName)
}

// resolveLdapMemberGroupID handles type=ldap by resolving or creating the
// backing usergroup and returning its numeric id. A DN supplied in group_name
// is accepted for back-compat with older configurations.
func resolveLdapMemberGroupID(apiClient *client.Client, ldapDN, groupName string) (int, error) {
	if ldapDN == "" && client.LooksLikeDN(groupName) {
		log.Printf("[WARN] harbor_project_member_group: treating group_name %q as ldap_group_dn for back-compat; please move the DN to the ldap_group_dn field", groupName)
		ldapDN = groupName
		groupName = ""
	}
	if ldapDN == "" {
		return 0, fmt.Errorf("type=ldap requires either group_id or ldap_group_dn")
	}
	id, adopted, err := apiClient.ResolveOrCreateLdapGroup(ldapDN, groupName)
	if err != nil {
		return 0, err
	}
	action := "created"
	if adopted {
		action = "adopted existing"
	}
	log.Printf("[DEBUG] harbor_project_member_group: %s ldap usergroup id=%d dn=%q", action, id, ldapDN)
	return id, nil
}

// resolveNamedMemberGroupID handles type=internal and type=oidc by looking up
// an existing Harbor usergroup by name. It does not create usergroups; callers
// must manage them with harbor_group.
func resolveNamedMemberGroupID(apiClient *client.Client, groupType, groupName string) (int, error) {
	if groupName == "" {
		return 0, fmt.Errorf("type=%s requires either group_id or group_name", groupType)
	}
	groups, err := apiClient.ListGroups()
	if err != nil {
		return 0, err
	}
	wantType := client.GroupType(groupType)
	for _, g := range groups {
		if g.GroupType == wantType && g.Groupname == groupName {
			return g.ID, nil
		}
	}
	return 0, fmt.Errorf("usergroup %q of type %q not found; create it first with harbor_group", groupName, groupType)
}

func resourceMembersGroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := checkProjectid(d.Get("project_id").(string))
	path := projectid + "/members"

	groupID, err := resolveMemberGroupID(apiClient, d)
	if err != nil {
		return err
	}

	body := client.ProjectMembersGroupBodyByID(d.Get("role").(string), groupID)

	_, headers, _, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return fmt.Errorf("create project member (group_id=%d): %w", groupID, err)
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceMembersGroupRead(d, m)
}

func resourceMembersGroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, respCode, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if respCode == 404 && err != nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	var member models.ProjectMembersBodyResponses
	if err := json.Unmarshal([]byte(resp), &member); err != nil {
		return fmt.Errorf("decode project member %s: %w", d.Id(), err)
	}

	d.Set("role", client.RoleTypeNumber(member.RoleID))
	d.Set("project_id", checkProjectid(strconv.Itoa(member.ProjectID)))
	d.Set("member_id", member.ID)

	// For group members, use the backing usergroup as the source of truth so
	// that `type`, `group_name` and `ldap_group_dn` stay stable across refresh
	// and import and don't drift back to the short form Harbor returns on the
	// members endpoint.
	if member.EntityType == "g" && member.EntityID > 0 {
		g, err := apiClient.GetGroupByID(member.EntityID)
		if err != nil {
			return err
		}
		d.Set("group_id", g.ID)
		d.Set("group_name", g.Groupname)
		d.Set("ldap_group_dn", g.LdapGroupDn)
		if typeName := client.GroupTypeName(g.GroupType); typeName != "" {
			d.Set("type", typeName)
		}
		return nil
	}

	// Fallback for older Harbor responses that omit entity_id: keep the
	// previous best-effort behaviour.
	d.Set("group_name", member.EntityName)
	return nil
}

func resourceMembersGroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	// Only the role is mutable on an existing member; everything else is
	// either ForceNew or derived from the backing usergroup.
	body := models.ProjectMembersBodyPost{
		RoleID: client.RoleType(d.Get("role").(string)),
	}
	if _, _, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200); err != nil {
		return fmt.Errorf("update project member %s: %w", d.Id(), err)
	}

	return resourceMembersGroupRead(d, m)
}

func resourceMembersGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, respCode, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if respCode != 404 && err != nil { // We can't delete something that doesn't exist. Hence the 404-check
		return err
	}
	return nil
}
