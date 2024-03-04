data "harbor_registry" "main" {
  name          = "test_docker_harbor"
}

output "harbor_registry_id" {
  value   = data.harbor_registry.main.id
}
