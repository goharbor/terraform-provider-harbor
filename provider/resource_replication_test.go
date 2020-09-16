package provider

import (
	"fmt"
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

					testAccCheckResourceExists("harbor_replication.push"),
					resource.TestCheckResourceAttr(
						"harbor_replication.push", "name", "test_push"),
					resource.TestCheckResourceAttr(
						"harbor_replication.push", "action", "push"),
				),
			},
			// {
			// 	Config: testAccCheckReplicationUpdate(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckResourceExists("harbor_replication.push"),
			// 		resource.TestCheckResourceAttr(
			// 			"harbor_replication.main", "name", "test_push"),
			// 	),
			// },
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

	return fmt.Sprintf(`

	resource "harbor_registry" "main" {
		provider_name = "harbor"
		name = "harbor-test-replication"
		endpoint_url = "https://harbor-dev.bestsellerit.com"
	  }

	  resource "harbor_replication" "push" {
		name  = "test_push"
		action = "push"
		registry_id = harbor_registry.main.registry_id

	}
	
	`)
}

// func testAccCheckReplicationUpdate() string {

// 	return fmt.Sprintf(`

// 	resource "harbor_registry" "main" {
// 		provider_name = "harbor"
// 		name = "harbor-test"
// 		endpoint_url = "https://harbor-dev.bestsellerit.com"
// 	  }

// 	  resource "harbor_replication" "push" {
// 		name  = "test_push"
// 		action = "push"
// 		registry_id = harbor_registry.main.registry_id
// 	}
// 	`)
// }
