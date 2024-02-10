resource "harbor_registry" "main" {
  provider_name = "docker-hub"
  name          = "test_docker_harbor"
  endpoint_url  = "https://hub.docker.com"
}
