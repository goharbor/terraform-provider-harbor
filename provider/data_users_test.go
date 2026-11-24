package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
)

// newHarborUsersServer answers GET /users the way Harbor does for
// q=username=...,email=...: q is unescaped once more after the query string
// is decoded (lib/q.Build in Harbor), matches are exact, and the page after
// the last one is empty.
func newHarborUsersServer(t *testing.T, users []map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users" {
			http.NotFound(w, r)
			return
		}
		query := r.URL.Query()
		filters := map[string]string{}
		if q := query.Get("q"); q != "" {
			if unescaped, err := url.QueryUnescape(q); err == nil {
				q = unescaped
			}
			for _, kv := range strings.Split(q, ",") {
				k, v, _ := strings.Cut(kv, "=")
				filters[k] = v
			}
		}

		items := []map[string]interface{}{}
		if query.Get("page") == "1" {
			for _, u := range users {
				if v, ok := filters["username"]; ok && u["username"] != v {
					continue
				}
				if v, ok := filters["email"]; ok && u["email"] != v {
					continue
				}
				items = append(items, u)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	}))
}

func TestDataUsersReadEscapesQueryValues(t *testing.T) {
	server := newHarborUsersServer(t, []map[string]interface{}{
		{"user_id": 3, "username": "ci", "email": "dev+ci@example.com"},
		{"user_id": 4, "username": "dev", "email": "dev@example.com"},
		{"user_id": 5, "username": "ops&audit", "email": "ops@example.com"},
	})
	defer server.Close()
	apiClient := client.NewClient(server.URL, "user", "pass", "", "", false, "", nil)

	tests := []struct {
		attr, value, wantID string
	}{
		{"email", "dev+ci@example.com", "/users/3"},
		{"email", "dev@example.com", "/users/4"},
		{"username", "ops&audit", "/users/5"},
	}
	for _, tt := range tests {
		t.Run(tt.attr+"="+tt.value, func(t *testing.T) {
			d := dataUsers().TestResourceData()
			if err := d.Set(tt.attr, tt.value); err != nil {
				t.Fatal(err)
			}
			if err := dataUsersRead(d, apiClient); err != nil {
				t.Fatalf("dataUsersRead: %v", err)
			}
			users := d.Get("users").([]interface{})
			if len(users) != 1 {
				t.Fatalf("got %d users, want 1", len(users))
			}
			if got := users[0].(map[string]interface{})["id"]; got != tt.wantID {
				t.Errorf("id = %v, want %s", got, tt.wantID)
			}
		})
	}
}
