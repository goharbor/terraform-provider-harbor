// +build internal_auth

package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const harborProjectMemberUserMain = "harbor_project_member_user.main"

func testAccCheckMemberUserDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_project_member_user" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func TestAccMemberUserBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberUserMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberUserMain, "role", "developer"),
				),
			},
		},
	})
}

func TestAccMemberUserUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMemberUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberUserMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberUserMain, "role", "developer"),
				),
			},
			{
				Config: testAccCheckMemberUserUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectMemberUserMain),
					resource.TestCheckResourceAttr(
						harborProjectMemberUserMain, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckMemberUserBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test"
	  }

	  resource "harbor_user" "main" {
		username  = "john"
		password  = "Password12345!"
		full_name = "John"
		email     = "john@contoso.com"
	  }

	  resource "harbor_project_member_user" "main" {
		project_id = harbor_project.main.id
		role       = "developer"
		user_name = harbor_user.main.username
	  }	  
	  
	`)
}

func testAccCheckMemberUserUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test"
	  }

	  resource "harbor_user" "main" {
		username  = "john"
		password  = "Password12345!"
		full_name = "John"
		email     = "john@contoso.com"
	  }

	  resource "harbor_project_member_user" "main" {
		project_id = harbor_project.main.id
		role       = "guest"
		user_name = harbor_user.main.username
	  }	  
	  
	`)
}
