# Data Source: harbor_registry

## Example Usage

```hcl
data "harbor_registry" "main" {
  name          = "test_docker_harbor"
}

output "harbor_registry_id" {
  value   = data.harbor_registry.main.id
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The name of the register.

## Attributes Reference
In addition to all argument, the folloing attributes are exported:

* **registry_id** - The id of the register within harbor.

* **status** - The health status of the external container register

* **endpoint_url** - The url endpoint for the external container register
  
* **description** - The description of the external container register.

* **insecure** - If the certificate of the external container register can be verified.

* **type** - The type of the provider type.