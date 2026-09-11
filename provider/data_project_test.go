package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.harbor_project.test", "name", "test_data_source_project"),
					resource.TestCheckResourceAttr("data.harbor_project.test", "type", "Project"),
					resource.TestCheckResourceAttr("data.harbor_project.test", "vulnerability_scanning", "false"),
				),
			},
		},
	})
}

func testAccDataSourceProjectConfigBasic() string {
	return `
resource "harbor_project" "test" {
	name                   = "test_data_source_project"
	vulnerability_scanning = false
}

data "harbor_project" "test" {
	name = harbor_project.test.name
}
`
}

func TestAccDataSourceProject_vulnerabilityScanning(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectConfigVulnerabilityScanning(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.harbor_project.test", "name", "test_data_source_project_vulnerability_scanning"),
					resource.TestCheckResourceAttr("data.harbor_project.test", "type", "Project"),
					resource.TestCheckResourceAttr("data.harbor_project.test", "vulnerability_scanning", "true"),
				),
			},
		},
	})
}

func testAccDataSourceProjectConfigVulnerabilityScanning() string {
	return `
resource "harbor_project" "test" {
	name                   = "test_data_source_project_vulnerability_scanning"
	vulnerability_scanning = true
}

data "harbor_project" "test" {
	name = harbor_project.test.name
}
`
}
