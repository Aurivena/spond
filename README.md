<p align="center">
  <img src="assets/logo.png" alt="Spond Logo" width="160" height="160"/>
</p>

# Spond
[![Go Reference](https://pkg.go.dev/badge/github.com/Aurivena/spond.svg)](https://pkg.go.dev/github.com/Aurivena/spond)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aurivena/spond)](https://goreportcard.com/report/github.com/Aurivena/spond)

## Description

- It is a compact library for standardized communication between the server (in Go) and web clients via JSON responses.
- Solves the problem of unified success/error structures for the API, custom response codes and unified error handling.

---

### Install

```bash
go get github.com/Aurivena/spond@v2.0.3
```

---

## Opportunities

- A single JSON response format (success and error)
- Easy expansion of the list of status codes and messages
- Minimum dependencies, pure Go-style

---

## Usage example
### Work with API-output

```go
import "github.com/Aurivena/spond/v2/core"

sp := spond.NewSpond()
// Success output
sp.SendResponseSuccess(w, spond.Success, map[string]string{"foo": "bar"})
// Error
sp.SendResponseError(c, sp.BuildError(spond.BadRequest, "Error", "incorect data","Change pls their input data"))
````

## Extension
### Append new code output

```go
import "github.com/Aurivena/spond/v2/core"

sp := spond.NewSpond()
err := sp.AppendCode(7777, "Мой статус")
if err != nil {
    panic(err)
}
```
## Project Structure

```
spond/
├── core/        # Core logic: helpers, response builders, encoders
└── envelope/    # Domain-level error and status handling
```
## Testing

```bash
go test ./...
```
Unit tests are available for all main features.
