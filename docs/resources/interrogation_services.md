# Resource: harbor_interrogation_services

## Example Usage
```hcl
resource "harbor_interrogation_services" "main" {
  vulnerability_scan_policy = "daily"

}
```

## Argument Reference
The following arguments are supported:
* `default_scanner` - (Optional) Sets the default interrogation service **Clair**

* `vulnerability_scan_policy` - (Optional) The frequency of the vulnerability scanning is done. This can be `daily`, `weekly`, `monthly` or can be a custom cron string.