# Resource: harbor_retention_policy

## Example Usage

```hcl
resource "harbor_project" "main" {
  name                = "acctest"
}

resource "harbor_retention_policy" "main" {
	scope = harbor_project.main.id
	schedule = "Daily"
	rule {
		n_days_since_last_pull = 5
		repo_matching = "**"
		tag_matching = "latest"
	}
	rule {
		n_days_since_last_push = 10
		repo_matching = "**"
		tag_matching = "{latest,snapshot}"
	}

}
```

## Argument Reference
The following arguments are supported:

* `scope` - (Required) The project id of which you would like to apply this policy.

* `schedule` - (Optional) The schedule of when you would like the policy to run. This can be `Hourly`, `Daily`, `Weekly` or can be a custom cron string.

* `rule` - (Required) Al collection of rule blocks as documented below.

---
`rule` supports the following:
* `disabled`- (Optional) Specify if the rule is disable or not. Defaults to `false`
* `n_days_since_last_pull` - (Optional) retains the artifacts pulled within the lasts n days.
* `n_days_since_last_push` - (Optional) retains the artifacts pushed within the lasts n days.
* `most_recently_pulled` - (Optional) retain the most recently pulled n artifacts.
* `most_recently_pushed` - (Optional) retain the most recently pushed n artifacts.
* `always_retain` - (Optional) retain always.
* `repo_matching`- (Optional) For the repositories matching.
* `repo_excluding` - (Optional) For the repositories excuding.
* `tag_matching`- (Optional) For the tag matching.
* `tag_excluding` - (Optional) For the tag excuding.
* `untagged_artifacts`- (Optional) with untagged artifacts. Defaults to `true`

~> Multiple tags or repositories must be provided as a comma-separated list wrapped into curly brackets `{ }`. Otherwise, the value is interpreted as a single value.

---

## Import
Harbor retention policy can be imported using the `retention_policy id` eg,

`
terraform import harbor_retention_policy.main /retentions/10
`
