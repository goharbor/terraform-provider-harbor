package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// getLdapURL returns the LDAP URL reachable from Harbor (inside Kind cluster).
// In CI, HARBOR_LDAP_URL points to the OpenLDAP container IP on the kind network.
// For local development, it falls back to ldap://localhost:389.
func getLdapURL() string {
	if v := os.Getenv("HARBOR_LDAP_URL"); v != "" {
		return v
	}
	return "ldap://localhost:389"
}

const resourceConfigAuthMain = "harbor_config_auth.main"

func testAccCheckConfigAuthDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_config_auth" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", "/configurations", nil, 200)
		if err != nil {
			return fmt.Errorf("resource was not deleted\n%s", resp)
		}
	}

	return nil
}

func TestAccConfigAuthDb(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConfigAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConfigAuthDb(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceConfigAuthMain),
					resource.TestCheckResourceAttr(
						resourceConfigAuthMain, "auth_mode", "db_auth"),
				),
			},
		},
	})
}

func testAccCheckConfigAuthDb() string {
	return fmt.Sprintf(`
	resource "harbor_config_auth" "main" {
		auth_mode = "db_auth"
	}
	`)
}

// LDAP tests
const resourceConfigAuthLdapMain = "harbor_config_auth.ldap"
const harborGroupLdapMain = "harbor_group.main"
const harborProjectMemberGroupLdapMain = "harbor_project_member_group.main"

func testAccCheckConfigAuthLdapDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_config_auth" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", "/configurations", nil, 200)
		if err != nil {
			return fmt.Errorf("resource was not deleted\n%s", resp)
		}
	}

	return nil
}

func TestAccConfigAuthLdap(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConfigAuthLdapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConfigAuthLdap(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceConfigAuthLdapMain),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "auth_mode", "ldap_auth"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_scope", "onelevel"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_base_dn", "ou=people,dc=planetexpress,dc=com"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_filter", "(objectClass=posixGroup)"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_gid", "cn"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_admin_dn", "cn=admin_staff,ou=people,dc=planetexpress,dc=com"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_membership", "memberof"),
					resource.TestCheckResourceAttr(
						resourceConfigAuthLdapMain, "ldap_group_scope", "subtree"),
				),
			},
		},
	})
}

func testAccCheckConfigAuthLdap() string {
	ldapURL := getLdapURL()
	return fmt.Sprintf(`
	resource "harbor_config_auth" "ldap" {
		auth_mode             = "ldap_auth"
		ldap_url              = "%s"
		ldap_base_dn          = "ou=people,dc=planetexpress,dc=com"
		ldap_uid              = "uid"
		ldap_verify_cert      = false
		ldap_search_dn        = "cn=admin,dc=planetexpress,dc=com"
		ldap_search_password  = "GoodNewsEveryone"
		ldap_scope            = "onelevel"
		ldap_group_base_dn    = "ou=people,dc=planetexpress,dc=com"
		ldap_group_filter     = "(objectClass=posixGroup)"
		ldap_group_gid        = "cn"
		ldap_group_admin_dn   = "cn=admin_staff,ou=people,dc=planetexpress,dc=com"
		ldap_group_membership = "memberof"
		ldap_group_scope      = "subtree"
	}
	`, ldapURL)
}

// ─── harbor_group LDAP tests ───────────────────────────────────────────────

func testAccCheckGroupLdapDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_group" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("resource was not deleted\n%s", resp)
		}
	}

	return nil
}

func TestAccGroupLdap(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupLdapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGroupLdap(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborGroupLdapMain),
					resource.TestCheckResourceAttr(
						harborGroupLdapMain, "group_name", "ship_crew"),
					resource.TestCheckResourceAttr(
						harborGroupLdapMain, "group_type", "1"),
					resource.TestCheckResourceAttr(
						harborGroupLdapMain, "ldap_group_dn", "cn=ship_crew,ou=people,dc=planetexpress,dc=com"),
				),
			},
		},
	})
}

func testAccCheckGroupLdap() string {
	ldapURL := getLdapURL()
	return fmt.Sprintf(`
	resource "harbor_config_auth" "ldap" {
		auth_mode            = "ldap_auth"
		ldap_url             = "%s"
		ldap_base_dn         = "ou=people,dc=planetexpress,dc=com"
		ldap_uid             = "uid"
		ldap_verify_cert     = false
		ldap_search_dn       = "cn=admin,dc=planetexpress,dc=com"
		ldap_search_password = "GoodNewsEveryone"
		ldap_group_base_dn   = "ou=people,dc=planetexpress,dc=com"
		ldap_group_filter    = "(objectClass=posixGroup)"
		ldap_group_gid       = "cn"
		ldap_group_admin_dn  = "cn=admin_staff,ou=people,dc=planetexpress,dc=com"
		ldap_group_membership = "memberof"
		ldap_group_scope     = "subtree"
	}

	resource "harbor_group" "main" {
		depends_on    = [harbor_config_auth.ldap]
		group_name    = "ship_crew"
		group_type    = 1
		ldap_group_dn = "cn=ship_crew,ou=people,dc=planetexpress,dc=com"
	}
	`, ldapURL)
}

// ─── harbor_project_member_group LDAP tests ────────────────────────────────

func testAccCheckMemberGroupLdapDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_project_member_group" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("resource was not deleted\n%s", resp)
		}
	}

	return nil
}

func TestAccMemberGroupLdap(t *testing.T) {
	randStr := strings.ToLower(randomString(4))
	projectName := "acctest_ldap_project_" + randStr

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupLdapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupLdap(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupLdapMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupLdapMain, "role", "developer"),
				),
			},
		},
	})
}

func TestAccMemberGroupLdapUpdate(t *testing.T) {
	randStr := strings.ToLower(randomString(4))
	projectName := "acctest_ldap_project_" + randStr

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupLdapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupLdap(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupLdapMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupLdapMain, "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberGroupLdapUpdate(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupLdapMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupLdapMain, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckMemberGroupLdap(projectName string) string {
	ldapURL := getLdapURL()
	return fmt.Sprintf(`
	resource "harbor_config_auth" "ldap" {
		auth_mode            = "ldap_auth"
		ldap_url             = "%s"
		ldap_base_dn         = "ou=people,dc=planetexpress,dc=com"
		ldap_uid             = "uid"
		ldap_verify_cert     = false
		ldap_search_dn       = "cn=admin,dc=planetexpress,dc=com"
		ldap_search_password = "GoodNewsEveryone"
		ldap_group_base_dn   = "ou=people,dc=planetexpress,dc=com"
		ldap_group_filter    = "(objectClass=posixGroup)"
		ldap_group_gid       = "cn"
		ldap_group_admin_dn  = "cn=admin_staff,ou=people,dc=planetexpress,dc=com"
		ldap_group_membership = "memberof"
		ldap_group_scope     = "subtree"
	}

	resource "harbor_project" "main" {
		name = "%s"
	}

	resource "harbor_project_member_group" "main" {
		depends_on    = [harbor_config_auth.ldap]
		project_id    = harbor_project.main.id
		group_name    = "ship_crew"
		role          = "developer"
		type          = "ldap"
		ldap_group_dn = "cn=ship_crew,ou=people,dc=planetexpress,dc=com"
	}
	`, ldapURL, projectName)
}

func testAccCheckMemberGroupLdapUpdate(projectName string) string {
	ldapURL := getLdapURL()
	return fmt.Sprintf(`
	resource "harbor_config_auth" "ldap" {
		auth_mode            = "ldap_auth"
		ldap_url             = "%s"
		ldap_base_dn         = "ou=people,dc=planetexpress,dc=com"
		ldap_uid             = "uid"
		ldap_verify_cert     = false
		ldap_search_dn       = "cn=admin,dc=planetexpress,dc=com"
		ldap_search_password = "GoodNewsEveryone"
		ldap_group_base_dn   = "ou=people,dc=planetexpress,dc=com"
		ldap_group_filter    = "(objectClass=posixGroup)"
		ldap_group_gid       = "cn"
		ldap_group_admin_dn  = "cn=admin_staff,ou=people,dc=planetexpress,dc=com"
		ldap_group_membership = "memberof"
		ldap_group_scope     = "subtree"
	}

	resource "harbor_project" "main" {
		name = "%s"
	}

	resource "harbor_project_member_group" "main" {
		depends_on    = [harbor_config_auth.ldap]
		project_id    = harbor_project.main.id
		group_name    = "ship_crew"
		role          = "guest"
		type          = "ldap"
		ldap_group_dn = "cn=ship_crew,ou=people,dc=planetexpress,dc=com"
	}
	`, ldapURL, projectName)
}
