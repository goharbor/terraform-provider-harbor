package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborRobotAccount = "harbor_robot_account.main"

func TestAccRobotSystem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotSystem(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.main"),

					testAccCheckResourceExists(harborRobotAccount),
					resource.TestCheckResourceAttr(
						harborRobotAccount, "name", "test_robot_system"),
				),
			},
		},
	})
}

func TestAccRobotProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotProject(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.main"),

					testAccCheckResourceExists(harborRobotAccount),
					resource.TestCheckResourceAttr(
						harborRobotAccount, "name", "test_robot_project"),
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
		if rs.Type != "harbor_project" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckRobotSystem() string {
	return fmt.Sprintf(`
	resource "harbor_robot_account" "main" {
	  name        = "test_robot_system"
	  description = "system level robot account"
	  level       = "system"
	  permissions {
		access {
		  action   = "push"
		  resource = "repository"
		}
		access {
		  action   = "read"
		  resource = "helm-chart"
		}
		access {
		  action   = "read"
		  resource = "helm-chart-version"
		}
		kind      = "project"
		namespace = harbor_project.main.name
	  }
	  permissions {
		access {
		  action   = "pull"
		  resource = "repository"
		}
		kind      = "project"
		namespace = "*"
	  }
	}

	resource "harbor_project" "main" {
	  name = "test_basic"
	}
	`)
}

func testAccCheckRobotProject() string {
	return fmt.Sprintf(`
	resource "harbor_robot_account" "main" {
	  name        = "test_robot_project"
	  description = "project level robot account"
	  level       = "project"
	  permissions {
		access {
		  action   = "pull"
		  resource = "repository"
		}
		access {
		  action   = "push"
		  resource = "repository"
		}
		kind      = "project"
		namespace = harbor_project.main.name
	  }
	}

	resource "harbor_project" "main" {
	  name = "test_basic"
	}
	`)
}
