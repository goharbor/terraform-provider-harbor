package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const harborProjectWebhook = "harbor_project_webhook.main"

func TestAccProjectWebhook(t *testing.T) {
	randStr := randomString(6)
	projectName := "acctest_webhook_" + strings.ToLower(randStr)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProjectWebhook(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project_webhook.main"),

					testAccCheckResourceExists(harborProjectWebhook),
					resource.TestCheckResourceAttr(
						harborProjectWebhook, "name", "acctest_webhook"),
				),
			},
			{
				Config: testAccCheckProjectWebhookUpdate(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborProjectWebhook),
					resource.TestCheckResourceAttr(
						harborProjectWebhook, "name", "acctest_webhook"),
				),
			},
		},
	})
}

func testAccCheckProjectWebhookDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_project_webhook" {
			continue
		}
		if rs.Type != "harbor_project" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func testAccCheckProjectWebhook(projectName string) string {
	return fmt.Sprintf(`

	resource "harbor_project_webhook" "main" {
		name        = "acctest_webhook"
		address     = "https://webhook.domain.com"
		project_id  = harbor_project.main.id
		notify_type = "http"
	  
		events_types = [
		  "DELETE_ARTIFACT",
		  "PULL_ARTIFACT",
		  "PUSH_ARTIFACT",
		  "DELETE_CHART",
		  "DOWNLOAD_CHART",
		  "UPLOAD_CHART",
		  "QUOTA_EXCEED",
		  "QUOTA_WARNING",
		  "REPLICATION",
		  "SCANNING_FAILED",
		  "SCANNING_COMPLETED",
		  "TAG_RETENTION"
		]
	  
	}

	resource "harbor_project" "main" {
	  name = "%v"
	}
	`, projectName)
}

func testAccCheckProjectWebhookUpdate(projectName string) string {
	return fmt.Sprintf(`

	resource "harbor_project_webhook" "main" {
		name        = "acctest_webhook"
		address     = "https://webhook.domain.com"
		project_id  = harbor_project.main.id
		notify_type = "http"
	  
		events_types = [
		  "DELETE_ARTIFACT",
		  "PULL_ARTIFACT",
		  "PUSH_ARTIFACT",
		]
	  
	}

	resource "harbor_project" "main" {
	  name = "%v"
	}
	`, projectName)
}
