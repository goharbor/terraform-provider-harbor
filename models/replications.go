package models

var PathReplication = "/replication/policies"
var PathExecution = "/replication/executions"

type ReplicationBody struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id"`
	SrcRegistry struct {
		ID int `json:"id,omitempty"`
	} `json:"src_registry,omitempty"`
	DestRegistry struct {
		ID int `json:"id,omitempty"`
	} `json:"dest_registry,omitempty"`
	DestNamespace        string `json:"dest_namespace,omitempty"`
	DestNamespaceReplace int    `json:"dest_namespace_replace_count"`
	Trigger              struct {
		Type            string `json:"type,omitempty"`
		TriggerSettings struct {
			Cron string `json:"cron,omitempty"`
		} `json:"trigger_settings,omitempty"`
	} `json:"trigger,omitempty"`
	Enabled     bool                 `json:"enabled"`
	Deletion    bool                 `json:"deletion,omitempty"`
	Override    bool                 `json:"override,omitempty"`
	CopyByChunk bool                 `json:"copy_by_chunk,omitempty"`
	Filters     []ReplicationFilters `json:"filters,omitempty"`
	Speed       int                  `json:"speed,omitempty"`
}

type ReplicationFilters struct {
	Type       string      `json:"type,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Decoration string      `json:"decoration,omitempty"`
}

type ExecutionBody struct {
	PolicyID int `json:"policy_id"`
}
