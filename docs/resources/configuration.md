# Resource: harbor_configuration

## Example Usage
How to configure oidc
```hcl
resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  oidc_name          = "azure"
  oidc_endpoint      = "https://login.microsoftonline.com/{GUID goes here}/v2.0"
  oidc_client_id     = "OIDC Client ID goes here"
  oidc_client_secret = "ODDC Client Secret goes here"
  oidc_scope         = "openid,email"
  oidc_verify_cert   = true
  oidc_auto_onboard  = true
  oidc_user_claim    = "name"
  oidc_admin_group   = "administrators"
}
```

How to configure ldap
```hcl
resource "harbor_config_auth" "ldap" {
  auth_mode            = "ldap_auth"
  ldap_url             = "openldap.default.svc.cluster.local:389"
  ldap_search_dn       = "cn=admin,dc=example,dc=org"
  ldap_search_password = "Not@SecurePassw0rd"
  ldap_base_dn         = "dc=example,dc=org"
  ldap_uid             = "email"
  ldap_verify_cert     = "false"
}
```

## Argument Reference
The following arguments are supported:

* `auth_mode` - (Required) Harbor authentication mode. Can be `"oidc_auth"`, `"db_auth"` or `"ldap_auth"`. (Default: **"db_auth"**)

* `oidc_name` - (Optional) The name of the oidc provider name. (Required - if auth_mode set to **oidc_auth**)

* `oidc_endpoint` - (Optional) The URL of an OIDC-complaint server. (Required - if auth_mode set to **oidc_auth**)

* `oidc_client_id` - (Optional) The client id for the oidc server. (Required - if auth_mode set to **oidc_auth**)

* `oidc_client_serect` - (Optional) The client secert for the oidc server. (Required - if auth_mode set to **oidc_auth**)

* `oidc_groups_claim` - (Optional) The name of the claim in the token whose values is the list of group names.

`NOTE: "oidc_groups_claim" can only be used with harbor version v1.10.1 and above`

* `oidc_scope` - (Optional) The scope sent to OIDC server during authentication. It has to contain “openid”. (Required - if auth_mode set to **oidc_auth**)

* `oidc_verify_cert` - (Optional) Set to **"false"** if your OIDC server is using a self-signed certificate. (Required - if auth_mode set to **oidc_auth**)

* `oidc_auto_onboard` - (Optional) Default is **"false"**, set to **"true"** if you want to skip the user onboarding screen, so user cannot change its username

* `oidc_user_claim` - (Optional) Default is **"name"**
  
* `oidc_admin_group` - (Optional) All members of this group get Harbor admin permissions.


* `ldap_url` - (Optional) The ldap server. Required when auth_mode is set to ldap.
* `ldap_base_dn` - (Optional) A user's DN who has the permission to search the LDAP/AD server. 
* `ldap_uid`- (Optional) The attribute used in a search to match a user. It could be uid, cn, email, sAMAccountName or other attributes depending on your LDAP/AD.
* `ldap_verify_cert`- (Optional) Verify Cert from LDAP Server.
* `ldap_search_dn` - (Optional) The base DN from which to look up a user in LDAP/AD.
* `ldap_search_password` - (Optional) The password for the user that will perform the LDAP search
* `ldap_filter` - (Optional) ldap filters
* `ldap_scope` - (Optional) LDAP Group Scope
* `ldap_group_base_dn` - (Optional) The base DN from which to look up a group in LDAP/AD.
* `ldap_group_filter` - (Optional) The filter to look up an LDAP/AD group.
* `ldap_group_gid` - (Optional) - The attribute used in a search to match a user
* `ldap_group_admin_dn` - (Optional) Specify an LDAP group DN. All LDAP user in this group will have harbor admin privilege
* `ldap_group_membership` - (Optional) The attribute indicates the membership of LDAP group
* `ldap_group_scope` - (Optional) The scope to search for groups
			
