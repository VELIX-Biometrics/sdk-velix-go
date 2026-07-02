package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestPersonsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"data":  []map[string]any{{"id": "p1", "name": "João Silva", "createdAt": "2026-01-01"}},
				"total": 1, "page": 1, "limit": 20,
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "key"})
	list, err := client.Persons.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Total != 1 || list.Data[0].Name != "João Silva" {
		t.Errorf("unexpected result: %+v", list)
	}
}

func TestPersonsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "new-id", "name": "Maria Santos", "createdAt": "2026-01-01"},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "key"})
	p, err := client.Persons.Create(context.Background(), velix.CreatePersonInput{Name: "Maria Santos"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != "new-id" {
		t.Errorf("expected id=new-id, got %s", p.ID)
	}
}
