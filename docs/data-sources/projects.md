# Data Source: harbor_projects

## Example Usage
```hcl
data "harbor_projects" "proxycache" {
    type = "ProxyCache"
}

output "proxy_cache_projects" {
    value = data.harbor_projects.proxycache
}
```

## Argument Reference
The following arguments are supported:

* **name** - (Optional) The name of the project.

* **type** - (Optional) The type of the project : Project or ProxyCache.

* **public** - (Optional) If the project has public accessibility.

* **vulnerability_scanning** - (Optional) If the images will be scanned for vulnerabilities when push to harbor.

## Attributes Reference
In addition to all argument, the following attributes are exported:

* **projects** - A list of projects matching previous arguments. Each **project** object provides the attributes documented below.

---

**project** object exports the following:

* **name** - The name of the project.

* **type** - The type of the project : Project or ProxyCache.

* **public** - If the project has public accessibility.

* **vulnerability_scanning** - If the images will be scanned for vulnerabilities when push to harbor.
