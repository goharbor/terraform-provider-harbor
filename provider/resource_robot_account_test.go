package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborRobotAccount = "harbor_robot_account.main"

func TestAccRobotSystem(t *testing.T) {
	randStr := randomString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotSystem("acctest_robot_" + strings.ToLower(randStr)),
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
	randStr := randomString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotProject("acctest_robot_" + strings.ToLower(randStr)),
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

func TestAccRobotProjectWriteOnlySecret(t *testing.T) {
	randStr := randomString(4)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotProjectWriteOnlySecret("acctest_robot_" + strings.ToLower(randStr)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.main"),

					testAccCheckResourceExists(harborRobotAccount),
					resource.TestCheckResourceAttr(
						harborRobotAccount, "name", "test_robot_project_wo"),
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

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckRobotSystem(projectName string) string {
	return fmt.Sprintf(`
	resource "harbor_robot_account" "main" {
	  name        = "test_robot_system"
	  description = "system level robot account"
	  level       = "system"
	  permissions {
		access {
		  action   = "create"
		  resource = "label"
		}
		kind      = "system"
		namespace = "/"
	  }
	  permissions {
		access {
		  action   = "push"
		  resource = "repository"
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
	  name = "%v"
	}
	`, projectName)
}

func testAccCheckRobotProject(projectName string) string {
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
	  name = "%v"
	}
	`, projectName)
}

func testAccCheckRobotProjectWriteOnlySecret(projectName string) string {
	return fmt.Sprintf(`
	resource "harbor_robot_account" "main" {
	  name              = "test_robot_project_wo"
	  description       = "project level robot account with write-only secret"
	  level             = "project"
	  secret_wo         = "robotSecret12345"
	  secret_wo_version = 1
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
	  name = "%v"
	}
	`, projectName)
}
