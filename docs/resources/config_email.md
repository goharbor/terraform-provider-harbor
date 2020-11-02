# Resource: harbor_config_system

## Example Usage

```hcl
resource "harbor_config_email" "main" {
	email_host = "server.acme.com"
	email_from = "dont_reply@acme.com"
}
```

## Argument Reference
The following arguments are supported:

* **email_host** - (Required) The FQDN of the email server
* **email_port** - (Optional) The smtp port for the email server `Default: 25`
* **email_username** - (Optional) The username for the email server
* **email_password** - (Optional) The password for the email server
* **email_from** - (Required) - The email from address ie, `dont_reply@acme.com` 
* **email_ssl** - (Optional) Enable SSL for email server connection