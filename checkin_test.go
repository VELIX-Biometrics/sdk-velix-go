package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestCheckinFacial_Passed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/checkin/tenant-slug/identify" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "vx_test_key" {
			t.Errorf("missing api key header")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"passed": true, "personId": "person-uuid"},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	result, err := client.Checkin.Facial(context.Background(), "tenant-slug", "base64frame")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Error("expected passed=true")
	}
	if result.PersonID != "person-uuid" {
		t.Errorf("expected personId=person-uuid, got %s", result.PersonID)
	}
}

func TestCheckinFacial_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid api key", "code": "UNAUTHORIZED"})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "bad"})
	_, err := client.Checkin.Facial(context.Background(), "slug", "frame")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*velix.AuthError); !ok {
		t.Errorf("expected AuthError, got %T", err)
	}
}
