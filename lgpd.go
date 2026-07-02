package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LGPDModule solicitações de exclusão de dados via API key (Velix.ID).
type LGPDModule struct{ c *VelixClient }

// DeletionRequest solicita a exclusão dos dados de uma pessoa. A pessoa deve
// possuir uma Identity ativa vinculada ao tenant dono da API key, caso
// contrário a API retorna 403. Escopo exigido: lgpd:write.
// POST /v1/api/deletion-request.
func (m *LGPDModule) DeletionRequest(ctx context.Context, personID string) (*DeletionRequestResponse, error) {
	raw, err := m.c.do(ctx, http.MethodPost, "/v1/api/deletion-request", DeletionRequestBody{PersonID: personID})
	if err != nil {
		return nil, err
	}
	var result DeletionRequestResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode deletion request result: %w", err)
	}
	return &result, nil
}
