# Data Source: harbor_groups

## Example Usage
```hcl
data "harbor_groups" "example" {
  group_name = "example-group"
}

output "group_ids" {
  value = [data.harbor_groups.example.*.id]
}
```

## Argument Reference
The following arguments are supported:

* `group_name` - (Optional) The name of the group to filter by.
* `ldap_group_dn` - (Optional) The LDAP group DN to filter by.

## Attributes Reference
In addition to all arguments, the following attributes are exported:

* `groups` - (Computed) A list of groups matching the previous arguments. Each `group` object provides the attributes documented below.

---

**group** object exports the following:

* `id` - The ID of the group.
* `group_name` - The name of the group.
* `group_type` - The type of the group.
* `ldap_group_dn` - The LDAP group DN of the group.

This data source retrieves a list of Harbor groups and filters them based on the `group_name` and `ldap_group_dn` arguments. It returns a list of `group` objects, each containing the `id`, `group_name`, `group_type`, and `ldap_group_dn` attributes of a group that matches the filter criteria.