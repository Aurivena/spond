package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Aurivena/spond/v2/envelope"
)

// codeExists check is there status
func (s *Spond) codeExists(code envelope.StatusCode) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.statusMessages[code]
	return exist
}

// validate checks the length of the title and message.
// Returns the error when restrictions are violated.
func validate(title, message string) error {
	if len(title) == 0 || len(title) > envelope.MaxTitleLength {
		return fmt.Errorf("%w", envelope.TitleInvalid)
	}
	if len(message) == 0 || len(message) > envelope.MaxMessageLength {
		return fmt.Errorf("%w", envelope.MessageInvalid)
	}
	return nil
}

// write encodes response as JSON and sends it to client.
// Always sets Content-Type to application/json; charset=utf-8.
func write(w http.ResponseWriter, output any, code envelope.StatusCode) {
	// set data for future json
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(output); err != nil {
		// fallback: plain text error if JSON encoding fails
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(code))
	w.Write(buff.Bytes())
}
