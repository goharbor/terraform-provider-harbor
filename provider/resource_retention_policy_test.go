package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceHarborRetentionMain = "harbor_retention_policy.main"

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
						resourceHarborRetentionMain, "schedule", "Daily"),
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

// func TestDestinationNamespace(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		// CheckDestroy: testAccCheckLabelDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testReplicationPolicyDestinationNamespace(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckResourceExists(resourceHarborRetentionMain),
// 					resource.TestCheckResourceAttr(
// 						resourceHarborRetentionMain, "schedule", scheduleType),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckRetentionBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name                = "acctest_retention_pol"
	  }
	  
	  resource "harbor_retention_policy" "main" {
		  scope = harbor_project.main.id
		  schedule = "Daily"
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

// func testReplicationPolicyDestinationNamespace() string {
// 	return fmt.Sprintf(`
// 	resource "harbor_project" "main" {
// 		name                = "acctest_retention_pol"
// 	  }

// 	  resource "harbor_retention_policy" "main" {
// 		  scope = harbor_project.main.id
// 		  schedule = "event_base"
// 		  rule {
// 			  n_days_since_last_pull = 5
// 			  repo_matching = "**"
// 			  tag_matching = "latest"
// 		  }
// 		  rule {
// 			  n_days_since_last_push = 10
// 			  repo_matching = "**"
// 			  tag_matching = "latest"
// 		  }

// 	  }
// 	`)
// }
