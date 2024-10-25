data "harbor_robot_accounts" "example" {
  name = "example-robot"
}

output "robot_account_ids" {
  value = [data.harbor_robot_accounts.example.robot_accounts.*.id]
}
