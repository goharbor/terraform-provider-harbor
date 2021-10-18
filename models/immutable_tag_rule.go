package models

var PathImmutableTagRules = "/immutabletagrules"

type ImmutableTagRule struct {
	Id                           int                            `json:"id,omitempty"`
	Disabled                     bool                           `json:"disabled"`
	ScopeSelectors               ScopeSelectors                 `json:"scope_selectors"`
	ImmutableTagRuleTagSelectors []ImmutableTagRuleTagSelectors `json:"tag_selectors"`
	Action                       string                         `json:"action"`
	Template                     string                         `json:"template"`
}

type ImmutableTagRuleTagSelectors struct {
	Kind       string `json:"kind"`
	Decoration string `json:"decoration"`
	Pattern    string `json:"pattern"`
}
