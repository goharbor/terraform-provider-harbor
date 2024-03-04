resource "harbor_config_auth" "ldap" {
  auth_mode            = "ldap_auth"
  primary_auth_mode    = true
  ldap_url             = "openldap.default.svc.cluster.local:389"
  ldap_search_dn       = "cn=admin,dc=example,dc=org"
  ldap_search_password = "Not@SecurePassw0rd"
  ldap_base_dn         = "dc=example,dc=org"
  ldap_uid             = "email"
  ldap_verify_cert     = false
}
