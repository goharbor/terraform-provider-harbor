package provider

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourcePurgeAuditMain = "harbor_purge_audit_log.main"

func TestAccPurgeAuditUpdate(t *testing.T) {
	harborVersion := getTestHarborVersion()

	var config string
	var check resource.TestCheckFunc

	if compareVersion(harborVersion, "2.13.0") < 0 {
		config = testAccCheckPurgeAuditBasic()
		check = resource.ComposeTestCheckFunc(
			testAccCheckResourceExists(resourcePurgeAuditMain),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "schedule", "Daily"),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "audit_retention_hour", "24"),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "include_operations", "create,pull"),
		)
	} else {
		config = testAccCheckPurgeAuditEventTypes()
		check = resource.ComposeTestCheckFunc(
			testAccCheckResourceExists(resourcePurgeAuditMain),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "schedule", "Daily"),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "audit_retention_hour", "24"),
			resource.TestCheckResourceAttr(resourcePurgeAuditMain, "include_event_types", "create,pull"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}

// Helper to compare semantic versions
func compareVersion(v1, v2 string) int {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")
	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(s1[i])
		n2, _ := strconv.Atoi(s2[i])
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}

func getTestHarborVersion() string {
	return os.Getenv("HARBOR_APP_VERSION")
}

func testAccCheckPurgeAuditBasic() string {
	return `
    resource "harbor_purge_audit_log" "main" {
        schedule              = "Daily"
        audit_retention_hour  = 24
        include_operations    = "create,pull"
    }
    `
}

func testAccCheckPurgeAuditEventTypes() string {
	return `
    resource "harbor_purge_audit_log" "main" {
        schedule              = "Daily"
        audit_retention_hour  = 24
        include_event_types   = "create,pull"
    }
    `
}
