package core

import (
	"fmt"

	"github.com/Aurivena/spond/envelope"
)

func (s *Spond) codeExists(code envelope.StatusCode) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.statusMessages[code]
	return exist
}

// validate checks the length of the title and message.
// Returns the error text when restrictions are violated.
func validate(title, message string) error {
	if len(title) == 0 || len(title) > envelope.MaxTitleLength {
		return fmt.Errorf("%w", envelope.TitleInvalid)
	}
	if len(message) == 0 || len(message) > envelope.MaxMessageLength {
		return fmt.Errorf("%w", envelope.MessageInvalid)
	}
	return nil
}
