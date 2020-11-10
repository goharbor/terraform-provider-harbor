# Resource: harbor_garbage_collection

## Example Usage
```hcl
resource "harbor_garbage_collection" "main" {
  schedule         = "Daily"
  delete_untagged  = true

}
```

## Argument Reference
The following arguments are supported:
* **schedule** - (Required) Sets the schedule how often the Garbage Collection will run.  Can be to `"hourly"`, `"daily"`, `"weekly"` or can be a custom cron string ie, `"5 4 * * *"` 

* **delete_untagged** - (Optional) Allow garbage collection on untagged artifacts.