<p align="center">
  <img src="assets/logo.png" alt="Spond Logo" width="160" height="160"/>
</p>

# Spond
[![Go Reference](https://pkg.go.dev/badge/github.com/Aurivena/spond.svg)](https://pkg.go.dev/github.com/Aurivena/spond/v3)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aurivena/spond)](https://goreportcard.com/report/github.com/Aurivena/spond/v3)


---

### Install

```bash
go get github.com/Aurivena/spond/v3@v3.0.1
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
import (
    "net/http"
    "github.com/Aurivena/spond/v3/netsp"
)

//Success response
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"foo": "bar"}
    netsp.SendResponseSuccess(w, http.StatusOK, data)
}

//Error response
func errorHandler(w http.ResponseWriter, r *http.Request) {
    // BuildError validates title and message lengths automatically
    err := netsp.BuildError(
        http.StatusBadRequest, 
        "Invalid Input", 
        "The provided data is incorrect", 
        "Please check the documentation and try again",
    )
    netsp.SendResponseError(w, err)
}
```

## Extension
### Append new code output

```go
import "github.com/Aurivena/spond/v3/netsp"

if err := netsp.AppendCode(418, "I'm a teapot"); err != nil {
    log.Fatal(err)
}
```
## Project Structure

```
spond/
├── netsp/        # Core logic: helpers, response builders, encoders
```
## Testing

```bash
go test ./...
```
Unit tests are available for all main features.
