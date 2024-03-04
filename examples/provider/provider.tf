provider "harbor" {
  url      = "https://harbor.aceme_corpartion.com"
  username = "insert_admin_username_here"
  password = "insert_password_here"
  bearer_token = "insert_bearer_token_here"
  insecure = true
  api_version = 2
}
