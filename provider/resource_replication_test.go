package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccReplicationyBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckReplicationBasic(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckResourceExists("harbor_replication.pull"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull", "name", "test_pull"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull", "action", "pull"),
				),
			},
			{
				Config: testAccCheckReplicationUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_replication.pull"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull", "name", "test_pull"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull", "action", "pull"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull", "schedule", "0 0 0 * * *"),
				),
			},
		},
	})
}

func testAccCheckReplicationDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_replication" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckReplicationBasic() string {
	endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	config := fmt.Sprintf(`

resource "harbor_registry" "main" {
	provider_name = "harbor"
	name = "harbor-test-replication"
	endpoint_url = "%s"
  }

  resource "harbor_replication" "pull" {
	name  = "test_pull"
	action = "pull"
	registry_id = harbor_registry.main.registry_id

}`, endpoint)
	return config
}

func testAccCheckReplicationUpdate() string {
	endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	config := fmt.Sprintf(`

	resource "harbor_registry" "main" {
		provider_name = "harbor"
		name = "harbor-test"
		endpoint_url = "%s"
	  }

	  resource "harbor_replication" "pull" {
		name  = "test_pull"
		action = "pull"
		registry_id = harbor_registry.main.registry_id
		schedule = "0 0 0 * * *"
		
	}
	`, endpoint)
	return config
}
