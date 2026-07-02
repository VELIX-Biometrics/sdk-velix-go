package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EventsModule CRUD de eventos.
type EventsModule struct{ c *VelixClient }

// List retorna todos os eventos do tenant.
func (m *EventsModule) List(ctx context.Context) ([]Event, error) {
	raw, err := m.c.do(ctx, http.MethodGet, "/v1/events", nil)
	if err != nil {
		return nil, err
	}
	var events []Event
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, fmt.Errorf("velix: decode events: %w", err)
	}
	return events, nil
}

// Get retorna um evento pelo ID.
func (m *EventsModule) Get(ctx context.Context, id string) (*Event, error) {
	raw, err := m.c.do(ctx, http.MethodGet, "/v1/events/"+id, nil)
	if err != nil {
		return nil, err
	}
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("velix: decode event: %w", err)
	}
	return &e, nil
}

// Create cria um novo evento.
func (m *EventsModule) Create(ctx context.Context, input CreateEventInput) (*Event, error) {
	raw, err := m.c.do(ctx, http.MethodPost, "/v1/events", input)
	if err != nil {
		return nil, err
	}
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("velix: decode event: %w", err)
	}
	return &e, nil
}

// Configure atualiza configurações de um evento.
func (m *EventsModule) Configure(ctx context.Context, id string, cfg EventConfigInput) error {
	_, err := m.c.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/events/%s/config", id), cfg)
	return err
}
