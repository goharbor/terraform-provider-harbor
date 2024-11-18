data "harbor_project_member_users" "example" {
  project_id = "1"
}

output "project_member_user_ids" {
  value = [data.harbor_project_member_users.example.project_member_users.*.id]
}
