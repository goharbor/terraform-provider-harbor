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
