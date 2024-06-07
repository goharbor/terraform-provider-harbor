package provider

import (
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborPreheatInstanceMain = "harbor_preheat_instance.main"

func testAccCheckPreheatInstanceDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_preheat_instance" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 200)
		if err != nil {
			return fmt.Errorf("Resource was not deleted \n %s", resp)
		}
	}

	return nil
}

func TestAccPreheatInstanceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPreheatInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPreheatInstanceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborPreheatInstanceMain),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "name", "test"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "vendor", "dragonfly"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "endpoint", "http://example.com"),
				),
			},
			{
				Config: testAccCheckPreheatInstanceUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborPreheatInstanceMain),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "name", "test-updated"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "vendor", "kraken"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "endpoint", "http://example-updated.com"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "auth_mode", "BASIC"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "username", "test-user"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "password", "test-password"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "default", "true"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "enabled", "false"),
					resource.TestCheckResourceAttr(
						harborPreheatInstanceMain, "insecure", "true"),
				),
			},
		},
	})
}

func testAccCheckPreheatInstanceBasic() string {
	return fmt.Sprintf(`
	resource "harbor_preheat_instance" "main" {
		name     = "test"
		vendor   = "dragonfly"
		endpoint = "http://example.com"
	}
	`)
}

func testAccCheckPreheatInstanceUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_preheat_instance" "main" {
		name       = "test-updated"
		vendor     = "kraken"
		endpoint   = "http://example-updated.com"
		auth_mode  = "BASIC"
		username   = "test-user"
		password   = "test-password"
		default    = true
		enabled    = false
		insecure   = true
	}
	`)
}
