package provider

import "testing"

func TestLdapGroupNameEquivalent(t *testing.T) {
	const dn = "cn=harbor_users,cn=groups,dc=example,dc=com"
	const short = "harbor_users"

	cases := []struct {
		name      string
		groupType string
		ldapDN    string
		old, new  string
		want      bool
	}{
		// The main back-compat case: state holds short name (written by the
		// new Read), config still has the DN from before the fix.
		{"ldap DN <-> short both directions (old=short)", "ldap", dn, short, dn, true},
		{"ldap DN <-> short both directions (old=DN)", "ldap", dn, dn, short, true},
		{"ldap identical DN on both sides", "ldap", dn, dn, dn, true},
		{"ldap identical short on both sides", "ldap", dn, short, short, true},

		// Real diffs must not be suppressed.
		{"ldap unrelated name", "ldap", dn, short, "someone_else", false},
		{"ldap DN vs unrelated", "ldap", dn, dn, "cn=other,dc=example,dc=com", false},

		// Non-ldap types must never be suppressed, even with matching values.
		{"internal never suppressed", "internal", "", short, short, false},
		{"oidc never suppressed even with DN", "oidc", dn, dn, short, false},

		// ldap without a DN in state: nothing to compare against.
		{"ldap without DN not suppressed", "ldap", "", short, short, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ldapGroupNameEquivalent(tc.groupType, tc.ldapDN, tc.old, tc.new); got != tc.want {
				t.Errorf("ldapGroupNameEquivalent(%q, %q, %q, %q) = %v, want %v",
					tc.groupType, tc.ldapDN, tc.old, tc.new, got, tc.want)
			}
		})
	}
}
