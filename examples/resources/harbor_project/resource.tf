resource "harbor_project" "main" {
  name                        = "main"
  public                      = false               # (Optional) Default value is false
  vulnerability_scanning      = true                # (Optional) Default value is true. Automatically scan images on push
  enable_content_trust        = true                # (Optional) Default value is false. Deny unsigned images from being pulled (notary)
  enable_content_trust_cosign = false               # (Optional) Default value is false. Deny unsigned images from being pulled (cosign)
}
