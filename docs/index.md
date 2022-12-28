# Harbor Provider
The Harbor provider is used to configure an instance of Harbor. The provider needs to be configured with the proper credentials before it can be used.

## Resources
* [Resource: harbor_configuration](resources/configuration.md)
* [Resource: harbor_config_system](resources/config_system.md)
* [Resource: harbor_config_email](resources/config_email.md)
* [Resource: harbor_garbage_collection](resources/garbage_collection.md)
* [Resource: harbor_immutable_tag_rule](resources/immutable_tag_rule.md)
* [Resource: harbor_interrogation_services](resources/interrogation_services.md)
* [Resource: harbor_label](resources/label.md)
* [Resource: harbor_project_member_group](resources/project_member_group.md)
* [Resource: harbor_project_member_user](resources/project_member_user.md)
* [Resource: harbor_project](resources/project.md)
* [Resource: harbor_registry](resources/registry.md)
* [Resource: harbor_replication](resources/replication.md)
* [Resource: harbor_retention_policy](resources/retention_policy.md)
* [Resource: harbor_robot_account](resources/robot_account.md)
* [Resource: harbor_tasks](resources/tasks.md)
* [Resource: harbor_user](resources/user.md)

## Authentication
```hcl
provider "harbor" {
  url      = "https://harbor.aceme_corpartion.com"
  username = "insert_admin_username_here"
  password = "insert_password_here"
}
```
## Argument Reference
The following arguments are supported:

* **url** - (Required) The url of harbor
* **username** - (Required) The username to be used to access harbor
* **password** - (Required) The password to be used to access harbor
* **insecure** - (Optional) Choose to ignore certificate errors
* **api_version** - (Optional) Choose which version of the api you would like to use 1 or 2. Default is `2`
