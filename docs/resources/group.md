# Resource: harbor_group

## Example Usage
```hcl

resource "harbor_group" "storage-group" {
  group_name = "storage-group"
  group_type = 3
}
```

## Argument Reference
The following arguments are supported:

* **group_name** - (Required) The name of the group.

* **group_type** - (Required) 3. Note: group type 3 is OIDC group.

* **ldap_group_dn** - (Optional) The distinguished name of the group within AD/LDAP 

## Import
An OIDC group can be imported using the `group id` eg,

`
terraform import harbor_group.storage-group /usergroups/19
`
