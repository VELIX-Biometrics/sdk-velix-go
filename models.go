package velix

// CheckinResult é o resultado de uma identificação biométrica.
type CheckinResult struct {
	Passed   bool   `json:"passed"`
	PersonID string `json:"personId,omitempty"`
	Message  string `json:"message,omitempty"`
}

// Person representa um colaborador cadastrado no tenant.
type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	ExternalID  string `json:"externalId,omitempty"`
	EnrolledAt  string `json:"enrolledAt,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

// CreatePersonInput parâmetros para criação de pessoa.
type CreatePersonInput struct {
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}

// UpdatePersonInput parâmetros para atualização de pessoa.
type UpdatePersonInput struct {
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}

// EnrollInput parâmetros para enroll biométrico.
type EnrollInput struct {
	Frames []string `json:"frames"` // base64 JPEG
}

// Event representa um evento gerenciado pelo tenant.
type Event struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TenantID  string `json:"tenantId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// CreateEventInput parâmetros para criação de evento.
type CreateEventInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	StartAt     string `json:"startAt,omitempty"`
	EndAt       string `json:"endAt,omitempty"`
}

// EventConfigInput parâmetros para configuração de evento.
type EventConfigInput struct {
	CheckinMode    string `json:"checkinMode,omitempty"`
	AllowWalkIn    bool   `json:"allowWalkIn,omitempty"`
	BadgeTemplate  string `json:"badgeTemplate,omitempty"`
}

// TenantSettings configurações do tenant.
type TenantSettings struct {
	RequireLiveness       bool    `json:"requireLiveness"`
	BiometricQualityLevel string  `json:"biometricQualityLevel"`
	GeofenceRadiusMetros  float64 `json:"geofenceRadiusMetros"`
	AllowOfflinePunch     bool    `json:"allowOfflinePunch"`
	Timezone              string  `json:"timezone"`
	WebhookURL            string  `json:"webhookUrl,omitempty"`
}

// ListResponse envelope paginado genérico.
type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// ListOptions parâmetros de paginação/filtro.
type ListOptions struct {
	Page   int
	Limit  int
	Search string
}
