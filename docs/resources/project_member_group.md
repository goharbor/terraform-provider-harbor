# Resource: harbor_project_member_group

## Example Usage
```hcl
resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_project_member_group" "main" {
  project_id    = harbor_project.main.project_id
  group_name    = "testing1"
  role          = "master"
  type          = "oidc"
}

```

## Argument Reference
The following arguments are supported:

* **group_name** - (Required) The name of the group member entity

* **project_id** - (Required) The project id of the project that the entity will have access to.

* **role** - (Required) The premissions that the entity will be granted.

* **type** - (Required) The group type.  Can be set to **"ldap"**, **"internal"** or **"oidc"** 

* **ldap_group_dn** - (Optional) The distinguished name of the group within AD/LDAP 

`NOTE: oidc group type can only be used with harbor version v1.10.1 and above`