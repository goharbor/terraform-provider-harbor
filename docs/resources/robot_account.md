# Resource: harbor_robot_account

Harbor supports different level of robot accounts. Currently `system` and `project` level robot accounts are supported.

## System Level

Introduced in harbor 2.2.0, system level robot accounts can have basically all available permissions in harbor and are not dependent on a single project.

```hcl
resource "harbor_robot_account" "system" {
  name        = "example-system"
  description = "system level robot account"
  level       = "system"
  permissions {
    access {
      action   = "push"
      resource = "repository"
    }
    access {
      action   = "read"
      resource = "helm-chart"
    }
    access {
      action   = "read"
      resource = "helm-chart-version"
    }
    kind      = "project"
    namespace = "my-project-name"
  }
  permissions {
    access {
      action   = "pull"
      resource = "repository"
    }
    kind      = "project"
    namespace = "*"
  }
}
```

The above example, creates a system level robot account with permissions to
- pull repository across all projects
- push repository to project "my-project-name"
- read helm-chart and helm-chart-version in project "my-project-name"

## Project Level

Other than system level robot accounts, project level robot accounts can interact on project level only.

```hcl
resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "project" {
  name        = "example-project"
  description = "project level robot account"
  level       = "project"
  permissions {
    access {
      action   = "pull"
      resource = "repository"
    }
    access {
      action   = "push"
      resource = "repository"
    }
    kind      = "project"
    namespace = harbor_project.main.name
  }
}
```

The above example creates a project level robot account with permissions to
- pull repository on project "main"
- push repository on project "main"

For a full list of available actions and resources have a look at: https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go


## Argument Reference
The following arguments are supported:

* **name** - (string, required) The of the project that will be created in harbor.

* **level** - (string, required) Level of the robot account, currently either `system` or `project`.

* **description** - (string, optional) The description of the robot account will be displayed in harbor.

* **duration** - (int, optional) By default, the robot account will not expire. Set it to the amount of days until the account should expire.

* **disable** - (bool, optional) Disables the robot account when set to `true`.

* **permissions** - (block, required) Permissions to be applied to the robot account.
  ```
  permissions {
    access {
      action   = "action"   // eg. `push`, `pull`, `read`, etc.
      resource = "resource" // eg. `repository`, `helm-chart`, `read`, etc.
      effect   = "effect"   // either `allow` or `deny`
    }
    access {
      ...
    }
    kind      = "project"
    namespace = harbor_project.main.name
  }
  permissions {
    ...
  }
  ```
  **Note, that for `project` level accounts, only one `permission` block is allowed!**

  For a full list of available actions and resources have a look at: https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go



## Attributes Reference
In addition to all argument, the following attributes are exported:

* **token** - The token of the robot account.