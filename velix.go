// Package velix é o SDK oficial Go para a plataforma VELIX Biometrics.
package velix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second
const userAgent = "velix-go-sdk/0.1.0-alpha.1"

// Config opções de inicialização do VelixClient.
type Config struct {
	APIURL  string        // ex: "https://api.velixbiometrics.com"
	APIKey  string        // ex: "vx_live_..."
	JWT     string        // alternativa ao APIKey (Bearer token)
	Timeout time.Duration // default: 30s
}

// VelixClient cliente principal do SDK.
type VelixClient struct {
	cfg     Config
	http    *http.Client
	Checkin *CheckinModule
	Persons *PersonsModule
	Events  *EventsModule
	Tenants *TenantsModule
}

// NewClient cria um VelixClient configurado.
func NewClient(cfg Config) *VelixClient {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}
	c := &VelixClient{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
	c.Checkin = &CheckinModule{c}
	c.Persons = &PersonsModule{c}
	c.Events = &EventsModule{c}
	c.Tenants = &TenantsModule{c}
	return c
}

// do executa uma requisição autenticada com retry em 429/503.
func (c *VelixClient) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("velix: marshal body: %w", err)
		}
		bodyBytes = b
	}

	resp, err := withRetry(ctx, func() (*http.Response, error) {
		var bodyReader io.Reader
		if bodyBytes != nil {
			// corpo recriado a cada tentativa: bytes.Reader é consumido no
			// envio anterior e um retry com o mesmo reader enviaria corpo vazio.
			bodyReader = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.cfg.APIURL+path, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Content-Type", "application/json")
		if c.cfg.APIKey != "" {
			req.Header.Set("x-api-key", c.cfg.APIKey)
		} else if c.cfg.JWT != "" {
			req.Header.Set("Authorization", "Bearer "+c.cfg.JWT)
		}
		return c.http.Do(req)
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("velix: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}
		_ = json.Unmarshal(raw, &apiErr)
		return nil, classifyError(resp.StatusCode, apiErr.Message, apiErr.Code)
	}

	// unwrap envelope { data: T }
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err == nil && envelope.Data != nil {
		return envelope.Data, nil
	}
	return raw, nil
}
