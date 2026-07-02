package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EventsModule convidados de evento via API key (Velix Events). Cobertura
// MÍNIMA — apenas criar e consultar convidado, único mapeado na spec pública
// (task #593). Não expõe CRUD de eventos: não existe superfície de API key
// para isso hoje em api-velix-identity-core.
type EventsModule struct{ c *VelixClient }

// CreateGuest cria um convidado de evento. Escopo exigido: events:write.
// POST /v1/api/events/{id}/guests.
func (m *EventsModule) CreateGuest(ctx context.Context, eventID string, req CreateGuestRequest) (*GuestResponse, error) {
	raw, err := m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/api/events/%s/guests", eventID), req)
	if err != nil {
		return nil, err
	}
	var result GuestResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode guest: %w", err)
	}
	return &result, nil
}

// GetGuest consulta um convidado de evento, incluindo status de checkin.
// Escopo exigido: events:read. GET /v1/api/events/{id}/guests/{guestId}.
func (m *EventsModule) GetGuest(ctx context.Context, eventID, guestID string) (*GuestResponse, error) {
	raw, err := m.c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/api/events/%s/guests/%s", eventID, guestID), nil)
	if err != nil {
		return nil, err
	}
	var result GuestResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode guest: %w", err)
	}
	return &result, nil
}
