# sdk-velix-go — Go SDK ![version](https://img.shields.io/badge/version-0.1.0--alpha.1-orange)

> ⚠️ **Alpha / pre-release**, mas já publicado e confirmado funcionando de ponta a ponta contra a API real de staging (onboarding, checkin, LGPD, me, events). **pkg.go.dev:** https://pkg.go.dev/github.com/VELIX-Biometrics/sdk-velix-go

Official Go SDK for the VELIX Biometrics platform — facial access control B2B SaaS.

## Requirements

- Go 1.22+
- Zero external dependencies (stdlib only)

## Installation

```bash
go get github.com/VELIX-Biometrics/sdk-velix-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    velix "github.com/VELIX-Biometrics/sdk-velix-go"
)

func main() {
    ctx := context.Background()
    client := velix.NewClient(velix.Config{
        APIURL: "https://api.velixbiometrics.com",
        APIKey: os.Getenv("VELIX_API_KEY"),
    })

    result, err := client.Checkin.Identify(ctx, velix.CheckinIdentifyRequest{ImageBase64: frameBase64})
    if err != nil {
        panic(err)
    }
    fmt.Println(result.Match, result.SubjectID)
}
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `VELIX_API_URL` | Yes | API base URL (`https://api.velixbiometrics.com`) |
| `VELIX_API_KEY` | Yes | Tenant API key (`vx_live_...`) |

```go
client := velix.NewClient(velix.Config{
    APIURL: os.Getenv("VELIX_API_URL"),
    APIKey: os.Getenv("VELIX_API_KEY"),
})
```

## Modules

| Module | Methods | Endpoint |
|--------|---------|----------|
| `client.Onboarding` | `Create(ctx, req)` | `POST /v1/api/onboarding` (escopo `onboarding:write`) |
| `client.Checkin` | `Identify(ctx, req)` | `POST /v1/api/checkin/identify` (escopo `checkin:write`) |
| `client.LGPD` | `DeletionRequest(ctx, personID)` | `POST /v1/api/deletion-request` (escopo `lgpd:write`) |
| `client.Me` | `Get(ctx, personID)` | `GET /v1/api/me/{personId}` (escopo `me:read`) |
| `client.Events` | `CreateGuest(ctx, eventID, req)`, `GetGuest(ctx, eventID, guestID)` | `POST`/`GET /v1/api/events/{id}/guests` (escopos `events:write`/`events:read`) |

`client.Time` existe mas retorna erro — `api-velix-time` ainda não tem proxy público via BFF.

| `client.Contexts` | `Create/Get/List/Update/Remove(ctx, ...)`, `Authorize(ctx, contextID, payload)`, `ListAuthorizationDecisions`, `CreateLinkRequest` | `/v1/contexts/*` (BearerAuth) |
| `client.Memberships` | `Create`, `ListByContext`, `ListByIdentity`, `UpdateStatus`, `AddRoles`, `RemoveRoles` | `/v1/contexts/:id/memberships`, `/v1/identities/:id/memberships`, `/v1/memberships/*` |
| `client.ContextRoles` | `Create`, `List`, `LinkPermissions` | `/v1/context-roles*` |
| `client.ContextPermissions` | `Create`, `List` | `/v1/context-permissions` |
| `client.AuthorizationTokens` | `Validate` | `POST /v1/authorization-tokens/validate` |

## Identity Context

```go
ctxData, _ := client.Contexts.Create(ctx, map[string]any{"name": "Matriz SP", "contextType": "location"})
decision, _ := client.Contexts.Authorize(ctx, contextID, map[string]any{
    "identityId": "identity-uuid",
    "permission": "access:enter",
})
membership, _ := client.Memberships.Create(ctx, contextID, map[string]any{
    "identityId": "identity-uuid",
    "roleIds":    []string{"role-uuid"},
})
// saída de contexto (definitiva, sem carência)
_, _ = client.Memberships.UpdateStatus(ctx, membershipID, "revoked")
// vínculo cross-tenant — fica PENDING até a pessoa consentir via magic link
_, _ = client.Contexts.CreateLinkRequest(ctx, contextID, map[string]any{"identityId": "identity-uuid"})
_, _ = client.AuthorizationTokens.Validate(ctx, "vat_...", false)
```

## Onboarding Module

```go
result, err := client.Onboarding.Create(ctx, velix.OnboardingRequest{
    Name:         "João Silva",
    DocumentType: "CPF",
    Document:     "00000000000",
    Frames:       []string{frame1, frame2, frame3}, // mínimo 1, máximo 5
})
// result.PersonID, result.IdentityID, result.Enrolled, result.FramesResults
```

## Checkin Module

```go
result, err := client.Checkin.Identify(ctx, velix.CheckinIdentifyRequest{ImageBase64: frameBase64})
// result.Match, result.SubjectID, result.SubjectName, result.Liveness.OK, result.Model
```

Score de similaridade e de liveness nunca são expostos — apenas os booleanos `Match`/`Liveness.OK`.

## LGPD Module

```go
result, err := client.LGPD.DeletionRequest(ctx, personID)
// result.ProtocolNumber
```

## Me Module

```go
person, err := client.Me.Get(ctx, personID)
```

## Events Module

```go
guest, err := client.Events.CreateGuest(ctx, eventID, velix.CreateGuestRequest{
    Name:  "João Silva",
    Email: "joao@company.com",
})
fetched, err := client.Events.GetGuest(ctx, eventID, guest.ID)
```

## Error Handling

```go
import "errors"

result, err := client.Checkin.Identify(ctx, velix.CheckinIdentifyRequest{ImageBase64: frame})
if err != nil {
    var authErr *velix.AuthError
    var bioErr  *velix.BiometricError
    var rlErr   *velix.RateLimitError

    switch {
    case errors.As(err, &authErr):
        fmt.Println("Invalid API key")
    case errors.As(err, &bioErr):
        fmt.Println("Face not recognized or liveness failed")
    case errors.As(err, &rlErr):
        fmt.Printf("Rate limit — retry after %v\n", rlErr.RetryAfter)
    default:
        fmt.Println("Unexpected error:", err)
    }
}
```

## Running Tests

```bash
go test ./...
go test ./... -v          # verbose
go test ./... -cover      # with coverage
```

## Local Development

```bash
git clone https://github.com/VELIX-Biometrics/sdk-velix-go.git
cd sdk-velix-go
go build ./...
go vet ./...
go test ./...
```

## Get an API Key

Talk to the Velix team during onboarding — there is no self-service sandbox.
