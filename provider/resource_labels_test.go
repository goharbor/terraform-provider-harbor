package provider

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const resourceHarborLabelMain = "harbor_label.main"

func testAccCheckLabelDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)
	r := "harbor_label"

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

func TestAccLabelBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLabelBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborLabelMain),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "name", "accTest"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "color", "#FFFFFF"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "description", "Description to for acceptance test"),
				),
			},
		},
	})
}

func TestAccLabelUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLabelBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborLabelMain),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "name", "accTest"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "color", "#FFFFFF"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "description", "Description to for acceptance test"),
				),
			},
			{
				Config: testAccCheckLabelUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborLabelMain),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "name", "accTest"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "color", "#FF0000"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "description", "Description to for acceptance test"),
				),
			},
		},
	})
}

func TestAccLabelProjectUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLabelProjectBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborLabelMain),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "name", "accTest"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "color", "#FFFFFF"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "description", "Description to for acceptance test"),
				),
			},
			{
				Config: testAccCheckLabelProjectUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceHarborLabelMain),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "name", "accTest"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "color", "#FF0000"),
					resource.TestCheckResourceAttr(
						resourceHarborLabelMain, "description", "Description to for acceptance test"),
				),
			},
		},
	})
}

func testAccCheckLabelBasic() string {
	return fmt.Sprintf(`
	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FFFFFF"
		description 	= "Description to for acceptance test"
	  }
	`)
}

func testAccCheckLabelUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FF0000"
		description 	= "Description to for acceptance test"
	}
	`)
}

func testAccCheckLabelProjectBasic() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest"
	}

	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FFFFFF"
		description 	= "Description to for acceptance test"
#		project_id	= harbor_project.main.id
	}
	`)
}

func testAccCheckLabelProjectUpdate() string {
	return fmt.Sprintf(`
	resource "harbor_project" "main" {
		name = "acctest"
	}

	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FF0000"
		description 	= "Description to for acceptance test"
	#	project_id	= harbor_project.main.id
	}
	`)
}
