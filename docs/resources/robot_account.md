# Resource: harbor_robot_account

## Example Usage
```hcl
resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "account" {
  name        = "${harbor_project.main.name}"
  description = "Robot account used to push images to harbor"
  project_id  = harbor_project.main.id
  actions      = ["push"]
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The of the project that will be created in harbor.

* **description** - (Optional) The description of the robot account will be displayed in harbor.

* **project_id** - (Required) The project id of the project that the robot account will be associated with.

* **actions** - (Optional) A list of actions that the robot account will be able to perform on the project. 
    You to have set `["pull"]` as minimal requirement, if `["push"]` is set you don't need to set pull. Other combinations can be `["push","create","read"]` or `["push","read"]` or `["pull","read"]`
    ```
    pull    = permission to pull from docker registry
    push    = permission to push to docker registry
    create  = permission to created helm charts
    read    = permission to read helm charts
    ```


## Attributes Reference
In addition to all argument, the following attributes are exported:

* **token** - The token of the robot account.