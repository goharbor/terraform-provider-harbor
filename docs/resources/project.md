# Resource: harbor_project

## Example Usage
```hcl
resource "harbor_project" "main" {
    name                    = "main"
    public                  = false               # (Optional) Default value is false
    vulnerability_scanning  = true                # (Optional) Default vale is true. Automatically scan images on push
    enable_content_trust    = true                # (Optional) Default vale is false. Deny unsigned images from being pulled
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

* `name` - (Required) The name of the project that will be created in harbor.

* `public` - (Optional) The project will be public accessibility. Can be set to `"true"` or `"false"` (Default: false)

* `vulnerability_scanning` - (Optional) Images will be scanned for vulnerabilities when push to harbor. Can be set to `"true"` or `"false"` (Default: true)

* `deployment_security` - (Optional) Prevent deployment of images with vulnerability severity equal or higher than the specified value. Images must be scanned before this takes effect. Possible values: `critical`, `high`, `medium`, `low`, `none`. (Default: `""` - empty)

* `registry_id` - (Optional) To enabled project as Proxy Cache

* `storage_quota` - (Optional) The storage quota of the project in GB's

* `enable_content_trust` - (Optional) Enables Content Trust for project. When enabled it queries the embedded docker notary server. Can be set to `"true"` or `"false"` (Default: false)

* `force_destroy` - (Optional, Default: `false`) A boolean that indicates all repositories should be deleted from the project so that the project can be destroyed without error. These repositories are *not* recoverable.

* `cve_allowlist` - (Optional) Project allowlist allows vulnerabilities in this list to be ignored in this project when pushing and pulling images. Should be in the format or `["CVE-123", "CVE-145"]` or `["CVE-123"]`

## Attributes Reference
In addition to all argument, the following attributes are exported:

* `project_id` - The id of the project with harbor.
* `cve_allowlist` - A list of allowed CVE IDs.

## Import
Harbor project can be imported using the `project id` eg,

`
terraform import harbor_project.main /projects/1
`
