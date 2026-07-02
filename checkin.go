package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CheckinModule métodos de identificação biométrica.
type CheckinModule struct{ c *VelixClient }

type facialPayload struct {
	Frame          string   `json:"frame"`
	LivenessSamples []string `json:"livenessSamples,omitempty"`
}

// Facial identifica uma pessoa por frame facial (base64 JPEG).
func (m *CheckinModule) Facial(ctx context.Context, tenantSlug, frameBase64 string, livenessSamples ...string) (*CheckinResult, error) {
	payload := facialPayload{Frame: frameBase64, LivenessSamples: livenessSamples}
	raw, err := m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/checkin/%s/identify", tenantSlug), payload)
	if err != nil {
		return nil, err
	}
	var result CheckinResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode checkin result: %w", err)
	}
	return &result, nil
}
