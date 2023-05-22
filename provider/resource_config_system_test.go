package provider

import (
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceConfigSystemMain = "harbor_config_system.main"

func testAccCheckConfigSystemDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)
	r := "harbor_config_system"

	for _, rs := range s.RootModule().Resources {
		if rs.Type != r {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", "/configurations", nil, 200)
		if err != nil {
			return fmt.Errorf("Resource was not deleted\n%s", resp)
		}
	}

	return nil
}

func TestAccConfigSystem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConfigSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConfigSystem(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceConfigSystemMain),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "project_creation_restriction", "adminonly"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "read_only", "false"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "robot_token_expiration", "30"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "robot_name_prefix", "robot$"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "scanner_skip_update_pulltime", "false"),
				),
			},
		},
	})
}

func testAccCheckConfigSystem() string {
	return fmt.Sprintf(`
	resource "harbor_config_system" "main" {
		project_creation_restriction  = "adminonly"
		read_only                     = false
		robot_token_expiration        = 30
		robot_name_prefix             = "robot$"
		scanner_skip_update_pulltime  = false
	}
	`)
}
