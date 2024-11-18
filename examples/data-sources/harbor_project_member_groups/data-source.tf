data "harbor_project_member_groups" "example" {
  project_id = "1"
}

output "project_member_group_ids" {
  value = [data.harbor_project_member_groups.example.project_member_groups.*.id]
}
