//go:build external_auth
// +build external_auth

package provider

import (
	"fmt"
	"strings"
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
	projectName, groupName := getProjectNGroupNames()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic(projectName, groupName),
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
	projectName, groupName := getProjectNGroupNames()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberGroupUpdate(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckMemberGroupBasic(projectName string, groupName string) string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "%s"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "%s"
	  role          = "developer"
	  type          = "oidc"
	}
	 
	`, projectName, groupName)
}

func testAccCheckMemberGroupUpdate(projectName string, groupName string) string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "%s"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "%s"
	  role          = "guest"
	  type          = "oidc"
	}

	`, projectName, groupName)
}

func getProjectNGroupNames() (string, string) {
	randStr := randomString(4)
	projectName := "acctest_project_" + strings.ToLower(randStr)
	groupName := "acctest_group_" + strings.ToLower(randStr)

	return projectName, groupName
}
