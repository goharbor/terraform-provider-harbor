package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceHarborRetentionMain = "harbor_retention_policy.main"

// func testAccCheckRetentionDestroy(s *terraform.State) error {
// 	apiClient := testAccProvider.Meta().(*client.Client)
// 	r := "harbor_retention_policy"

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != r {
// 			continue
// 		}

// 		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
// 		if err != nil {
// 			return fmt.Errorf("Resouse was not delete \n %s", resp)
// 		}

// 	}

// 	return nil
// }

func TestAccRetentionUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRetentionBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborRetentionMain),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "schedule", "daily"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.n_days_since_last_pull", "5"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.repo_matching", "**"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.tag_matching", "latest"),
				),
			},
			// {
			// 	Config: testAccCheckLabelUpdate(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckResourceExists(resourceHarborRetentionMain),
			// 		resource.TestCheckResourceAttr(
			// 			resourceHarborRetentionMain, "name", "accTest"),
			// 		resource.TestCheckResourceAttr(
			// 			resourceHarborRetentionMain, "color", "#FF0000"),
			// 		resource.TestCheckResourceAttr(
			// 			resourceHarborRetentionMain, "description", "Description to for acceptance test"),
			// 	),
			// },
		},
	})
}

func TestDestinationNamespace(t *testing.T) {
	var scheduleType = "event_based"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testReplicationPolicyDestinationNamespace(scheduleType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborRetentionMain),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "schedule", scheduleType),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "dest_namespace", scheduleType),
				),
			},
		},
	})
}

func testAccCheckRetentionBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name                = "acctest_retention_pol"
	  }
	  
	  resource "harbor_retention_policy" "main" {
		  scope = harbor_project.main.id
		  schedule = "daily"
		  rule {
			  n_days_since_last_pull = 5
			  repo_matching = "**"
			  tag_matching = "latest"
		  }
		  rule {
			  n_days_since_last_push = 10
			  repo_matching = "**"
			  tag_matching = "latest"
		  }
	  
	  }
	`)
}

func testReplicationPolicyDestinationNamespace(scheduleType string) string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name                = "acctest_retention_pol"
	  }
	  
	  resource "harbor_retention_policy" "main" {
		  scope = harbor_project.main.id
		  schedule = "event_base"
		  dest_namespace = "%s"
		  rule {
			  n_days_since_last_pull = 5
			  repo_matching = "**"
			  tag_matching = "latest"
		  }
		  rule {
			  n_days_since_last_push = 10
			  repo_matching = "**"
			  tag_matching = "latest"
		  }
	  
	  }
	`, scheduleType)
}
