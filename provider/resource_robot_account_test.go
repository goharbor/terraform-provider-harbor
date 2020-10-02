package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const harborRobotAccount = "harbor_robot_account.main"

func TestAccRobotBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.main"),

					testAccCheckResourceExists(harborRobotAccount),
					resource.TestCheckResourceAttr(
						harborRobotAccount, "name", "test_robot_account"),
					// resource.TestCheckResourceAttr(
					// 	harborRobotAccount, "action", "push"),
				),
			},
		},
	})
}

func TestAccRobotMultipleAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotMultipleAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.main"),

					testAccCheckResourceExists(harborRobotAccount),
					resource.TestCheckResourceAttr(
						harborRobotAccount, "name", "test_robot_account"),
					// resource.TestCheckResourceAttr(
					// 	harborRobotAccount, "action", "push"),
				),
			},
		},
	})
}

func testAccCheckRobotDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_robot_account" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckRobotBasic() string {
	return fmt.Sprintf(`

	resource "harbor_robot_account" "main" {
		name        = "test_robot_account"
		description = "Robot account to be used to push images"
		project_id  = harbor_project.main.id
		actions     = ["pull"]

	  }

	  resource "harbor_project" "main" {
		name = "test_basic"
	  }
	  
	`)
}

func testAccCheckRobotMultipleAction() string {
	return fmt.Sprintf(`

	resource "harbor_robot_account" "main" {
		name        = "test_robot_account"
		description = "Robot account to be used to push images"
		project_id  = harbor_project.main.id
		actions      = ["push","read","create"]
	  }

	  resource "harbor_project" "main" {
		name = "test_basic"
	  }
	  
	`)
}
