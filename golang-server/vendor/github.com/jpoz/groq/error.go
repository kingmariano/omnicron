package groq

import "fmt"

type Error struct {
	Message          string `json:"message"`
	Type             string `json:"type"`
	FailedGeneration string `json:"failed_generation,omitempty"`
}

func (e Error) Error() string {
	if e.FailedGeneration == "" {
		return fmt.Sprintf("%s: %s", e.Type, e.Message)
	}
	return fmt.Sprintf("%s: %s\n\nfailed generation:\n%s", e.Type, e.Message, e.FailedGeneration)
}
