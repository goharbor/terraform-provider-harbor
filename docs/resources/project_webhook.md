# Resource: harbor_project_webhook

## Example Usage
```hcl
resource "harbor_project_webhook" "main" {
  name        = "test_webhook"
  address     = "https://webhook.domain.com"
  project_id  = harbor_project.main.id
  notify_type = "http"

  events_types = [
    "DELETE_ARTIFACT",
    "PULL_ARTIFACT",
    "PUSH_ARTIFACT",
    "QUOTA_EXCEED",
    "QUOTA_WARNING",
    "REPLICATION",
    "SCANNING_FAILED",
    "SCANNING_COMPLETED",
    "TAG_RETENTION"
  ]

}

resource "harbor_project" "main" {
  name = "test-project"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required, string) The name of the webhook that will be created in harbor.
* `address` - (Required, string) The address of the webhook
* `description` _ (Optional, string) A description of the webhook
* `enabled` - (Optional, bool), To enable / disable the webhook. Default `true` 
* `project_id` - (Required, string) The project id (**/projects/ID**) of the harbor that webhook related to.
* `notify_type` - (Required, string) The notification type either `http` or `slack`
* `events_types` - (Required, list(string)) The type events you want to subscript to can be `DELETE_ARTIFACT`, `PULL_ARTIFACT`, `PUSH_ARTIFACT`, `QUOTA_EXCEED`, `QUOTA_WARNING`, `REPLICATION`, `SCANNING_FAILED`, `SCANNING_COMPLETED`, `TAG_RETENTION`
* `auth_header` - (Required, string) authentication header for you the webhook
* `skip_cert_verify` - (Optional - bool) checks the for validate SSL certificate.