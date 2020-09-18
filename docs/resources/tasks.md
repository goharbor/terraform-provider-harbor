# Resource: harbor_tasks

## Example Usage
```hcl
resource "harbor_tasks" "main" {
  vulnerability_scan_policy = "daily"
}
```

## Argument Reference
The following arguments are supported:

* **vulnerability_scan_policy** - (Optional) The frequency of the vulnerbility scanning is done. Can be to **"hourly"**, **"daily"** or **"weekly"**