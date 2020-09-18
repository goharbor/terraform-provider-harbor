# Resource: harbor_configuration

## Example Usage

```hcl
resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  oidc_name          = "azure"
  oidc_endpoint      = "https://login.microsoftonline.com/{GUID goes here}/v2.0"
  oidc_client_id     = "OIDC Client ID goes here"
  oidc_client_secret = "ODDC Client Secret goes here"
  oidc_scope         = "openid,email"
  oidc_verify_cert   = true
}
```

## Argument Reference
The following arguments are supported:

* **auth_mode** - (Requried) Harbor authenication mode. Can be **"oidc_auth"** or **"db_auth"**. (Default: **"db_auth"**)

* **oidc_name** - (Optional) The name of the oidc provider name. (Required - if auth_mode set to **oidc_auth**)

* **oidc_endpoint** - (Optional) The URL of an OIDC-complaint server. (Required - if auth_mode set to **oidc_auth**)

* **oidc_client_id** - (Optional) The client id for the odic server. (Required - if auth_mode set to **oidc_auth**)

* **oidc_client_serect** - (Optional) The client secert for the odic server. (Required - if auth_mode set to **oidc_auth**)

* **oidc_groups_claim** - (Optional) The name of the claim in the token whose values is the list of group names.

`NOTE: "oidc_groups_claim" can only be used with harbor version v1.10.1 and above`

* **oidc_scope** - (Optional) The scope sent to OIDC server during authentication. It has to contain “openid”. (Required - if auth_mode set to **oidc_auth**)

* **oidc_verify_cert** - (Optional) Set to **"false"** if your OIDC server is using a self-signed certificate. (Required - if auth_mode set to **oidc_auth**)