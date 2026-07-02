package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CheckinModule identificação biométrica via API key (Velix.ID).
type CheckinModule struct{ c *VelixClient }

// Identify identifica uma pessoa por frame facial. Escopo exigido na API key:
// checkin:write. POST /v1/api/checkin/identify.
func (m *CheckinModule) Identify(ctx context.Context, req CheckinIdentifyRequest) (*CheckinIdentifyResponse, error) {
	raw, err := m.c.do(ctx, http.MethodPost, "/v1/api/checkin/identify", req)
	if err != nil {
		return nil, err
	}
	var result CheckinIdentifyResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode checkin result: %w", err)
	}
	return &result, nil
}
