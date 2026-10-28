package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceProjects_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectsConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.harbor_projects.test", "name", "test_data_source_projects"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.name", "test_data_source_projects"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.type", "Project"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.public", "false"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.vulnerability_scanning", "false"),
				),
			},
		},
	})
}

func testAccDataSourceProjectsConfigBasic() string {
	return `
resource "harbor_project" "test" {
	name                   = "test_data_source_projects"
	public                 = false
	vulnerability_scanning = false
}

data "harbor_projects" "test" {
	name = harbor_project.test.name
}
`
}

func TestAccDataSourceProjects_vulnerabilityScanning(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectsConfigVulnerabilityScanning(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.harbor_projects.test", "name", "test_data_source_projects_vulnerability_scanning"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "vulnerability_scanning", "true"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.name", "test_data_source_projects_vulnerability_scanning"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.type", "Project"),
					resource.TestCheckResourceAttr("data.harbor_projects.test", "projects.0.vulnerability_scanning", "true"),
				),
			},
		},
	})
}

func testAccDataSourceProjectsConfigVulnerabilityScanning() string {
	return `
resource "harbor_project" "test" {
	name                   = "test_data_source_projects_vulnerability_scanning"
	vulnerability_scanning = true
}

data "harbor_projects" "test" {
	name                   = harbor_project.test.name
	vulnerability_scanning = true
}
`
}
