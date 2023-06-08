package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourcePurgeAuditMain = "harbor_purge_audit_log.main"

func TestAccPurgeAuditUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPurgeAuditBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourcePurgeAuditMain),
					resource.TestCheckResourceAttr(
						resourcePurgeAuditMain, "schedule", "Daily"),
					resource.TestCheckResourceAttr(
						resourcePurgeAuditMain, "audit_retention_hour", "24"),
					resource.TestCheckResourceAttr(
						resourcePurgeAuditMain, "include_operations", "create, pull"),
				),
			},
		},
	})
}

func testAccCheckPurgeAuditBasic() string {
	return fmt.Sprintf(`
	resource "purge_audit" "main" {
		schedule              = "Daily"
		audit_retention_hour  = 24
		include_operations    = "create,pull"
	}
	`)
}
