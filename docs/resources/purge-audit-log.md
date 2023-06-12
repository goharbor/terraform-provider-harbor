# Resource: harbor_purge_audit_log

## Example Usage
```hcl
resource "harbor_purge_audit_log" "main" {
  schedule              = "Daily"
  audit_retention_hour  = 24
  include_operations    = "create,pull"
}
```

## Argument Reference
The following arguments are supported:
* **schedule** - (Required) Sets the schedule how often the Garbage Collection will run.  Can be to `"Hourly"`, `"Daily"`, `"Weekly"` or can be a custom cron string ie, `"5 4 * * *"` 

* **audit_retention_hour** - (Required) to configure how long audit logs should be kept. For example, if you set this to 24 Harbor will only purge audit logs that are 24 or more hours old.

* **include_operations** - (Required) valid values are `create` `delete` `pull`, thoses values can be comma separated. When Create, Delete, or Pull is set, Harbor will include audit logs for those operations in the purge.