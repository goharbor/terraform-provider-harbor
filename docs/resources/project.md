# Resource: harbor_project

## Example Usage
```hcl
resource "haror_project" "main" {
    name                    = "main"
    public                  = false               # (Optional) Default value is false
    vulnerability_scanning  = true                # (Optional) Default vale is true. Automatically scan images on push 
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The of the project that will be created in harbor.

* **public** - (Optional) The project will be public accessibility. Can be set to **"true"** or **"false"** (Default: false)

* **vulnerability_scanning** - (Optional) Images will be scanned for vulnerabilities when push to harbor. Can be set to **"true"** or **"false"** (Default: true)

## Attributes Reference
In addition to all argument, the folloing attributes are exported:

* **project_id** - The id of the project with harbor.

## Import
Harbor project can be imported using the `project id` eg,

`
terraform import haror_project.main /projects/1
`