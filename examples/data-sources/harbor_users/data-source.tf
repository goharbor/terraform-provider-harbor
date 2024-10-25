data "harbor_users" "example" {
  username = "example-user"
}

output "users_ids" {
  value = [data.harbor_users.example.users.*.id]
}
