# Harbor Provider
The Harbor provider is used to configure an instance of Harbor. The provider needs to be configured with the proper credentials before it can be used.

## Resources
* [Resource: harbor_configuration](resources/configuration.md)
* [Resource: harbor_config_system](resources/config_system.md)
* [Resource: harbor_tasks](resources/tasks.md)
* [Resource: harbor_project](resources/project.md)
* [Resource: harbor_robot_account](resources/robot_account.md)

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

* **url** - (Required) The url of the harbor 
* **username** - (Required) The username to be used to access harbor
* **password** - (Required) The password to be used to access harbor
* **insecure** - (Optional) Choose to igorne certificate errors
* **api_version** - (Optional) Choose which version of the api you would like to use 1 or 2. Default is `2`
