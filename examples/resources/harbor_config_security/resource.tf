resource "harbor_config_security" "main" {
  cve_allowlist = ["CVE-456", "CVE-123"]
  expires_at = "1701167767"
}
