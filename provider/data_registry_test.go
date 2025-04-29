package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRegistry_basic(t *testing.T) {
	for _, rName := range []string{"Foo", "Foo bar"} {
		t.Run(rName, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccDataSourceRegistryConfig_basicDataSource(rName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.harbor_registry.test", "name", rName),
							resource.TestCheckResourceAttr("data.harbor_registry.test", "type", "docker-hub"),
							resource.TestCheckResourceAttr("data.harbor_registry.test", "url", "https://hub.docker.com"),
						),
					},
				},
			})
		})
	}
}

func testAccDataSourceRegistryConfig_basicDataSource(rName string) string {
	return fmt.Sprintf(`
resource "harbor_registry" "test" {
	provider_name = "docker-hub"
    name = "%s"
	endpoint_url = "https://hub.docker.com"
}

data "harbor_registry" "test" {
    name = harbor_registry.test.name
}
`, rName)
}
