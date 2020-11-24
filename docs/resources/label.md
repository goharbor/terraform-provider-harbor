# Resource: harbor_label

## Example Usage

* Create a global label within harbor
```hcl
	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FF0000"
		description 	= "Description to for acceptance test"
	}
```

* Creates a label for project 
```hcl
	resource "harbor_project" "main" {
		name = "acctest"
	}

	resource "harbor_label" "main" {
		name  		= "accTest"
		color 		= "#FFFFFF"
		description = "Description for acceptance test"
		project_id	= harbor_project.main.id
	}
```

## Argument Reference
The following arguments are supported:

* **name** - (Required) The of name of the label within harbor.

* **color** - (Optional) The color of the label within harbor (Default: #FFFFF)

* **description** - (Optional) The Description of the label within harbor

* **project_id** - The id of the project with harbor.

## Import
Harbor label can be imported using the `label id` eg,

`
terraform import harbor_label.main /labels/1
`
