package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ContextModule cobre /v1/contexts/* — Identity Context (Velix.ID). BearerAuth (JWT de sessão).
type ContextModule struct{ c *VelixClient }

// Create cria um novo contexto. POST /v1/contexts.
func (m *ContextModule) Create(ctx context.Context, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, "/v1/contexts", payload)
}

// Get busca um contexto por id. GET /v1/contexts/{id}.
func (m *ContextModule) Get(ctx context.Context, id string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodGet, "/v1/contexts/"+url.PathEscape(id), nil)
}

// List lista contextos do tenant. GET /v1/contexts.
func (m *ContextModule) List(ctx context.Context) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodGet, "/v1/contexts", nil)
}

// Update atualiza um contexto. PATCH /v1/contexts/{id}.
func (m *ContextModule) Update(ctx context.Context, id string, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPatch, "/v1/contexts/"+url.PathEscape(id), payload)
}

// Remove remove (soft delete) um contexto. DELETE /v1/contexts/{id}.
func (m *ContextModule) Remove(ctx context.Context, id string) error {
	_, err := m.c.do(ctx, http.MethodDelete, "/v1/contexts/"+url.PathEscape(id), nil)
	return err
}

// Authorize avalia autorização de uma identidade em um contexto (Authorization Engine).
// POST /v1/contexts/{contextId}/authorize.
func (m *ContextModule) Authorize(ctx context.Context, contextID string, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/contexts/%s/authorize", url.PathEscape(contextID)), payload)
}

// ListAuthorizationDecisions busca/audita decisões de autorização de um contexto.
// GET /v1/contexts/{contextId}/authorization-decisions.
func (m *ContextModule) ListAuthorizationDecisions(ctx context.Context, contextID string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/contexts/%s/authorization-decisions", url.PathEscape(contextID)), nil)
}

// CreateLinkRequest solicita vínculo cross-tenant (consentimento pendente da pessoa).
// POST /v1/contexts/{contextId}/link-requests. Nunca cria membership diretamente:
// retorna 202 (PENDING) aguardando consentimento via magic link/notificação. A API
// pública não expõe approve/reject — isso acontece fora do SDK.
func (m *ContextModule) CreateLinkRequest(ctx context.Context, contextID string, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/contexts/%s/link-requests", url.PathEscape(contextID)), payload)
}

// ContextMembershipModule cobre memberships de contexto e identidade.
type ContextMembershipModule struct{ c *VelixClient }

// Create vincula uma identidade a um contexto. POST /v1/contexts/{contextId}/memberships.
func (m *ContextMembershipModule) Create(ctx context.Context, contextID string, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/contexts/%s/memberships", url.PathEscape(contextID)), payload)
}

// ListByContext lista memberships de um contexto. GET /v1/contexts/{contextId}/memberships.
func (m *ContextMembershipModule) ListByContext(ctx context.Context, contextID string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/contexts/%s/memberships", url.PathEscape(contextID)), nil)
}

// ListByIdentity lista memberships de uma identidade. GET /v1/identities/{identityId}/memberships.
func (m *ContextMembershipModule) ListByIdentity(ctx context.Context, identityID string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/identities/%s/memberships", url.PathEscape(identityID)), nil)
}

// UpdateStatus atualiza o status de um membership. PATCH /v1/memberships/{membershipId}/status.
// status="revoked" é a saída de contexto (definitiva, sem carência, task #834).
func (m *ContextMembershipModule) UpdateStatus(ctx context.Context, membershipID, status string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/memberships/%s/status", url.PathEscape(membershipID)), map[string]any{"status": status})
}

// AddRoles adiciona roles a um membership. POST /v1/memberships/{membershipId}/roles.
func (m *ContextMembershipModule) AddRoles(ctx context.Context, membershipID string, roleIDs []string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/memberships/%s/roles", url.PathEscape(membershipID)), map[string]any{"roleIds": roleIDs})
}

// RemoveRoles remove roles de um membership. POST /v1/memberships/{membershipId}/roles/remove.
func (m *ContextMembershipModule) RemoveRoles(ctx context.Context, membershipID string, roleIDs []string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/memberships/%s/roles/remove", url.PathEscape(membershipID)), map[string]any{"roleIds": roleIDs})
}

// ContextRoleModule cobre /v1/context-roles*.
type ContextRoleModule struct{ c *VelixClient }

// Create cria uma role de contexto. POST /v1/context-roles.
func (m *ContextRoleModule) Create(ctx context.Context, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, "/v1/context-roles", payload)
}

// List lista roles de contexto por context_type. GET /v1/context-roles?contextType=....
func (m *ContextRoleModule) List(ctx context.Context, contextType string) (json.RawMessage, error) {
	q := url.Values{"contextType": {contextType}}
	return m.c.do(ctx, http.MethodGet, "/v1/context-roles?"+q.Encode(), nil)
}

// LinkPermissions vincula permissions a uma role. POST /v1/context-roles/{roleId}/permissions.
func (m *ContextRoleModule) LinkPermissions(ctx context.Context, roleID string, permissionIDs []string) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/context-roles/%s/permissions", url.PathEscape(roleID)), map[string]any{"permissionIds": permissionIDs})
}

// ContextPermissionModule cobre /v1/context-permissions.
type ContextPermissionModule struct{ c *VelixClient }

// Create cria uma permission de contexto. POST /v1/context-permissions.
func (m *ContextPermissionModule) Create(ctx context.Context, payload map[string]any) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, "/v1/context-permissions", payload)
}

// List lista permissions de contexto, opcionalmente filtrando por categoria.
// GET /v1/context-permissions?category=....
func (m *ContextPermissionModule) List(ctx context.Context, category string) (json.RawMessage, error) {
	path := "/v1/context-permissions"
	if category != "" {
		path += "?" + (url.Values{"category": {category}}).Encode()
	}
	return m.c.do(ctx, http.MethodGet, path, nil)
}

// AuthorizationTokenModule cobre /v1/authorization-tokens/validate.
type AuthorizationTokenModule struct{ c *VelixClient }

// Validate valida (e opcionalmente consome) um token de autorização vat_*.
// POST /v1/authorization-tokens/validate.
func (m *AuthorizationTokenModule) Validate(ctx context.Context, token string, consume bool) (json.RawMessage, error) {
	return m.c.do(ctx, http.MethodPost, "/v1/authorization-tokens/validate", map[string]any{"token": token, "consume": consume})
}
