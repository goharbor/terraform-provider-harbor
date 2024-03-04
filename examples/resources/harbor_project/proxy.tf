resource "harbor_project" "main" {
  name        = "acctest"
  registry_id = harbor_registry.docker.registry_id
}

resource "harbor_registry" "docker" {
  provider_name = "docker-hub"
  name          = "test"
  endpoint_url  = "https://hub.docker.com"
}
