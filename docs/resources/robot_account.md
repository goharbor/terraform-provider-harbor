# Resource: harbor_robot_account

Harbor supports different levels of robot accounts. Currently `system` and `project` level robot accounts are supported.

## Example Usage

### System Level
Introduced in harbor 2.2.0, system level robot accounts can have basically [all available permissions](https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go) in harbor and are not dependent on a single project.

```hcl
resource "harbor_project" "main" {
    name = "main"
}

resource "harbor_robot_account" "system" {
  name        = "example-system"
  description = "system level robot account"
  level       = "system"
  permissions {
    access {
      action   = "create"
      resource = "labels"
    }
    kind      = "system"
    namespace = "/"
  }
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
    namespace = harbor_project.main.name
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
- permission to create labels on system level
- pull repository across all projects
- push repository to project "my-project-name"
- read helm-chart and helm-chart-version in project "my-project-name"

### Project Level

Other than system level robot accounts, project level robot accounts can interact on project level only.
The [available permissions](https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go) are mostly the same as for system level robots.


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


## Argument Reference
The following arguments are supported:

* **name** - (string, required) The of the project that will be created in harbor.

* **level** - (string, required) Level of the robot account, currently either `system` or `project`.

* **description** - (string, optional) The description of the robot account will be displayed in harbor.

* **duration** - (int, optional) By default, the robot account will not expire. Set it to the amount of days until the account should expire.

* **disable** - (bool, optional) Disables the robot account when set to `true`.

* **permissions** - (block, required) [Permissions](#permissions-arguments) to be applied to the robot account. 
  ```
  permissions {
    access {
      action   = "action"
      resource = "resource"
      effect   = "effect"
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

### Permissions Arguments
* **access** - (block, required) Define one or multiple [access blocks](#access-arguments).

* **kind** - (string, required) Either `system` or `project`.

* **namespace** - (string, required) namespace is the name of your project.
                  For kind `system` permissions, always use `/` as namespace.
                  Use `*` to match all projects.

* **secret** - (string, optional) The secret of the robot account used for authentication. Defaults to random generated string from Harbor 
  
### Access Arguments
* **action** - (string, required) Eg. `push`, `pull`, `read`, etc. Check [available actions](https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go).

* **resource** - (string, required) Eg. `repository`, `helm-chart`, `labels`, etc. Check [available resources](https://github.com/goharbor/harbor/blob/master/src/common/rbac/const.go).

* **effect** - (string, optional) Either `allow` or `deny`. Defaults to `allow`.


## Attributes Reference
In addition to all argument, the following attributes are exported:


* **full_name** - Full name of the robot account which harbor generates including the robot prefix. Eg. `robot$project+name` or `harbor@project+name` (depending on your robot prefix).