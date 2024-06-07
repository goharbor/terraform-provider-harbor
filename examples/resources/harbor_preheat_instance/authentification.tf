resource "harbor_preheat_instance" "example" {
  name       = "example-preheat-instance"
  vendor     = "dragonfly"
  endpoint   = "http://example.com"
  auth_mode  = "BASIC"
  username   = "example-user"
  password   = "example-password"
}
