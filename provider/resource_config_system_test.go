package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConfigSystemBannerMessage(t *testing.T) {
	const initialMessage = "Registry status: \"maintenance\"\nPlease retry later."
	const updatedMessage = "Registry status: \"available\"\nAll services restored."
	config := func(message string, closable bool, messageType string) string {
		encoded, err := json.Marshal(message)
		if err != nil {
			t.Fatal(err)
		}
		return fmt.Sprintf(`
resource "harbor_config_system" "main" {
  disabled_audit_log_event_types = "create_artifact,delete_artifact"
  banner_message {
    message   = %s
    closable  = %t
    type      = %q
    from_date = "2026-11-01T00:00:00Z"
    to_date   = "2026-11-08T00:00:00Z"
  }
}`, encoded, closable, messageType)
	}
	check := func(message, closable, messageType string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckResourceExists(resourceConfigSystemMain),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.#", "1"),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.0.message", message),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.0.closable", closable),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.0.type", messageType),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.0.from_date", "2026-11-01T00:00:00Z"),
			resource.TestCheckResourceAttr(resourceConfigSystemMain, "banner_message.0.to_date", "2026-11-08T00:00:00Z"),
		)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConfigSystemDestroy,
		Steps: []resource.TestStep{
			{Config: config(initialMessage, true, "warning"), Check: check(initialMessage, "true", "warning")},
			{Config: config(updatedMessage, false, "info"), Check: check(updatedMessage, "false", "info")},
		},
	})
}

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
						resourceConfigSystemMain, "token_expiration", "30"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "session_timeout", "60"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "robot_name_prefix", "robot$"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "scanner_skip_update_pulltime", "false"),
					resource.TestCheckResourceAttr(
						resourceConfigSystemMain, "disabled_audit_log_event_types", "create_artifact,delete_artifact"),
				),
			},
		},
	})
}

func testAccCheckConfigSystem() string {
	return fmt.Sprintf(`
	resource "harbor_config_system" "main" {
		project_creation_restriction     = "adminonly"
		read_only                        = false
		robot_token_expiration           = 30
		token_expiration                 = 30
		session_timeout                  = 60
		robot_name_prefix                = "robot$"
		scanner_skip_update_pulltime     = false
		disabled_audit_log_event_types   = "create_artifact,delete_artifact"
	}
	`)
}
