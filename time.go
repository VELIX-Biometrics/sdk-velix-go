package velix

import (
	"context"
	"errors"
)

// ErrTimeNotImplemented é retornado por todos os métodos de TimeModule.
//
// Velix Time NÃO possui nenhum endpoint na superfície pública /v1/api/*
// (ver nota "COBERTURA PARCIAL" em lib-velix-contracts/openapi/public-api.yaml,
// task #593/#616). O GatewayModule de api-velix-identity-core só faz proxy BFF
// para edge, intelligence, copilot e marketplace — não existe proxy para
// api-velix-time hoje. Os escopos time:read/time:write estão reservados no
// ApiScope para quando essa lacuna for endereçada.
var ErrTimeNotImplemented = errors.New("velix: Time (Velix Time) não está implementado — nenhum endpoint /v1/api/time/* existe hoje na API pública (ver task #593/#616)")

// TimeModule placeholder para Velix Time. TODO(#616): implementar quando a
// API pública expuser endpoints de tempo/ponto — hoje não existem.
type TimeModule struct{ c *VelixClient }

// PunchRecords é um stub — sempre retorna ErrTimeNotImplemented.
// TODO(#616): substituir por chamada real quando o endpoint existir.
func (m *TimeModule) PunchRecords(ctx context.Context) error {
	return ErrTimeNotImplemented
}
