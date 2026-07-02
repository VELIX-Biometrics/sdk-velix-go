package velix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PersonsModule CRUD de pessoas + enroll biométrico.
type PersonsModule struct{ c *VelixClient }

// List retorna lista paginada de pessoas do tenant.
func (m *PersonsModule) List(ctx context.Context, opts *ListOptions) (*ListResponse[Person], error) {
	path := "/v1/persons"
	if opts != nil {
		q := url.Values{}
		if opts.Page > 0 {
			q.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Search != "" {
			q.Set("search", opts.Search)
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}
	raw, err := m.c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[Person]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("velix: decode persons list: %w", err)
	}
	return &result, nil
}

// Get retorna uma pessoa pelo ID.
func (m *PersonsModule) Get(ctx context.Context, id string) (*Person, error) {
	raw, err := m.c.do(ctx, http.MethodGet, "/v1/persons/"+id, nil)
	if err != nil {
		return nil, err
	}
	var p Person
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("velix: decode person: %w", err)
	}
	return &p, nil
}

// Create cria uma nova pessoa no tenant.
func (m *PersonsModule) Create(ctx context.Context, input CreatePersonInput) (*Person, error) {
	raw, err := m.c.do(ctx, http.MethodPost, "/v1/persons", input)
	if err != nil {
		return nil, err
	}
	var p Person
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("velix: decode person: %w", err)
	}
	return &p, nil
}

// Update atualiza dados de uma pessoa.
func (m *PersonsModule) Update(ctx context.Context, id string, input UpdatePersonInput) (*Person, error) {
	raw, err := m.c.do(ctx, http.MethodPut, "/v1/persons/"+id, input)
	if err != nil {
		return nil, err
	}
	var p Person
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("velix: decode person: %w", err)
	}
	return &p, nil
}

// Delete remove uma pessoa do tenant.
func (m *PersonsModule) Delete(ctx context.Context, id string) error {
	_, err := m.c.do(ctx, http.MethodDelete, "/v1/persons/"+id, nil)
	return err
}

// Enroll cadastra biometria facial para uma pessoa (mínimo 3 frames).
func (m *PersonsModule) Enroll(ctx context.Context, id string, frames []string) error {
	_, err := m.c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/persons/%s/enroll", id), EnrollInput{Frames: frames})
	return err
}
