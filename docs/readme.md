# Harbor Provider
The Harbor provider is used to configure an instance of Harbor. The provider needs to be configured with the proper credentials before it can be used.

## Resources
* [Resource: harbor_configuration](harbor_configuration.md)
* [Resource: harbor_config_system](harbor_config_system.md)
* [Resource: harbor_tasks](harbor_tasks.md)
* [Resource: harbor_project](harbor_project.md)
* [Resource: harbor_robot_account](harbor_robot_account.md)
* [Resource: harbor_project_member](harbor_project_member.md)

## Authentication
```
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
* **api_version** - (Optional) Chosse which version of the api you would like to use 1 or 2.
