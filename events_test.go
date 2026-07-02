package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestEventsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "vx_test_key" {
			t.Error("missing api key header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"items": []map[string]any{
					{"id": "evt-1", "name": "Tech Summit", "status": "active"},
				},
				"total": 1, "page": 1, "limit": 20,
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	result, err := client.Events.List(context.Background(), velix.ListOptions{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total=1, got %d", result.Total)
	}
	if len(result.Items) != 1 || result.Items[0].ID != "evt-1" {
		t.Error("unexpected items")
	}
}

func TestEventsGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events/evt-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "evt-1", "name": "Tech Summit", "status": "active"},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	event, err := client.Events.Get(context.Background(), "evt-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.ID != "evt-1" {
		t.Errorf("expected id=evt-1, got %s", event.ID)
	}
}

func TestEventsGet_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "not found", "code": "NOT_FOUND"})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	_, err := client.Events.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*velix.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestEventsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "evt-new", "name": "New Event", "status": "draft"},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	event, err := client.Events.Create(context.Background(), velix.CreateEventInput{
		Name: "New Event",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.ID != "evt-new" {
		t.Errorf("expected id=evt-new, got %s", event.ID)
	}
}
