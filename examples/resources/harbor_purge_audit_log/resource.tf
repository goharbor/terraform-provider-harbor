resource "harbor_purge_audit_log" "main" {
  schedule              = "Daily"
  audit_retention_hour  = 24
  include_operations    = "create,pull"
}
