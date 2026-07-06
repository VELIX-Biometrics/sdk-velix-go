//go:build ignore

// Smoke test de contrato — rode com `go run smoke_test.go` dentro do módulo.
// Usa o SDK Go de verdade (pacote local, não o módulo publicado).
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	velix "github.com/VELIX-Biometrics/sdk-velix-go"
)

const img = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func result(step string, ok bool, detail string) {
	status := "PASS"
	if !ok {
		status = "FAIL"
	}
	fmt.Printf("RESULT:go:%s:%s:%s\n", step, status, detail)
}

func reachable(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, s := range []string{"route not found", "no route", "401", "403"} {
		if strings.Contains(msg, s) {
			return false
		}
	}
	return true
}

func main() {
	ctx := context.Background()
	client := velix.NewClient(velix.Config{
		APIURL: os.Getenv("API_BASE_URL"),
		APIKey: os.Getenv("VELIX_API_KEY"),
	})

	var personID string
	if r, err := client.Onboarding.Create(ctx, velix.OnboardingRequest{Name: "Smoke Test Go", Frames: []string{img, img, img}}); err != nil {
		result("onboarding", false, err.Error())
	} else {
		personID = r.PersonID
		result("onboarding", personID != "", "person_id="+personID)
	}

	if r, err := client.Checkin.Identify(ctx, velix.CheckinIdentifyRequest{ImageBase64: img}); err != nil {
		result("checkin", false, err.Error())
	} else {
		result("checkin", true, fmt.Sprintf("match=%v", r.Match))
	}

	if personID != "" {
		if _, err := client.LGPD.DeletionRequest(ctx, personID); err != nil {
			result("lgpd", false, err.Error())
		} else {
			result("lgpd", true, "deletion-request ok")
		}

		if _, err := client.Me.Get(ctx, personID); err != nil {
			result("me", false, err.Error())
		} else {
			result("me", true, "got response")
		}
	}

	dummy := "00000000-0000-0000-0000-000000000000"
	if _, err := client.Events.CreateGuest(ctx, dummy, velix.CreateGuestRequest{Name: "Guest Smoke", Email: "guest@smoke.test"}); err != nil {
		result("events_create", reachable(err), err.Error())
	} else {
		result("events_create", true, "endpoint reachable")
	}

	if _, err := client.Events.GetGuest(ctx, dummy, dummy); err != nil {
		result("events_get", reachable(err), err.Error())
	} else {
		result("events_get", true, "endpoint reachable")
	}
}
