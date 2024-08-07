---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "harbor_config_system Resource - terraform-provider-harbor"
subcategory: ""
description: |-
  
---

# harbor_config_system (Resource)

<!-- schema generated by tfplugindocs -->

## Example Usage

```terraform
resource "harbor_config_system" "main" {
  project_creation_restriction = "adminonly"
  robot_token_expiration       = 30
  robot_name_prefix            = "harbor@"
  storage_per_project          = 100
}
```

## Schema

### Optional

- `project_creation_restriction` (String) Who can create projects within Harbor. Can be `"adminonly"` or `"everyone"`
- `read_only` (Boolean) Whether or not the system is in read only mode.
- `robot_name_prefix` (String) Robot account prefix.
- `robot_token_expiration` (Number) The amount of time in days a robot account will expire.
- `scanner_skip_update_pulltime` (Boolean) Whether or not to skip update pull time for scanner.
- `storage_per_project` (Number) Default quota space per project in GIB. Default is -1 (unlimited).
- `audit_log_forward_endpoint` (String) The endpoint to forward audit logs to.
- `skip_audit_log_database` (Boolean) Whether or not to skip audit log database.
- `banner_message` (Block Set) (see [below for nested schema](#nestedblock--banner_message))

<a id="nestedblock--banner_message"></a>

### Nested Schema for `banner_message`

### Required
- `message` (String) The message to display in the banner.

### Optional
- `closable` (Boolean) Whether or not the banner message is closable.
- `type` (String) The type of banner message. Can be `"info"`, `"warning"`, `"success"` or `"danger"`.
- `from_date` (String) The date the banner message will start displaying. (Format: `MM/DD/YYYY`)
- `to_date` (String) The date the banner message will stop displaying. (Format: `MM/DD/YYYY`)

#### Notes
`scanner_skip_update_pulltime` can only be used with harbor version v2.8.0 and above

`robot_token_expiration` if the time is set to high you will get a 500 internal server error message when creating robot accounts

### Read-Only

- `id` (String) The ID of this resource.
