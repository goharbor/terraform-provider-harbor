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
const enableContentTrustCosign = "enable_content_trust_cosign"

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

func TestAccProjectProxyCacheFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProjectProxyCacheFilter("doublestar", "library/{nginx,redis}"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "proxy_cache_filter_kind", "doublestar"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "proxy_cache_filter_pattern", "library/{nginx,redis}"),
				),
			},
			{
				Config: testAccCheckProjectProxyCacheFilter("regex", "^library/(nginx|redis)$"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborProjectMain),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "proxy_cache_filter_kind", "regex"),
					resource.TestCheckResourceAttr(
						resourceHarborProjectMain, "proxy_cache_filter_pattern", "^library/(nginx|redis)$"),
				),
			},
		},
	})
}

func testAccCheckProjectProxyCacheFilter(kind, pattern string) string {
	return fmt.Sprintf(`
	resource "harbor_registry" "main" {
		provider_name = "docker-hub"
		name          = "docker-hub-test-proxy-filter"
		endpoint_url  = "https://hub.docker.com"
	  }

	resource "harbor_project" "main" {
		name                       = "test_proxy_filter"
		registry_id                = harbor_registry.main.registry_id
		proxy_cache_filter_kind    = "%s"
		proxy_cache_filter_pattern = %q
	  }
	`, kind, pattern)
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
