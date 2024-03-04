resource "harbor_group" "storage-group" {
  group_name = "storage-group"
  group_type = 3
  ldap_group_dn = "storage-group"
}
