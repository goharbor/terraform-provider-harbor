resource "random_password" "password" {
  length  = 12
  special = false
}

resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "system" {
  name        = "example-system"
  description = "system level robot account"
  level       = "system"
  secret      = resource.random_password.password.result
  permissions {
    access {
      action   = "create"
      resource = "labels"
    }
    kind      = "system"
    namespace = "/"
  }
  permissions {
    access {
      action   = "push"
      resource = "repository"
    }
    kind      = "project"
    namespace = harbor_project.main.name
  }
  permissions {
    access {
      action   = "pull"
      resource = "repository"
    }
    kind      = "project"
    namespace = "*"
  }
}
