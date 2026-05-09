package netoutput

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Write encodes response as JSON and sends it to client.
// Always sets Content-Type to application/json; charset=utf-8.
func WriteHTTP[T any](w http.ResponseWriter, resp netsp.Response[T]) {
	httpCode := resp.Code.HTTP()

	if httpCode == http.StatusNoContent || any(resp.Data) == nil {
		w.WriteHeader(int(httpCode))
		return
	}

	// set data for future jsos
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(resp.Data); err != nil {
		// fallback: plain text error if JSON encoding fails
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(httpCode))
	w.Write(buff.Bytes())
}

func WriteGRPC[T any](stream grpc.ServerStream, resp netsp.Response[T]) error {
	if err := stream.SendMsg(resp.Data); err != nil {
		return err
	}

	if resp.Code != netstatus.CodeSuccess {
		return status.Error(resp.Code.GRPC(), "request failed")
	}
	return nil
}

func WriteWS[T any](conn *websocket.Conn, resp netsp.Response[T]) error {
	return json.NewEncoder(conn).Encode(resp)
}
