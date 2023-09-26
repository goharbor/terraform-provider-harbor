package models

var PathConfig = "/configurations"

type ConfigBodyAuthPost struct {
	OidcVerifyCert        bool   `json:"oidc_verify_cert"`
	OidcAutoOnboard       bool   `json:"oidc_auto_onboard"`
	OidcUserClaim         string `json:"oidc_user_claim,omitempty"`
	LdapGroupSearchFilter string `json:"ldap_group_search_filter,omitempty"`
	AuthMode              string `json:"auth_mode,omitempty"`
	PrimaryAuthMode       bool   `json:"primary_auth_mode,omitempty"`
	SelfRegistration      bool   `json:"self_registration"`
	OidcScope             string `json:"oidc_scope,omitempty"`
	LdapSearchDn          string `json:"ldap_search_dn,omitempty"`
	ScanAllPolicy         struct {
		Type      string `json:"type,omitempty"`
		Parameter struct {
			DailyTime int `json:"daily_time,omitempty"`
		} `json:"parameter,omitempty"`
	} `json:"scan_all_policy,omitempty"`
	LdapTimeout                  int    `json:"ldap_timeout,omitempty"`
	LdapBaseDn                   string `json:"ldap_base_dn,omitempty"`
	LdapFilter                   string `json:"ldap_filter,omitempty"`
	LdapURL                      string `json:"ldap_url,omitempty"`
	OidcName                     string `json:"oidc_name,omitempty"`
	LdapUID                      string `json:"ldap_uid,omitempty"`
	OidcClientID                 string `json:"oidc_client_id,omitempty"`
	LdapGroupBaseDn              string `json:"ldap_group_base_dn,omitempty"`
	LdapGroupAttributeName       string `json:"ldap_group_attribute_name,omitempty"`
	LdapGroupMembershipAttribute string `json:"ldap_group_membership_attribute,omitempty"`
	LdapSearchPassword           string `json:"ldap_search_password,omitempty"`
	LdapGroupAdminDn             string `json:"ldap_group_admin_dn,omitempty"`
	OidcEndpoint                 string `json:"oidc_endpoint,omitempty"`
	OidcClientSecret             string `json:"oidc_client_secret,omitempty"`
	OidcGroupFilter              string `json:"oidc_group_filter,omitempty"`
	OidcGroupsClaim              string `json:"oidc_groups_claim,omitempty"`
	LdapScope                    int    `json:"ldap_scope,omitempty"`
	TokenExpiration              int    `json:"token_expiration,omitempty"`
	LdapGroupSearchScope         int    `json:"ldap_group_search_scope,omitempty"`
	LdapVerifyCert               bool   `json:"ldap_verify_cert"`
	LdapGroupGID                 string `json:"ldap_group_gid,omitempty"`
	OidcAdminGroup               string `json:"oidc_admin_group,omitempty"`
}

type ConfigBodySystemPost struct {
	ProjectCreationRestriction string `json:"project_creation_restriction,omitempty"`
	ReadOnly                   bool   `json:"read_only"`
	RobotTokenDuration         int    `json:"robot_token_duration,omitempty"`
	QuotaPerProjectEnable      bool   `json:"quota_per_project_enable"`
	RobotNamePrefix            string `json:"robot_name_prefix,omitempty"`
	StoragePerProject          int    `json:"storage_per_project,omitempty"`
	ScannerSkipUpdatePulltime  bool   `json:"scanner_skip_update_pulltime"`
}

type ConfigBodyEmailPost struct {
	EmailHost     string `json:"email_host,omitempty"`
	EmailPort     int    `json:"email_port,omitempty"`
	EmailUsername string `json:"email_username,omitempty"`
	EmailPassword string `json:"email_password,omitempty"`
	EmailFrom     string `json:"email_from,omitempty"`
	EmailSsl      bool   `json:"email_ssl"`
	EmailInsecure bool   `json:"email_insecure"`
	EmailIdentity string `json:"email_identity,omitempty"`
}

type ConfigBodyResponse struct {
	OidcVerifyCert struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"oidc_verify_cert,omitempty"`
	OidcAutoOnboard struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"oidc_auto_onboard,omitempty"`
	OidcUserClaim struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_user_claim,omitempty"`
	OidcGroupFilter struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_group_filter,omitempty"`
	OidcGroupsClaim struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_groups_claim,omitempty"`
	EmailIdentity struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"email_identity,omitempty"`
	LdapGroupSearchFilter struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_search_filter,omitempty"`
	AuthMode struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"auth_mode,omitempty"`
	PrimaryAuthMode struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"primary_auth_mode,omitempty"`
	SelfRegistration struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"self_registration,omitempty"`
	OidcScope struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_scope,omitempty"`
	LdapSearchDn struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_search_dn,omitempty"`
	StoragePerProject struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"storage_per_project,omitempty"`
	ScanAllPolicy struct {
		Type      string `json:"type,omitempty"`
		Parameter struct {
			DailyTime int `json:"daily_time,omitempty"`
		} `json:"parameter,omitempty"`
	} `json:"scan_all_policy",omitempty`
	VerifyRemoteCert struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"verify_remote_cert,omitempty"`
	LdapTimeout struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"ldap_timeout,omitempty"`
	LdapBaseDn struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_base_dn,omitempty"`
	LdapFilter struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_filter,omitempty"`
	ReadOnly struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"read_only,omitempty"`
	QuotaPerProjectEnable struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"quota_per_project_enable,omitempty"`
	LdapURL struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_url,omitempty"`
	OidcName struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_name,omitempty"`
	ProjectCreationRestriction struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"project_creation_restriction,omitempty"`
	LdapUID struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_uid,omitempty"`
	OidcClientID struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_client_id,omitempty"`
	LdapGroupBaseDn struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_base_dn,omitempty"`
	LdapGroupAttributeName struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_attribute_name,omitempty"`
	EmailInsecure struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"email_insecure,omitempty"`
	LdapGroupAdminDn struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_admin_dn,omitempty"`
	LdapGroupMembershipAttribute struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_membership_attribute,omitempty"`
	EmailUsername struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"email_username,omitempty"`
	OidcEndpoint struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_endpoint,omitempty"`
	LdapScope struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"ldap_scope,omitempty"`
	TokenExpiration struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"token_expiration,omitempty"`
	LdapGroupSearchScope struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"ldap_group_search_scope,omitempty"`
	EmailSsl struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"email_ssl,omitempty"`
	EmailPort struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"email_port,omitempty"`
	EmailHost struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"email_host,omitempty"`
	EmailFrom struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"email_from,omitempty"`
	LdapGroupGID struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"ldap_group_gid,omitempty"`
	OidcAdminGroup struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"oidc_admin_group,omitempty"`
	RobotNamePrefix struct {
		Editable bool   `json:"editable,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"robot_name_prefix,omitempty"`
	RobotTokenDuration struct {
		Editable bool `json:"editable,omitempty"`
		Value    int  `json:"value,omitempty"`
	} `json:"robot_token_duration,omitempty"`
	LdapVerifyCert struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"ldap_verify_cert,omitempty"`
	ScannerSkipUpdatePulltime struct {
		Editable bool `json:"editable,omitempty"`
		Value    bool `json:"value,omitempty"`
	} `json:"scanner_skip_update_pulltime,omitempty"`
}
