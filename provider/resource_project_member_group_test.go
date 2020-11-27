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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic(),
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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMemberGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberGroupUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberGroupMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberGroupMain, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckMemberGroupBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "testing"
	  role          = "developer"
	  type          = "oidc"
	}
	 
	`)
}

func testAccCheckMemberGroupUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test"
	}
	
	resource "harbor_project_member_group" "main" {
	  project_id    = harbor_project.main.id
	  group_name    = "testing"
	  role          = "guest"
	  type          = "oidc"
	}

	`)
}
