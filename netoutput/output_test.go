package netoutput_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type MockGRPCStream struct {
	grpc.ServerStream
	SentData any
}

func (m *MockGRPCStream) SendMsg(msg any) error {
	m.SentData = msg
	return nil
}

func TestWriteHTTPSuccess(t *testing.T) {
	w := httptest.NewRecorder()

	payload := map[string]string{"foo": "bar"}
	netoutput.WriteHTTP(w, netsp.Response[map[string]string]{
		Code: netstatus.CodeSuccess,
		Data: payload,
	})

	assert.Equal(t, netstatus.CodeSuccess.HTTP(), w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
}

func TestWriteHTTPSuccess_NoContent(t *testing.T) {
	w := httptest.NewRecorder()

	netoutput.WriteHTTP(w, netsp.Response[any]{
		Code: netstatus.CodeNotFound,
		Data: nil,
	})

	assert.Equal(t, netstatus.CodeNotFound.HTTP(), w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Type"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestWriteHTTPError(t *testing.T) {
	w := httptest.NewRecorder()

	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	appErr := netsp.Response[netsp.ErrorDetail]{
		Code: netstatus.CodeBadRequest,
		Data: netsp.ErrorDetail{Title: errTitle, Message: errMessage, Solution: ""},
	}

	netoutput.WriteHTTP(w, appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var out netsp.ErrorDetail
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, errTitle, out.Title)
	assert.Equal(t, errMessage, out.Message)
	assert.Equal(t, "", out.Solution)
}

func TestWriteGRPC_Success(t *testing.T) {
	stream := &MockGRPCStream{}
	payload := "hello grpc"

	err := netoutput.WriteGRPC(stream, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: payload,
	})

	assert.NoError(t, err)
	assert.Equal(t, payload, stream.SentData)
}

func TestWriteGRPC_Error(t *testing.T) {
	stream := &MockGRPCStream{}

	err := netoutput.WriteGRPC(stream, netsp.Response[any]{
		Code: netstatus.CodeNotFound,
		Data: nil,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NotFound")
}

func TestWriteWS_Success(t *testing.T) {
	resp := netsp.Response[string]{
		Code: netstatus.CodeSuccess,
		Data: "ws message",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"Code":0`)
	assert.Contains(t, string(data), `"Data":"ws message"`)
}
