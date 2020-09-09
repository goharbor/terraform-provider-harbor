// +build external_auth

package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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
					testAccCheckResourceExists("harbor_project_member_group.main"),
					resource.TestCheckResourceAttr(
						"harbor_project_member_group.main", "role", "developer"),
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
					testAccCheckResourceExists("harbor_project_member_group.main"),
					resource.TestCheckResourceAttr(
						"harbor_project_member_group.main", "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberGroupUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project_member_group.main"),
					resource.TestCheckResourceAttr(
						"harbor_project_member_group.main", "role", "guest"),
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
	  group_name    = "8582dd52-1da5-4afe-94fe-25b55097d43a"
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
	  group_name    = "8582dd52-1da5-4afe-94fe-25b55097d43a"
	  role          = "guest"
	  type          = "oidc"
	}
	 
	  
	`)
}
