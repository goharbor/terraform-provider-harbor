# Resource: harbor_replication



## Example Usage

```hcl
  resource "harbor_registry" "main" {
	provider_name = "docker-hub"
	name = "docker-hub-replication"
	endpoint_url = "https://hub.docker.com"
  }

  resource "harbor_replication" "push" {
	name  = "test_push"
	action = "push"
	registry_id = harbor_registry.main.registry_id
  }
```

## Argument Reference
The following arguments are supported:

* **name** - (Required)

* **action** - (Required)

* **schedule** - (Optional) The scheduled time of when the containter register will be push / pull. In cron base format. Hourly `"0 0 * * * *"`, Daily `"0 0 0 * * *"`, Monthly `"0 0 0 * * 0"`
* **override** - (Optional) Specify whether to override the resources at the destination if a resources with the same name exist. Can be set to `true` or `false` (Default: `true`)
* **enabled** - (Optional) Specify whether the replication is enabled. Can be set to `true` or `false` (Default: `true`)
* **description** (Optional) Write about decscription of the replication policy.

* **filters** - (Optional) A collection of `filters` block as documented below.

---

**filters** supports the following:

* **name** - (Optional) Filter on the name of the resource.
* **tag** - (Optional) Filter on the tag/verison of the resource.
* **labels** - (Optional) Filter on the resource according to labels.
* **resource** - (Optional) Filter on the recource type. Can be one of the following types. `image`,`chart`, `artifact`
				


## Attributes Reference
In addition to all argument, the folloing attributes are exported:

* **replication_policy_id**
  
## Import
Harbor project can be imported using the `replication id` eg,

`
terraform import haror_project.main /registries/7
`