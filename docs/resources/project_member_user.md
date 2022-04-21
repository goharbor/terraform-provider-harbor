# Resource: harbor_project_member_user

## Example Usage
```hcl
resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_project_member_user" "main" {
  project_id    = harbor_project.main.id
  user_name     = "testing1"
  role          = "projectadmin"
}

```

## Argument Reference
The following arguments are supported:

* **user_name** - (Required) The name of the member entity

* **project_id** - (Required) The project id of the project that the entity will have access to.

* **role** - (Required) The premissions that the entity will be granted.

## Import
Harbor project member user can be imported using the `project id` and `member id` eg,

`
terraform import harbor_project_member_user.main /projects/10/members/200
`
