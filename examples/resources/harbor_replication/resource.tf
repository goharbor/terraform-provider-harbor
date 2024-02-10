resource "harbor_registry" "main" {
  provider_name = "docker-hub"
  name          = "test_docker_harbor"
  endpoint_url  = "https://hub.docker.com"

}

resource "harbor_replication" "push" {
  name        = "test_push"
  action      = "push"
  registry_id = harbor_registry.main.registry_id
}

resource "harbor_replication" "alpine" {
  name        = "alpine"
  action      = "pull"
  registry_id = harbor_registry.main.registry_id
  schedule = "0 0/15 * * * *"
  filters {
    name = "library/alpine"
  }
  filters {
    tag = "3.*.*"
  }
  filters {
    resource = "artifact"
  }
  filters {
    labels = ["qa"]
  }
}

resource "harbor_replication" "alpine" {
  name        = "alpine"
  action      = "pull"
  registry_id = harbor_registry.main.registry_id
  schedule = "event_based"
  filters {
    name = "library/alpine"
  }
  filters {
    tag = "3.*.*"
  }
}
