package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const harborConfigEmail = "harbor_config_email.main"

func TestAccConfigEmail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRobotBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(harborConfigEmail),
					resource.TestCheckResourceAttr(
						harborConfigEmail, "email_host", "server.acme.com"),
				),
			},
		},
	})
}

// stAccCheckRtDestroy(s *terraform.State) error {
// 	apiClient := testAccProvider.Meta().(*client.Client)

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "harbor_config_email" {
// 			continue
// 		}

// 		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
// 		if err != nil {
// 			return fmt.Errorf("Resouse was not delete \n %s", resp)
// 		}

// 	}

// 	return nil
// }

func testAccCheckConfigEmail() string {
	return fmt.Sprintf(`

	resource "harbor_config_email" "main" {
		email_host = "server.acme.com"
		email_from = "dont_reply@acme.com"
	}
	  
	`)
}
