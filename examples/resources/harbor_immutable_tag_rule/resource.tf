resource "harbor_project" "main" {
  name = "acctest"
}

resource "harbor_immutable_tag_rule" "main" {
  disabled = true
  project_id = harbor_project.main.id
  repo_matching = "**"
  tag_matching = "v2.*"
  tag_excluding = "latest"
}
