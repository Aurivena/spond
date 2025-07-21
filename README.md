<p align="center">
  <img src="assets/logo.png" alt="Spond Logo" width="160" height="160"/>
</p>

# Spond
[![Go Reference](https://pkg.go.dev/badge/github.com/Aurivena/spond.svg)](https://pkg.go.dev/github.com/Aurivena/spond)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Description

- It is a compact library for standardized communication between the server (in Go) and web clients via JSON responses.
- Solves the problem of unified success/error structures for the API, custom response codes, minimalistic logging, and unified error handling.

---

### Install

```bash
go get github.com/Aurivena/spond@v1.0.2
```

---

## Opportunities

- A single JSON response format (success and error)
- Easy expansion of the list of status codes and messages
- Integration with Gin (or any other web framework)
- Built-in thread-safe logger (rotation, color output)
- Minimum dependencies, pure Go-style

---

## Usage example
### Work with API-output

```go
import "spond"

sp := spond.NewSpond()
// Success output
sp.SendResponseSuccess(c, spond.Success, map[string]string{"foo": "bar"})
// Error
sp.SendResponseError(c, sp.BuildError(spond.BadRequest, "Error", "incorect data"))
````

### Logger (rotation, color, thread-safe)

```go
import "spond/log"

logger := log.NewLog("log/io.log", 50*1024*1024) // 50 МБ
defer logger.Close()

logger.Info("Test INFO %s", "hello")
logger.Error("Some error: %v", err)
```

## Extension
### Append new code output

```go
import "spond/log"

sp := spond.NewSpond()
err := sp.AppendCode(7777, "Мой статус")
if err != nil {
    panic(err)
}
```
## Project Structure

```
spond/
  ├── log/          # Minimalistic logger (rotation, color, thread-safe)
  ├── response/     # Response structs, status codes, Success/Error types
  └── spond.go      # Core API handler logic
```
## Testing

```bash
go test ./...
```
Unit tests are available for all main features.
