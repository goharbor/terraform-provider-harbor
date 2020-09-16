# Resource: harbor_config_system

## Example Usage

```hcl
resource "harbor_config_system" "main" {
  project_creation_restriction = "adminonly"
  robot_token_expiration       = 30
}
```

## Argument Reference
The following arguments are supported:

* **project_creation_restriction** - (Optional) Who can create projects within Harbor. Can be **"adminonly"** or **"everyone"**

* **robot_token_expiration** - (Optional) The amount of time in days a robot account will expiry. 

`NOTE: If the time is set to high you will get a 500 internal server error message when creating robot accounts`