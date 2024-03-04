data "harbor_projects" "proxycache" {
    type = "ProxyCache"
}

output "proxy_cache_projects" {
    value = data.harbor_projects.proxycache
}
