package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestTenantsMe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tenants/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "tenant-uuid", "name": "Acme Corp", "slug": "acme",
				"plan": "enterprise", "maxPersons": 1000,
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	tenant, err := client.Tenants.Me(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant.ID != "tenant-uuid" {
		t.Errorf("expected id=tenant-uuid, got %s", tenant.ID)
	}
	if tenant.Slug != "acme" {
		t.Errorf("expected slug=acme, got %s", tenant.Slug)
	}
}

func TestTenantsUpdateSettings(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tenants/me/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "tenant-uuid", "requireLiveness": true, "timezone": "America/Sao_Paulo",
			},
		})
	}))
	defer srv.Close()

	requireLiveness := true
	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	tenant, err := client.Tenants.UpdateSettings(context.Background(), velix.UpdateTenantSettingsInput{
		RequireLiveness: &requireLiveness,
		Timezone:        "America/Sao_Paulo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tenant.RequireLiveness {
		t.Error("expected requireLiveness=true")
	}
}
