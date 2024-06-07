resource "harbor_preheat_instance" "example" {
  name     = "example-preheat-instance"
  vendor   = "dragonfly"
  endpoint = "http://example.com"
}
