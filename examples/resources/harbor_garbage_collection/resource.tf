resource "harbor_garbage_collection" "main" {
  schedule         = "Daily"
  delete_untagged  = true
}
