package models

var PathRetentions = "/retentions"

type Retention struct {
	Id        int     `json:"id,omitempty"`
	Algorithm string  `json:"algorithm"`
	Rules     []Rules `json:"rules"`
	Trigger   Trigger `json:"trigger"`
	Scope     Scope   `json:"scope"`
}
type Repository struct {
	Kind       string `json:"kind"`
	Decoration string `json:"decoration"`
	Pattern    string `json:"pattern"`
}
type ScopeSelectors struct {
	Repository []Repository `json:"repository"`
}
type TagSelectors struct {
	Kind       string `json:"kind"`
	Decoration string `json:"decoration"`
	Pattern    string `json:"pattern"`
	Extras     string `json:"extras"`
}
type Params struct {
	LatestPushedK      int `json:"latestPushedK,omitempty"`
	LatestPulledN      int `json:"latestPulledN,omitempty"`
	NDaysSinceLastPush int `json:"nDaysSinceLastPush,omitempty"`
	NDaysSinceLastPull int `json:"nDaysSinceLastPull,omitempty"`
}
type Rules struct {
	Disabled       bool           `json:"disabled"`
	Action         string         `json:"action"`
	ScopeSelectors ScopeSelectors `json:"scope_selectors"`
	TagSelectors   []TagSelectors `json:"tag_selectors"`
	Params         Params         `json:"params"`
	Template       string         `json:"template"`
}
type References struct {
}
type Settings struct {
	Cron string `json:"cron"`
}
type Trigger struct {
	Kind       string     `json:"kind"`
	References References `json:"references"`
	Settings   Settings   `json:"settings"`
}
type Scope struct {
	Level string `json:"level"`
	Ref   int    `json:"ref"`
}
