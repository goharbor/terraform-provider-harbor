resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  primary_auth_mode  = true
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
