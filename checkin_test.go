package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/velix-io/velix-sdk-go"
)

func TestCheckinIdentify_Matched(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/api/checkin/identify" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "vx_test_key" {
			t.Errorf("missing api key header")
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["imageBase64"] != "base64frame" {
			t.Errorf("expected imageBase64=base64frame, got %v", body["imageBase64"])
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"matched": true, "person_id": "person-uuid", "quality_score": 0.92},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	result, err := client.Checkin.Identify(context.Background(), velix.CheckinIdentifyRequest{ImageBase64: "base64frame"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matched {
		t.Error("expected matched=true")
	}
	if result.PersonID == nil || *result.PersonID != "person-uuid" {
		t.Errorf("expected person_id=person-uuid, got %v", result.PersonID)
	}
}

func TestCheckinIdentify_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid api key", "code": "UNAUTHORIZED"})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "bad"})
	_, err := client.Checkin.Identify(context.Background(), velix.CheckinIdentifyRequest{ImageBase64: "frame"})
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*velix.AuthError); !ok {
		t.Errorf("expected AuthError, got %T", err)
	}
}
