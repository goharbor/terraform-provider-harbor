# Resource: harbor_config_system

## Example Usage

```hcl
resource "harbor_config_system" "main" {
  project_creation_restriction = "adminonly"
  robot_token_expiration       = 30
  robot_name_prefix            = "harbor@"
}
```

## Argument Reference
The following arguments are supported:

* **project_creation_restriction** - (Optional) Who can create projects within Harbor. Can be **"adminonly"** or **"everyone"**

* **robot_token_expiration** - (Optional) The amount of time in days a robot account will expiry. 

* **robot_name_prefix** - (Optional) Robot account prefix.
`NOTE: If the time is set to high you will get a 500 internal server error message when creating robot accounts`

* **scanner_skip_update_pulltime** - (Optional) Whether or not to skip update pull time for scanner.
`NOTE: "scanner_skip_update_pulltime" can only be used with harbor version v2.8.0 and above`