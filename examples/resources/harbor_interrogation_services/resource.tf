resource "harbor_interrogation_services" "main" {
  default_scanner = "Clair"
  vulnerability_scan_policy = "Daily"
}
