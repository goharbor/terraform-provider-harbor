data "harbor_project" "main" {
    name    = "library" 
}

output "project_id" {
    value = data.harbor_project.main.id
}
