# Resource: harbor_immutable_tag_rule

## Example Usage

```hcl
resource "harbor_project" "main" {
  name                = "acctest"
}

resource "harbor_immutable_tag_rule" "main" {
	project_id = harbor_project.main.id
	repo_matching = "**"
	tag_excluding = "latest"
}
```

## Argument Reference
The following arguments are supported:

* `disabled`- (Optional) Specify if the rule is disable or not. Defaults to `false`
* `repo_matching`- (Optional) For the repositories matching.
* `repo_excluding` - (Optional) For the repositories excuding.
* `tag_matching`- (Optional) For the tag matching.
* `tag_excluding` - (Optional) For the tag excuding.
* `scope` - (Required) The project id of which you would like to apply this policy.

## Import
Harbor immutable tag rule can be imported using the `project and immutabletagrule ids` eg,

`
terraform import harbor_immutable_tag_rule.main /projects/4/immutabletagrules/25
`
