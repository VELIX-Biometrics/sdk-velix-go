package velix

import "fmt"

// VelixError erro base do SDK com código HTTP e mensagem da API.
type VelixError struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *VelixError) Error() string {
	return fmt.Sprintf("velix: %s (status %d, code %s)", e.Message, e.StatusCode, e.Code)
}

// AuthError credenciais inválidas ou expiradas.
type AuthError struct{ VelixError }

// BiometricError falha de identificação ou enroll biométrico.
type BiometricError struct{ VelixError }

// NotFoundError recurso não encontrado.
type NotFoundError struct{ VelixError }

// RateLimitError requisição bloqueada por rate limit.
type RateLimitError struct{ VelixError }

func classifyError(statusCode int, message, code string) error {
	base := VelixError{StatusCode: statusCode, Message: message, Code: code}
	switch statusCode {
	case 401, 403:
		return &AuthError{base}
	case 404:
		return &NotFoundError{base}
	case 429:
		return &RateLimitError{base}
	}
	if code == "BIOMETRIC_FAILED" || code == "LIVENESS_FAILED" {
		return &BiometricError{base}
	}
	return &base
}
