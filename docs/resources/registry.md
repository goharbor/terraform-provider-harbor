# Resource: harbor_registry



## Example Usage

```hcl
resource "harbor_registry" "main" {
  provider_name = "docker-hub"
  name          = "test_docker_harbor"
  endpoint_url  = "https://hub.docker.com"
}
```

## Argument Reference
The following arguments are supported:

* **provider_name** - (Required) The name of the provider type. Supported values include `alibaba`, `artifact-hub`, `aws`, `azure`, `docker-hub`, `docker-registry`, `gitlab`, `github`, `google`, `harbor`, `helm`, `huawei`, `jfrog`

* **name** - (Required) The name of the register.

* **endpoint_url** - (Required) The url endpoint for the external container register ie, `https://hub.docker.com`

* **description** - (Optional) The description of the external container register.

* **access_id** - (Optional) The username / access id for the external container register 

* **access_secret** - (Optional) The password / access keys / token for the external container register

* **insecure** - (Optional) Verifies the certificate of the external container register. Can be set to **"true"** or **"false"** (Default: false)

* **enabled** - (Optional) enables / disables the external container register within harbor. Can be set to **"true"** or **"false"** (Default: true)

## Attributes Reference
In addition to all argument, the following attributes are exported:

* **registry_id** - The id of the register within harbor.

* **status** - The health status of the external container register

## Import
Harbor project can be imported using the `registry id` eg,

`
terraform import harbor_registry.main /registries/7
`
