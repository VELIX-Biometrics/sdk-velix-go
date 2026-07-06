package velix

// ── Onboarding (Velix.ID) — POST /v1/api/onboarding ────────────────────────

// OnboardingRequest contrato real de OnboardingDto (src/modules/onboarding/dto/onboarding.dto.ts).
// Name, Email, Phone, Document, DocumentType e Frames são obrigatórios.
type OnboardingRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Document     string `json:"document"`
	DocumentType string `json:"document_type"` // CPF, CNPJ, RG, PASSPORT, OTHER
	// BirthDate ISO 8601 (ex: "1990-05-20"). Opcional, mas necessário para
	// calcular Age/IsMinor na resposta.
	BirthDate    string         `json:"birth_date,omitempty"`
	ExternalID   string         `json:"external_id,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Frames       []string       `json:"frames"`         // JPEG base64, sem prefixo data URI, mínimo 1
	Role         string         `json:"role,omitempty"` // member, admin, tenant_admin
	AccessGroups []string       `json:"access_groups,omitempty"`
}

// FrameResult resultado de processamento de um frame individual de onboarding.
type FrameResult struct {
	FrameIndex     int     `json:"frame_index"`
	QualityPassed  bool    `json:"quality_passed"`
	QualityScore   float64 `json:"quality_score"`
	LivenessPassed bool    `json:"liveness_passed"`
}

// OnboardingResponse conteúdo de Envelope.data para POST /v1/api/onboarding.
type OnboardingResponse struct {
	PersonID        string        `json:"person_id"`
	IdentityID      string        `json:"identity_id"`
	Enrolled        bool          `json:"enrolled"`
	FramesProcessed int           `json:"frames_processed"`
	FramesResults   []FrameResult `json:"frames_results"`
	EmbeddingID     *string       `json:"embedding_id"`
	Message         string        `json:"message"`
	// Age calculada a partir de BirthDate — nil se BirthDate não foi informado.
	Age     *int  `json:"age"`
	IsMinor *bool `json:"is_minor"`
}

// ── Checkin (Velix.ID) — POST /v1/api/checkin/identify ─────────────────────

// LivenessSample amostra de liveness ativo (contrato mantém camelCase no wire).
type LivenessSample struct {
	Action      string `json:"action"` // center, move_closer, move_away
	ImageBase64 string `json:"imageBase64"`
}

// LivenessBlock bloco opcional de prova de vida ativa.
type LivenessBlock struct {
	Token   string           `json:"token"`
	Samples []LivenessSample `json:"samples"`
}

// CheckinLocation geolocalização opcional do checkin.
type CheckinLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// CheckinIdentifyRequest contrato real de IdentifyFaceDto (src/modules/checkin/dto/identify-face.dto.ts).
type CheckinIdentifyRequest struct {
	ImageBase64 string           `json:"imageBase64"`
	Images      []string         `json:"images,omitempty"`
	TopK        int              `json:"topK,omitempty"`
	Liveness    *LivenessBlock   `json:"liveness,omitempty"`
	Location    *CheckinLocation `json:"location,omitempty"`
}

// LivenessResult bloco de resultado de prova de vida. Score NUNCA é exposto —
// apenas o booleano `ok`.
type LivenessResult struct {
	OK bool `json:"ok"`
}

// CheckinIdentifyResponse resultado da identificação — contrato real de
// CheckinService.identifyFace (checkin.service.ts). Score de similaridade e
// de liveness NUNCA são expostos — apenas os booleanos `match`/`liveness.ok`.
type CheckinIdentifyResponse struct {
	Match       bool           `json:"match"`
	SubjectID   *string        `json:"subjectId"`
	SubjectName *string        `json:"subjectName"`
	Liveness    LivenessResult `json:"liveness"`
	Model       string         `json:"model"`
}

// ── LGPD (Velix.ID) — POST /v1/api/deletion-request ────────────────────────

// DeletionRequestBody corpo de POST /v1/api/deletion-request.
type DeletionRequestBody struct {
	PersonID string `json:"person_id"`
}

// DeletionRequestResponse conteúdo de Envelope.data para POST /v1/api/deletion-request.
type DeletionRequestResponse struct {
	ProtocolNumber string `json:"protocol_number"`
	Message        string `json:"message"`
}

// ── Me (Velix.ID) — GET /v1/api/me/{personId} ───────────────────────────────

// MeResponse conteúdo de Envelope.data para GET /v1/api/me/{personId}.
type MeResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	PhotoURL  *string `json:"photo_url"`
	CreatedAt string  `json:"created_at"`
}

// ── Events (Velix Events) — /v1/api/events/{id}/guests ─────────────────────

// CreateGuestRequest contrato real de create-guest.dto.ts.
// Campos deste schema permanecem em camelCase no wire (birthDate, categoryId,
// companionOf), diferente do restante da superfície /v1/api/* que é snake_case.
type CreateGuestRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CPF         string `json:"cpf,omitempty"`
	Phone       string `json:"phone,omitempty"`
	BirthDate   string `json:"birthDate,omitempty"`
	CategoryID  string `json:"categoryId,omitempty"`
	CompanionOf string `json:"companionOf,omitempty"`
}

// GuestResponse EventGuest — retornado por POST .../guests e GET .../guests/{guestId}.
type GuestResponse struct {
	ID         string  `json:"id"`
	EventID    string  `json:"eventId"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Status     string  `json:"status"`
	CategoryID *string `json:"categoryId"`
}
