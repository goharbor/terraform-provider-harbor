# Resource: harbor_interrogation_services

## Example Usage
```
resource "harbor_interrogation_services" "main" {
  vulnerability_scan_policy = "daily"

}
```

## Argument Reference
The following arguments are supported:
* **default_scanner** - (Optinal) Sets the default interrogation service **Clair**

* **vulnerability_scan_policy** - (Optional) The frequency of the vulnerbility scanning is done. Can be to **"hourly"**, **"daily"** or **"weekly"**