package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
)

func TestUserReadRemovesTheResourceWhenTheUserIsGone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"code":"NOT_FOUND","message":"user 42 not found"}]}`))
	}))
	t.Cleanup(server.Close)

	apiClient := client.NewClient(server.URL, "admin", "Harbor12345", "", "", true, "", nil)
	d := resourceUser().TestResourceData()
	d.SetId("/users/42")

	if err := resourceUserRead(d, apiClient); err != nil {
		t.Fatalf("expected no error for a user deleted outside Terraform, got %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected the resource to be removed from state, got id %q", d.Id())
	}
}

func TestUserReadReportsOtherErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	apiClient := client.NewClient(server.URL, "admin", "Harbor12345", "", "", true, "", nil)
	d := resourceUser().TestResourceData()
	d.SetId("/users/42")

	if err := resourceUserRead(d, apiClient); err == nil {
		t.Fatal("expected the server error to be reported")
	}
	if d.Id() != "/users/42" {
		t.Fatalf("expected the id to be kept, got %q", d.Id())
	}
}
