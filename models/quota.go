package models

type StorageQuota struct {
	Hard Hard `json:"hard"`
}
type Hard struct {
	Storage int64 `json:"storage"`
}
