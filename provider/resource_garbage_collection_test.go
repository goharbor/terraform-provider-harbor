//go:build external_auth
// +build external_auth

package provider

import (
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborGCMain = "harbor_garbage_collection.main"

func testAccCheckGCDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_garbage_collection" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 200)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}
		if resp != "" {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func TestAccGCUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGCBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborGCMain),
					resource.TestCheckResourceAttr(
						harborGCMain, "schedule", "Daily"),
				),
			},
			{
				Config: testAccCheckGCUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborGCMain),
					resource.TestCheckResourceAttr(
						harborGCMain, "schedule", "Hourly"),
					resource.TestCheckResourceAttr(
						harborGCMain, "delete_untagged", "true"),
				),
			},
		},
	})
}

func testAccCheckGCBasic() string {
	return fmt.Sprintf(`
	resource "harbor_garbage_collection" "main" {
		schedule        = "Daily"
	}
	`)
}

func testAccCheckGCUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_garbage_collection" "main" {
		schedule        = "Hourly"
		delete_untagged = true
	}
	`)
}
