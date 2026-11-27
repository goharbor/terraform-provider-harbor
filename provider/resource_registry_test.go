package provider

import (
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

func TestAccRegistryWriteOnlySecret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRegistryWriteOnlySecret("docker-hub-test-reg-wo", 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborRegistryMain),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "name", "docker-hub-test-reg-wo"),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "access_id", "registryUser"),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "access_secret", ""),
					resource.TestCheckNoResourceAttr(
						harborRegistryMain, "access_secret_wo"),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "access_secret_wo_version", "1"),
				),
			},
			{
				Config: testAccCheckRegistryWriteOnlySecret("docker-hub-test-reg-wo", 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborRegistryMain),
					resource.TestCheckNoResourceAttr(
						harborRegistryMain, "access_secret_wo"),
					resource.TestCheckResourceAttr(
						harborRegistryMain, "access_secret_wo_version", "2"),
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

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
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

func testAccCheckRegistryWriteOnlySecret(name string, secretVersion int) string {
	endpoint := "https://hub.docker.com"
	return fmt.Sprintf(`
	resource "harbor_registry" "main" {
		provider_name            = "docker-hub"
		name                     = "%s"
		endpoint_url             = "%s"
		access_id                = "registryUser"
		access_secret_wo         = "registrySecret%d"
		access_secret_wo_version = %d
	  }
	`, name, endpoint, secretVersion, secretVersion)
}
