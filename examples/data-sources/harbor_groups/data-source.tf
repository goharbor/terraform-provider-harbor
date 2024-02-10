data "harbor_groups" "example" {
  group_name = "example-group"
}

output "group_ids" {
  value = [data.harbor_groups.example.*.id]
}
