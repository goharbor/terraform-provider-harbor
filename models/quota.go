package models

import "time"

type StorageQuota struct {
	Hard Hard `json:"hard"`
}
type Hard struct {
	Storage int64 `json:"storage"`
}
type Used struct {
	Storage int64 `json:"storage"`
}
type QuotaResponse struct {
	CreationTime	time.Time	`json:"creation_time"`
	Hard		Hard		`json:"hard"`
	ID		int		`json:"id"`
	Ref struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Owner string `json:"owner"`
        } `json:"ref"`
	UpdateTime	time.Time	`json:"update_time"`
	Used		Used		`json:"used"`
}
