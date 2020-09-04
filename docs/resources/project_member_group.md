# Resource: harbor_project_member_group

## Example Usage
```
resource "haror_project" "main" {
    name = "main"
}

resource "harbor_project_member_group" "main" {
  project_id    = harbor_project.main.id
  name          = "testing1"
  role          = "master"
  type          = "oidc"
}

```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The name of the member entity

* **project_id** - (Required) The project id of the project that the entity will have access to.

* **role** - (Required) The premissions that the entity will be granted.

* **type** - (Requried) The group type.  Can be set to **"ldap"**, **"internal"** or **"oidc"** 

`NOTE: odic group type can only be used with harbor version v1.10.1 and above`