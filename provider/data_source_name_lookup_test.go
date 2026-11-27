package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/goharbor/terraform-provider-harbor/client"
)

// newHarborListServer serves the list endpoints the name lookups use, the way
// Harbor does: ?name= is a substring match and a page holds 10 results unless
// page_size is given. Results come back in the order of names, and
// lookupTestNames puts the exact name last, as a real Harbor 2.15 listed
// "billing" after "a01-billing" ... "a12-billing".
func newHarborListServer(t *testing.T, names []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/scanner") {
			http.NotFound(w, r)
			return
		}
		query := r.URL.Query()
		page, pageSize := 1, 10
		if v, err := strconv.Atoi(query.Get("page")); err == nil {
			page = v
		}
		if v, err := strconv.Atoi(query.Get("page_size")); err == nil {
			pageSize = v
		}

		var matched []string
		for _, n := range names {
			if strings.Contains(n, query.Get("name")) {
				matched = append(matched, n)
			}
		}

		items := []map[string]interface{}{}
		for i := (page - 1) * pageSize; i < len(matched) && i < page*pageSize; i++ {
			id := i + 1
			switch r.URL.Path {
			case "/projects":
				items = append(items, map[string]interface{}{"project_id": id, "name": matched[i], "metadata": map[string]string{"public": "false"}})
			case "/registries":
				items = append(items, map[string]interface{}{"id": id, "name": matched[i], "type": "docker-hub", "url": "https://hub.docker.com"})
			case "/scanners":
				items = append(items, map[string]interface{}{"uuid": fmt.Sprintf("uuid-%d", id), "name": matched[i]})
			default:
				http.NotFound(w, r)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	}))
}

func lookupTestNames(target string) []string {
	var names []string
	for i := 1; i <= 12; i++ {
		names = append(names, fmt.Sprintf("a%02d-%s", i, target))
	}
	return append(names, target)
}

func TestDataProjectReadFindsNameOnLaterPage(t *testing.T) {
	server := newHarborListServer(t, lookupTestNames("billing"))
	defer server.Close()
	apiClient := client.NewClient(server.URL, "user", "pass", "", "", false, "", nil)

	d := dataProject().TestResourceData()
	if err := d.Set("name", "billing"); err != nil {
		t.Fatal(err)
	}
	if err := dataProjectRead(d, apiClient); err != nil {
		t.Fatalf("dataProjectRead: %v", err)
	}
	if got, want := d.Id(), "/projects/13"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got := d.Get("project_id").(int); got != 13 {
		t.Errorf("project_id = %d, want 13", got)
	}
}

func TestDataRegistryReadFindsNameOnLaterPage(t *testing.T) {
	server := newHarborListServer(t, lookupTestNames("dockerhub"))
	defer server.Close()
	apiClient := client.NewClient(server.URL, "user", "pass", "", "", false, "", nil)

	d := dataRegistry().TestResourceData()
	if err := d.Set("name", "dockerhub"); err != nil {
		t.Fatal(err)
	}
	if err := dataRegistryRead(d, apiClient); err != nil {
		t.Fatalf("dataRegistryRead: %v", err)
	}
	if got := d.Get("registry_id").(int); got != 13 {
		t.Errorf("registry_id = %d, want 13", got)
	}
	if got := d.Get("type").(string); got != "docker-hub" {
		t.Errorf("type = %q, want docker-hub", got)
	}
}

func TestGetScannerByNameFindsScannerPastFirstPage(t *testing.T) {
	server := newHarborListServer(t, lookupTestNames("trivy"))
	defer server.Close()
	apiClient := client.NewClient(server.URL, "user", "pass", "", "", false, "", nil)

	scanner, err := apiClient.GetScannerByName("trivy")
	if err != nil {
		t.Fatalf("GetScannerByName: %v", err)
	}
	if scanner.UUID != "uuid-13" {
		t.Errorf("uuid = %q, want uuid-13", scanner.UUID)
	}

	if _, err := apiClient.GetScannerByName("clair"); err == nil {
		t.Error("GetScannerByName(clair) succeeded, want not found")
	}
}
