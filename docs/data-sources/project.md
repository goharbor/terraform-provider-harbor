# Data Source: harbor_project

## Example Usage
```hcl
data "harbor_project" "main" {
    name    = "library" 
}

output "project_id" {
    value = data.harbor_project.main.id
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The name of the project.

## Attributes Reference
In addition to all argument, the following attributes are exported:

* **project_id** - The id of the project within harbor.

* **public** - If the project has public accessibility.

* **vulnerability_scanning** - If the images is scanned for vulnerabilities when push to harbor.

* **type** - The type of the project : Project or ProxyCache.
