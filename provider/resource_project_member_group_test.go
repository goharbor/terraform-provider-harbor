//go:build external_auth
// +build external_auth

package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborProjectMemberGroupMain = "harbor_project_member_group.main"

func testAccCheckMemberGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_project_member_group" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func TestAccMemberGroupBasic(t *testing.T) {
	randStr := RandString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic("acctest_project_"+strings.ToLower(randStr), "acctest_group_"+strings.ToLower(randStr)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "developer"),
				),
			},
		},
	})
}

func TestAccMemberGroupUpdate(t *testing.T) {
	randStr := RandString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic("acctest_project_"+strings.ToLower(randStr), "acctest_group_"+strings.ToLower(randStr)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberGroupUpdate("acctest_project_"+strings.ToLower(randStr), "acctest_group_"+strings.ToLower(randStr)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckMemberGroupBasic(projectName, groupName string) string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "%v"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "%v"
	  role          = "developer"
	  type          = "oidc"
	}
	 
	`, projectName, groupName)
}

func testAccCheckMemberGroupUpdate(projectName, groupName string) string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "%v"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "%v"
	  role          = "guest"
	  type          = "oidc"
	}

	`, projectName, groupName)
}
