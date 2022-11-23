package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborReplicationPull = "harbor_replication.pull"

func TestAccReplicationyBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckReplicationBasic(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckResourceExists(harborReplicationPull),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "name", "test_pull"),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "action", "pull"),
				),
			},
			{
				Config: testAccCheckReplicationUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborReplicationPull),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "name", "test_pull"),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "action", "pull"),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "schedule", "0 0 0 * * *"),
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

		resp, _, _, err := apiClient.SendRequest(context.Background(), "GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resource was not deleted \n %s", resp)
		}

	}

	return nil
}

func testAccCheckReplicationBasic() string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	config := fmt.Sprintf(`
	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name          = "docker-hub-test-replication"
		endpoint_url  = "%s"
	  }
	  
	resource "harbor_replication" "pull" {
		name        = "test_pull"
		action      = "pull"
		registry_id = harbor_registry.main.registry_id
	  }`, endpoint)
	return config
}

func testAccCheckReplicationUpdate() string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	config := fmt.Sprintf(`

	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name = "docker-hub-test-replication"
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

func TestDestinationNamespace(t *testing.T) {
	var scheduleType = "* 0/15 * * * *"
	var destNamepace = "gcp-project"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testReplicationPolicyDestinationNamespace(scheduleType, destNamepace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborReplicationPull),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "schedule", scheduleType),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "dest_namespace", destNamepace),
				),
			},
		},
	})
}

func TestDestinationNamespaceReplaceCount(t *testing.T) {
	var scheduleType = "* 0/15 * * * *"
	var destNamespaceReplaceCount = 0

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testReplicationPolicyDestinationNamespaceWithReplaceCount(scheduleType, destNamespaceReplaceCount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborReplicationPull),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "schedule", scheduleType),
					resource.TestCheckResourceAttr(
						harborReplicationPull, "dest_namespace_replace", fmt.Sprintf("%d", destNamespaceReplaceCount)),
				),
			},
		},
	})
}

func testReplicationPolicyDestinationNamespace(scheduleType string, destNamepace string) string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	return fmt.Sprintf(`
	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name = "docker-hub-test-rep-pol"
		endpoint_url = "%s"
	  }

	  resource "harbor_replication" "pull" {
		name  = "test_pull"
		action = "pull"
		registry_id = harbor_registry.main.registry_id
		schedule = "%s"
		dest_namespace = "%s"
	  }
	`, endpoint, scheduleType, destNamepace)
}

func testReplicationPolicyDestinationNamespaceWithReplaceCount(scheduleType string, destNamepace int) string {
	// endpoint := os.Getenv("HARBOR_REPLICATION_ENDPOINT")
	endpoint := "https://hub.docker.com"
	return fmt.Sprintf(`
	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name = "docker-hub-test-rep-pol"
		endpoint_url = "%s"
	  }

	  resource "harbor_replication" "pull" {
		name  = "test_pull"
		action = "pull"
		registry_id = harbor_registry.main.registry_id
		schedule = "%s"
		dest_namespace = "nobody_cares"
		dest_namespace_replace = "%d"
	  }
	`, endpoint, scheduleType, destNamepace)
}
