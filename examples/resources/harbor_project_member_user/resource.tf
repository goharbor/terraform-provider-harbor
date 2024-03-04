resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_project_member_user" "main" {
  project_id    = harbor_project.main.id
  user_name     = "testing1"
  role          = "projectadmin"
}
