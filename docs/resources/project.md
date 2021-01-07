# Resource: harbor_project

## Example Usage
```hcl
resource "harbor_project" "main" {
    name                    = "main"
    public                  = false               # (Optional) Default value is false
    vulnerability_scanning  = true                # (Optional) Default vale is true. Automatically scan images on push 
}
```

## Harbor project example as proxy cache
```hcl
resource "harbor_project" "main" {
  name        = "acctest"
  registry_id = harbor_registry.docker.registry_id
}

resource "harbor_registry" "docker" {
  provider_name = "docker-hub"
  name          = "test"
  endpoint_url  = "https://hub.docker.com"
}
```


## Argument Reference
The following arguments are supported:

* `name` - (Required) The of the project that will be created in harbor.

* `public` - (Optional) The project will be public accessibility. Can be set to `"true"` or `"false"` (Default: false)

* `vulnerability_scanning` - (Optional) Images will be scanned for vulnerabilities when push to harbor. Can be set to `"true"` or `"false"` (Default: true)

* `registry_id` - (Optional) To enabled project as Proxy Cache

## Attributes Reference
In addition to all argument, the following attributes are exported:

* `project_id` - The id of the project with harbor.

## Import
Harbor project can be imported using the `project id` eg,

`
terraform import harbor_project.main /projects/1
`
