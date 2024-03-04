resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "project" {
  name        = "example-project"
  description = "project level robot account"
  level       = "project"
  permissions {
    access {
      action   = "pull"
      resource = "repository"
    }
    access {
      action   = "push"
      resource = "repository"
    }
    kind      = "project"
    namespace = harbor_project.main.name
  }
}
