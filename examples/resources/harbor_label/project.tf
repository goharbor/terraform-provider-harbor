resource "harbor_project" "main" {
  name = "acctest"
}

resource "harbor_label" "main" {
  name      = "accTest"
  color     = "#FFFFFF"
  description = "Description for acceptance test"
  project_id  = harbor_project.main.id
}
