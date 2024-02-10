resource "harbor_project" "main" {
  name = "acctest"
}

resource "harbor_retention_policy" "main" {
  scope = harbor_project.main.id
  schedule = "Daily"
  rule {
    n_days_since_last_pull = 5
    repo_matching = "**"
    tag_matching = "latest"
  }
  rule {
    n_days_since_last_push = 10
    repo_matching = "**"
    tag_matching = "{latest,snapshot}"
  }
}
