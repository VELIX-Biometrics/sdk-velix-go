package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OnboardingModule cadastro biométrico via API key (Velix.ID).
type OnboardingModule struct{ c *VelixClient }

// Create realiza o onboarding biométrico de uma pessoa. Escopo exigido na API
// key: onboarding:write. POST /v1/api/onboarding.
func (m *OnboardingModule) Create(ctx context.Context, req OnboardingRequest) (*OnboardingResponse, error) {
	raw, err := m.c.do(ctx, http.MethodPost, "/v1/api/onboarding", req)
	if err != nil {
		return nil, err
	}
	var result OnboardingResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode onboarding result: %w", err)
	}
	return &result, nil
}
