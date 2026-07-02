# velix-sdk-go — Go SDK ![version](https://img.shields.io/badge/version-0.1.0--alpha.1-orange)

> ⚠️ **Alpha / pre-release.** This SDK targets a public API surface that does not yet fully exist on the VELIX backend (see internal task #593). Endpoints and auth may not work against production. Do not use in production integrations yet.

Official Go SDK for the VELIX Biometrics platform — facial access control B2B SaaS.

## Requirements

- Go 1.22+
- Zero external dependencies (stdlib only)

## Installation

```bash
go get github.com/velix-io/velix-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    velix "github.com/velix-io/velix-sdk-go"
)

func main() {
    client := velix.NewClient(velix.Config{
        APIURL: "https://api.velixbiometrics.com",
        APIKey: "vx_live_...",
    })

    result, err := client.Checkin.Facial(context.Background(), "tenant-slug", frameBase64)
    if err != nil {
        panic(err)
    }
    fmt.Println(result.Passed, result.PersonID)
}
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `VELIX_API_URL` | Yes | API base URL (`https://api.velixbiometrics.com`) |
| `VELIX_API_KEY` | Yes | Tenant API key (`vx_live_...` or `vx_sandbox_...`) |

```go
client := velix.NewClient(velix.Config{
    APIURL: os.Getenv("VELIX_API_URL"),
    APIKey: os.Getenv("VELIX_API_KEY"),
})
```

## Modules

| Module | Methods |
|--------|---------|
| `client.Checkin` | `Facial()`, `QR()`, `PIN()`, `GetHistory()` |
| `client.Persons` | `List()`, `Get()`, `Create()`, `Update()`, `Delete()`, `Enroll()` |
| `client.Events` | `List()`, `Get()`, `Create()`, `Configure()` |
| `client.Tenants` | `Me()`, `UpdateSettings()` |

## Checkin Module

```go
ctx := context.Background()

// Facial identification (base64 JPEG frame)
result, err := client.Checkin.Facial(ctx, "tenant-slug", frameBase64)
// result.Passed == true
// result.PersonID == "uuid"
// result.PersonName == "João Silva"

// QR code checkin
result, err := client.Checkin.QR(ctx, "tenant-slug", qrToken)

// PIN checkin
result, err := client.Checkin.PIN(ctx, "tenant-slug", pin)

// Paginated history
history, err := client.Checkin.GetHistory(ctx, "tenant-slug", &velix.ListOptions{Page: 1, Limit: 20})
```

## Persons Module

```go
// List with optional search
list, err := client.Persons.List(ctx, &velix.ListOptions{Page: 1, Limit: 20, Search: "João"})

// Get by ID
person, err := client.Persons.Get(ctx, "uuid")

// Create
created, err := client.Persons.Create(ctx, velix.CreatePersonInput{
    Name:       "João Silva",
    Email:      "joao@company.com",
    ExternalID: "EMP-001",
})

// Update
err = client.Persons.Update(ctx, "uuid", velix.UpdatePersonInput{Name: "João B. Silva"})

// Enroll biometrics (minimum 3 base64 frames)
err = client.Persons.Enroll(ctx, "uuid", []string{frame1, frame2, frame3})

// Delete
err = client.Persons.Delete(ctx, "uuid")
```

## Events Module

```go
list,    err := client.Events.List(ctx, &velix.ListOptions{Page: 1, Limit: 20})
event,   err := client.Events.Get(ctx, "uuid")
created, err := client.Events.Create(ctx, velix.CreateEventInput{Name: "Conference 2026"})
err = client.Events.Configure(ctx, "uuid", velix.EventConfig{CheckInOpen: true})
```

## Tenants Module

```go
tenant, err := client.Tenants.Me(ctx)
err = client.Tenants.UpdateSettings(ctx, velix.TenantSettings{RequireLiveness: true})
```

## Error Handling

```go
import "errors"

result, err := client.Checkin.Facial(ctx, "slug", frame)
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
git clone <repo>
cd velix-sdk-go
go build ./...
go vet ./...
go test ./...
```

## Get an API Key

Access the dashboard at **velixbiometrics.com** → Settings → API Keys → New Key.
