resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_project_member_group" "main" {
  project_id    = harbor_project.main.id
  group_name    = "testing1"
  role          = "projectadmin"
  type          = "oidc"
}
