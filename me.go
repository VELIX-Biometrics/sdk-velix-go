package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MeModule consulta de dados de pessoa via API key (Velix.ID).
type MeModule struct{ c *VelixClient }

// Get retorna os dados de uma pessoa pelo personId. Valida que a pessoa
// possui Identity vinculada ao tenant dono da API key, caso contrário a API
// retorna 403. Escopo exigido: me:read. GET /v1/api/me/{personId}.
func (m *MeModule) Get(ctx context.Context, personID string) (*MeResponse, error) {
	raw, err := m.c.do(ctx, http.MethodGet, "/v1/api/me/"+personID, nil)
	if err != nil {
		return nil, err
	}
	var result MeResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode me result: %w", err)
	}
	return &result, nil
}
