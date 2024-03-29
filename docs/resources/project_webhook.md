---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "harbor_project_webhook Resource - terraform-provider-harbor"
subcategory: ""
description: |-
  
---

# harbor_project_webhook (Resource)

<!-- schema generated by tfplugindocs -->

## Example Usage

```terraform
resource "harbor_project" "main" {
  name = "test-project"
}

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
```

## Schema

### Required

- `address` (String) The address of the webhook.
- `events_types` (List of String) The type events you want to subscript to can be 
- `name` (String) The name of the webhook that will be created in harbor.
- `notify_type` (String) The notification type either `http` or `slack`.
- `project_id` (String) The project id of the harbor that webhook related to.

##### `events_types` Options

- `DELETE_ARTIFACT`
- `PULL_ARTIFACT`
- `PUSH_ARTIFACT`
- `QUOTA_EXCEED`
- `QUOTA_WARNING`
- `REPLICATION`
- `SCANNING_FAILED`
- `SCANNING_COMPLETED`
- `TAG_RETENTION`

### Optional

- `auth_header` (String) authentication header for you the webhook.
- `description` (String) A description of the webhook.
- `enabled` (Boolean) To enable / disable the webhook. Default `true`.
- `skip_cert_verify` (Boolean) checks the for validate SSL certificate.

### Read-Only

- `id` (String) The ID of this resource.
