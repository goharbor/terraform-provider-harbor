package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckProjectDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)
	r := "harbor_project"

	for _, rs := range s.RootModule().Resources {
		if rs.Type != r {
			continue
		}

		resp, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
		if err != nil {
			return fmt.Errorf("Resouse was not delete \n %s", resp)
		}

	}

	return nil
}

func TestAccProjectBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProjectBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "name", "test_basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "public", "false"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "vulnerability_scanning", "false"),
				),
			},
		},
	})
}

func TestAccProjectUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProjectBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "name", "test_basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "public", "false"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "vulnerability_scanning", "false"),
				),
			},
			{
				Config: testAccCheckItemUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("harbor_project.basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "name", "test_basic"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "public", "true"),
					resource.TestCheckResourceAttr(
						"harbor_project.basic", "vulnerability_scanning", "true"),
				),
			},
		},
	})
}

func testAccCheckProjectBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "basic" {
		name = "test_basic"
		public = false
		vulnerability_scanning = false
	  }
	`)
}

func testAccCheckItemUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_project" "basic" {
		name = "test_basic"
		public = true
		vulnerability_scanning = true
	  }
`)
}
