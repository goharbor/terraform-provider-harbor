package provider

import (
	"fmt"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceHarborProjectMain = "harbor_project.main"
const enableContentTrust = "enable_content_trust"

func testAccCheckProjectDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)
	r := "harbor_project"

	for _, rs := range s.RootModule().Resources {
		if rs.Type != r {
			continue
		}

		resp, _, _, err := apiClient.SendRequest("GET", rs.Primary.ID, nil, 404)
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
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "name", "test_basic"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "public", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "vulnerability_scanning", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, enableContentTrust, "false"),
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
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "name", "test_basic"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "public", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "vulnerability_scanning", "false"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, enableContentTrust, "false"),
				),
			},
			{
				Config: testAccCheckItemUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "name", "test_basic"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "public", "true"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "vulnerability_scanning", "true"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, enableContentTrust, "true"),
				),
			},
		},
	})
}

func TestAccProjectCVEAllowlist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProjectCVEAllowlist(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "cve_allowlist.#", "2"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "cve_allowlist.0", "CVE-2023-12345"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "cve_allowlist.1", "CVE-2024-12345"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "cve_allowlist_expires_at", "2000000000"),
				),
			},
		},
	})
}

func testAccCheckProjectBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test_basic"
		public = false
		vulnerability_scanning = false
		%v = false
	  }
	`, enableContentTrust)
}

func testAccCheckItemUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "test_basic"
		public = true
		vulnerability_scanning = true
		%v = true
	  }
`, enableContentTrust)
}

func testAccCheckProjectCVEAllowlist() string {
	return `
	resource "harbor_project" "main" {
		name = "test_cve_allowlist"
		cve_allowlist = [
			"CVE-2023-12345",
			"CVE-2024-12345",
		]
		cve_allowlist_expires_at = 2000000000
	}
	`
}
