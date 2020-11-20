# Resource: harbor_user

## Example Usage
```hcl

resource "harbor_user" "main" {
  username = "john"
  password = "Password12345!"
  full_name = "John Smith"
  email = "john@smith.com"
}
```

## Argument Reference
The following arguments are supported:

* **username** - (Required) The username of the internal user.

* **password** - (Required) The password for the interal user.

* **full_name** - (Required) The Full Name of the internal user.

* **email** - (Required) The email address of the internal user.

* **comment** - (Optional) Any comments for that are need for the internal user.

* **admin** - (Optional) If the user will have admin rights within Harbor (Default: **false**)

## Import
An internal user harbor user can be imported using the `user id` eg,

`
terraform import harbor_user.main /users/19
`