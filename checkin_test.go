package velix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	velix "github.com/VELIX-Biometrics/sdk-velix-go"
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
			"data": map[string]any{
				"match":       true,
				"subjectId":   "person-uuid",
				"subjectName": "Ana Silva",
				"liveness":    map[string]any{"ok": true},
				"model":       "adaface",
			},
		})
	}))
	defer srv.Close()

	client := velix.NewClient(velix.Config{APIURL: srv.URL, APIKey: "vx_test_key"})
	result, err := client.Checkin.Identify(context.Background(), velix.CheckinIdentifyRequest{ImageBase64: "base64frame"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Match {
		t.Error("expected match=true")
	}
	if result.SubjectID == nil || *result.SubjectID != "person-uuid" {
		t.Errorf("expected subjectId=person-uuid, got %v", result.SubjectID)
	}
	if !result.Liveness.OK {
		t.Error("expected liveness.ok=true")
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
