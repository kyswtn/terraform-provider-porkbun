package client

import (
	"fmt"
)

type Status struct {
	Value   string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (s Status) HasFailed() bool {
	return s.Value != "SUCCESS"
}

func (s Status) Error() string {
	return fmt.Sprintf("porkbun: %s: %s", s.Value, s.Message)
}
