package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborRegistryMain = "harbor_registry.main"

func TestAccRegistryBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRegistryBasic(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckResourceExists(harborRegistryMain),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "name", "docker-hub-test-reg"),
				),
			},
			{
				Config: testAccCheckRegistryUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborRegistryMain),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "name", "docker-hub-test-update"),
				),
			},
		},
	})
}

func testAccCheckRegistryDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_registry" {
			continue
		}

		resp, _, _, err := apiClient.SendRequest(context.Background(), "GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckRegistryBasic() string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	config := fmt.Sprintf(`

	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name = "docker-hub-test-reg"
		endpoint_url = "%s"
	  }

	`, endpoint)
	return config
}

func testAccCheckRegistryUpdate() string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	config := fmt.Sprintf(`

	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name = "docker-hub-test-update"
		endpoint_url = "%s"
	  }

	`, endpoint)
	return config
}
