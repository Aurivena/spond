<p align="center">
  <img src="assets/logo.png" alt="Spond Logo" width="160" height="160"/>
</p>

# Spond
[![Go Reference](https://pkg.go.dev/badge/github.com/Aurivena/spond.svg)](https://pkg.go.dev/github.com/Aurivena/spond/v4)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aurivena/spond)](https://goreportcard.com/report/github.com/Aurivena/spond/v4)


---

### Install

```bash
go get github.com/Aurivena/spond/v4@v4.0.0
```

---

Возможности

- **Единый формат ответа**: Использование универсальной структуры `Response[T]` для всех типов вывода.
- **Поддержка нескольких протоколов**: Встроенные механизмы записи для HTTP, gRPC и WebSockets.
- **Маппинг статус-кодов**: Бесшовный переход между внутренними кодами библиотеки, HTTP-статусами и кодами gRPC.
- **Валидация ошибок**: Автоматическая проверка длины заголовков и сообщений об ошибках для поддержания консистентности API.
- **Минимум зависимостей**: Чистый стиль Go с минимальным количеством внешних библиотек.

---

## Примеры использования

### 1. HTTP ответы

```go
import (
    "net/http"
    "github.com/Aurivena/spond/v3/netoutput"
    "github.com/Aurivena/spond/v3/netsp"
    "github.com/Aurivena/spond/v3/netstatus"
)

// Успешный ответ
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"foo": "bar"}
    resp := netsp.Response[any]{
        Code: netstatus.CodeSuccess,
        Data: data,
    }
    netoutput.WriteHTTP(w, resp)
}

// Ответ с ошибкой
func errorHandler(w http.ResponseWriter, r *http.Request) {
    errDetail := netsp.ErrorDetail{
        Title:    "Invalid Input",
        Message:  "The provided data is incorrect",
        Solution: "Please check the documentation and try again",
    }
    // BuildError создает структуру Response с указанным кодом и данными ошибки
    resp := netsp.BuildError(netstatus.CodeBadRequest, errDetail)
    netoutput.WriteHTTP(w, *resp)
}
```

### 2. gRPC ответы

```go
import (
    "google.golang.org/grpc"
    "github.com/Aurivena/spond/v3/netoutput"
    "github.com/Aurivena/spond/v3/netsp"
    "github.com/Aurivena/spond/v3/netstatus"
)

func grpcHandler(stream grpc.ServerStream) error {
    data := "Hello World"
    resp := netsp.Response[string]{
        Code: netstatus.CodeSuccess,
        Data: data,
    }
    return netoutput.WriteGRPC(stream, resp)
}
```

### 3. WebSocket ответы

```go
import (
    "golang.org/x/net/websocket"
    "github.com/Aurivena/spond/v3/netoutput"
    "github.com/Aurivena/spond/v3/netsp"
    "github.com/Aurivena/spond/v3/netstatus"
)

func wsHandler(conn *websocket.Conn) error {
    resp := netsp.Response[string]{
        Code: netstatus.CodeSuccess,
        Data: "Connected!",
    }
    return netoutput.WriteWS(conn, resp)
}
```

---

## Структура проекта

```
spond/
├── netoutput/  # Транспорты: запись ответов для HTTP, gRPC и WebSocket
├── netsp/       # Ядро: структуры Response, построение ошибок и валидация
└── netstatus/  # Статус-коды и их маппинги (HTTP/gRPC)
```

---

## Тестирование

```bash
go test ./...
```
Юнит-тесты доступны для всех основных функций.
