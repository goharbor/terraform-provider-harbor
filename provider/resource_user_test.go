// +build internal_auth

package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceHarborUserMain = "harbor_user.main"

func testAccCheckUserDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "harbor_user" {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func TestAccUserBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborUserMain),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "username", "john"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "full_name", "John Smith"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "email", "john.smith@contoso.com"),
				),
			},
		},
	})
}

func TestAccUserUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborUserMain),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "username", "john"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "full_name", "John Smith"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "email", "john.smith@contoso.com"),
				),
			},
			{
				Config: testAccCheckUserUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborUserMain),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "username", "john"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "full_name", "John"),
					resource.TestCheckResourceAttr(
						resourceHarborUserMain, "email", "john@contoso.com"),
				),
			},
		},
	})
}

func testAccCheckUserBasic() string {
	return fmt.Sprintf(`
	resource "harbor_user" "main" {
		username  = "john"
		password  = "Password12345"
		full_name = "John Smith"
		email     = "john.smith@contoso.com"
	  }
	`)
}

func testAccCheckUserUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_user" "main" {
		username  = "john"
		password  = "Password12345!"
		full_name = "John"
		email     = "john@contoso.com"
	  }
	`)
}
