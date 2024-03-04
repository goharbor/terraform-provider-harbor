resource "harbor_config_system" "main" {
  project_creation_restriction = "adminonly"
  robot_token_expiration       = 30
  robot_name_prefix            = "harbor@"
  storage_per_project          = 100
}
