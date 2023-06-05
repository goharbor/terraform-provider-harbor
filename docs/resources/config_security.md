# Resource: harbor_config_security

## Example Usage

```hcl
resource "harbor_config_security" "main" {
  cve_allowlist = ["CVE-456", "CVE-123"]
  expires_at = "1701167767"
}
```

## Argument Reference

The following arguments are supported:

* `cve_allowlist` - (Required) System allowlist. Vulnerabilities in this list will be ignored when pushing and pulling images. Should be in the format or `["CVE-123", "CVE-145"]` or `["CVE-123"]`

* `expires_at` - (Optional) The time for expiration of the allowlist, in the form of seconds since epoch. This is an optional attribute, if it's not set the CVE allowlist does not expire.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the resource.
* `creation_time` - Time of creation of the list.
* `update_time` - Time of update of the list.

## Import

The list can be imported using the `id` eg,

`
terraform import harbor_config_security.main "7"
`

> Note that at this point of time Harbor doesn't has any api endpoint for deleting this list. Only updating the records.
