package provider

import (
	"fmt"
	"regexp"
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
			{
				Config: testAccCheckRetentionScheduleNone(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborRetentionMain),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "schedule", ""),
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

func testAccCheckRetentionScheduleNone() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
			name                = "acctest_retention_pol"
	  }

	  resource "harbor_retention_policy" "main" {
		  scope = harbor_project.main.id
		  schedule = ""
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

func TestAccRetentionMissingTagSelector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionNoTagSelector(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionNoTagSelector() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_no_tag"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"
		rule {
			n_days_since_last_pull = 5
			repo_matching          = "**"
		}
	}
	`)
}

func TestAccRetentionMissingRepoSelector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionNoRepoSelector(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionNoRepoSelector() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_no_repo"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"
		rule {
			n_days_since_last_pull = 5
			tag_matching           = "latest"
		}
	}
	`)
}

func TestAccRetentionMissingRetainParam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionNoRetainParam(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionNoRetainParam() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_no_param"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"
		rule {
			repo_matching = "**"
			tag_matching  = "latest"
		}
	}
	`)
}

// ============================================================
// Multi-rule valid tests
// ============================================================

func TestAccRetentionMultipleRulesValid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRetentionMultipleRulesValid(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborRetentionMain),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "schedule", "Daily"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.n_days_since_last_pull", "7"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.repo_matching", "**"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.tag_matching", "latest"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.1.n_days_since_last_push", "14"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.1.repo_matching", "myrepo/**"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.1.tag_matching", "v*"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.1.disabled", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.2.n_days_since_last_pull", "30"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.2.repo_matching", "archive/**"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.2.tag_matching", "*"),
					resource.TestCheckResourceAttr(
						resourceHarborRetentionMain, "rule.2.disabled", "true"),
				),
			},
		},
	})
}

func testAccCheckRetentionMultipleRulesValid() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_multi"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"

		rule {
			n_days_since_last_pull = 7
			repo_matching          = "**"
			tag_matching           = "latest"
		}

		rule {
			n_days_since_last_push = 14
			repo_matching          = "myrepo/**"
			tag_matching           = "v*"
		}

		rule {
			n_days_since_last_pull = 30
			repo_matching          = "archive/**"
			tag_matching           = "*"
			disabled               = true
		}
	}
	`)
}

// ============================================================
// Multi-rule invalid tests
// ============================================================

func TestAccRetentionMultiRuleMissingTagSelector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionMultiRuleMissingTag(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionMultiRuleMissingTag() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_multi_notag"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"

		rule {
			n_days_since_last_pull = 7
			repo_matching          = "**"
			tag_matching           = "latest"
		}

		rule {
			n_days_since_last_push = 10
			repo_matching          = "**"
		}
	}
	`)
}

func TestAccRetentionMultiRuleMissingRepoSelector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionMultiRuleMissingRepo(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionMultiRuleMissingRepo() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_multi_norepo"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"

		rule {
			n_days_since_last_pull = 7
			repo_matching          = "**"
			tag_matching           = "latest"
		}

		rule {
			n_days_since_last_push = 10
			tag_matching           = "v*"
		}
	}
	`)
}

func TestAccRetentionMultiRuleMissingRetainParam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRetentionMultiRuleMissingRetain(),
				ExpectError: regexp.MustCompile(`one of`),
			},
		},
	})
}

func testAccCheckRetentionMultiRuleMissingRetain() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest_retention_multi_noretain"
	}

	resource "harbor_retention_policy" "main" {
		scope    = harbor_project.main.id
		schedule = "Daily"

		rule {
			n_days_since_last_pull = 7
			repo_matching          = "**"
			tag_matching           = "latest"
		}

		rule {
			repo_matching = "**"
			tag_matching  = "v*"
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
