package models

var PathConfig = "/configurations"

type ConfigBodyPost struct {
	OidcVerifyCert        bool   `json:"oidc_verify_cert,omitempty"`
	EmailIdentity         string `json:"email_identity,omitempty"`
	LdapGroupSearchFilter string `json:"ldap_group_search_filter,omitempty"`
	AuthMode              string `json:"auth_mode,omitempty"`
	SelfRegistration      bool   `json:"self_registration,omitempty"`
	OidcScope             string `json:"oidc_scope,omitempty"`
	LdapSearchDn          string `json:"ldap_search_dn,omitempty"`
	StoragePerProject     string `json:"storage_per_project,omitempty"`
	ScanAllPolicy         struct {
		Type      string `json:"type,omitempty"`
		Parameter struct {
			DailyTime int `json:"daily_time,omitempty"`
		} `json:"parameter,omitempty"`
	} `json:"scan_all_policy,omitempty"`
	VerifyRemoteCert           bool   `json:"verify_remote_cert,omitempty"`
	LdapTimeout                int    `json:"ldap_timeout,omitempty"`
	LdapBaseDn                 string `json:"ldap_base_dn,omitempty"`
	LdapFilter                 string `json:"ldap_filter,omitempty"`
	ReadOnly                   bool   `json:"read_only,omitempty"`
	QuotaPerProjectEnable      bool   `json:"quota_per_project_enable,omitempty"`
	LdapURL                    string `json:"ldap_url,omitempty"`
	OidcName                   string `json:"oidc_name,omitempty"`
	ProjectCreationRestriction string `json:"project_creation_restriction,omitempty"`
	LdapUID                    string `json:"ldap_uid,omitempty"`
	OidcClientID               string `json:"oidc_client_id,omitempty"`
	LdapGroupBaseDn            string `json:"ldap_group_base_dn,omitempty"`
	LdapGroupAttributeName     string `json:"ldap_group_attribute_name,omitempty"`
	EmailInsecure              bool   `json:"email_insecure,omitempty"`
	LdapGroupAdminDn           string `json:"ldap_group_admin_dn,omitempty"`
	EmailUsername              string `json:"email_username,omitempty"`
	EmailPassword              string `json:"email_password,omitempty"`
	OidcEndpoint               string `json:"oidc_endpoint,omitempty"`
	OidcClientSecret           string `json:"oidc_client_secret,omitempty"`
	OidcGroupsClaim            string `json:"oidc_groups_claim,omitempty"`
	LdapScope                  int    `json:"ldap_scope,omitempty"`
	CountPerProject            string `json:"count_per_project,omitempty"`
	TokenExpiration            int    `json:"token_expiration,omitempty"`
	LdapGroupSearchScope       int    `json:"ldap_group_search_scope,omitempty"`
	EmailSsl                   bool   `json:"email_ssl,omitempty"`
	EmailPort                  int    `json:"email_port,omitempty"`
	EmailHost                  string `json:"email_host,omitempty"`
	EmailFrom                  string `json:"email_from,omitempty"`
	RobotTokenDuration         int    `json:"robot_token_duration,omitempty"`
}
