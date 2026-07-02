package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestEventsCreateGuest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/api/events/evt-1/guests" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") != "vx_test_key" {
			t.Error("missing api key header")
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["birthDate"] != "1990-01-01" {
			t.Errorf("expected birthDate=1990-01-01, got %v", body["birthDate"])
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "guest-1", "eventId": "evt-1", "name": "Jane Doe",
				"email": "jane@example.com", "status": "invited",
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	guest, err := client.Events.CreateGuest(context.Background(), "evt-1", velix.CreateGuestRequest{
		Name: "Jane Doe", Email: "jane@example.com", BirthDate: "1990-01-01",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if guest.ID != "guest-1" || guest.EventID != "evt-1" {
		t.Errorf("unexpected guest: %+v", guest)
	}
}

func TestEventsGetGuest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/api/events/evt-1/guests/guest-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "guest-1", "eventId": "evt-1", "name": "Jane Doe",
				"email": "jane@example.com", "status": "checked_in",
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	guest, err := client.Events.GetGuest(context.Background(), "evt-1", "guest-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if guest.Status != "checked_in" {
		t.Errorf("expected status=checked_in, got %s", guest.Status)
	}
}

func TestEventsGetGuest_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "not found", "code": "NOT_FOUND"})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	_, err := client.Events.GetGuest(context.Background(), "evt-1", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*velix.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}
