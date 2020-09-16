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
* **override** - (Optional)
* **enabled** - (Optional)
* **description** (Optional)	

"filters":
				
						"name":
						"tag": 
						"labels"ptional: true,
			
						"resource"


* **provider_name** - (Required) The name of the provider type. Supported values include `alibaba`, `aws`, `azure`, `docker-hub`, `docker-registry`, `gitlab`, `google`, `harbor`, `helm`, `huawei`, `jfrog`

* **name** - (Required) The name of the register.

* **endpoint_url** - (Required) The url endpoint for the external container register ie, `https://hub.docker.com`

* **description** - (Optional) The description of the external container register.

* **access_id** - (Optional) The username / access id for the external container register 

* **access_key** - (Optional) The password / access keys / token for the external container register

* **insecure** - (Optional) Verifies the certificate of the external container register. Can be set to **"true"** or **"false"** (Default: true)

* **enabled** - (Optional) enables / disables the external container register within harbor. Can be set to **"true"** or **"false"** (Default: true)

## Attributes Reference
In addition to all argument, the folloing attributes are exported:

* **replication_policy_id**
  
## Import
Harbor project can be imported using the `registry id` eg,

`
terraform import haror_project.main /registries/7
`