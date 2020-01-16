# Resource: harbor_robot_account

## Example Usage
```
resource "haror_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "account" {
  name        = "${harbor_project.main.name}"
  description = "Robot account used to push images to harbor"
  project_id  = harbor_project.main.project_id
  action      = "push"
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The of the project that will be created in harbor.

* **description** - (Optional) The description of the robot account will be displayed in harbor.

* **project_id** - (Required) The project id of the project that the robot account will be associated with.

* **action** - (Optional) The action that the robot account will be able to perform on the project. Can be **"pull"** or **"push"** (Default: **pull**).

## Attributes Reference
In addition to all argument, the folloing attributes are exported:

* **token** - The token of the robot account.