package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TenantsModule gerenciamento de configurações do tenant.
type TenantsModule struct{ c *VelixClient }

// Tenant dados do tenant autenticado.
type Tenant struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Slug     string         `json:"slug"`
	Settings TenantSettings `json:"settings"`
}

// Me retorna os dados do tenant da API key autenticada.
func (m *TenantsModule) Me(ctx context.Context) (*Tenant, error) {
	raw, err := m.c.do(ctx, http.MethodGet, "/v1/tenants/me", nil)
	if err != nil {
		return nil, err
	}
	var t Tenant
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, fmt.Errorf("velix: decode tenant: %w", err)
	}
	return &t, nil
}

// UpdateSettings atualiza configurações do tenant.
func (m *TenantsModule) UpdateSettings(ctx context.Context, settings TenantSettings) (*Tenant, error) {
	raw, err := m.c.do(ctx, http.MethodPut, "/v1/tenants/me/settings", settings)
	if err != nil {
		return nil, err
	}
	var t Tenant
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, fmt.Errorf("velix: decode tenant: %w", err)
	}
	return &t, nil
}
