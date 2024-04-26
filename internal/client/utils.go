package client

import (
	"encoding/json"
	"fmt"
)

func marshalAndJoin(left, right interface{}) ([]byte, error) {
	leftMarshaled, err := json.Marshal(left)
	if err != nil {
		return nil, err
	}

	if right == nil {
		return leftMarshaled, nil
	}

	rightMarshaled, err := json.Marshal(right)
	if err != nil {
		return nil, err
	}

	joined := fmt.Sprintf("%s,%s", leftMarshaled[:len(leftMarshaled)-1], rightMarshaled[1:])
	return []byte(joined), nil
}

type status struct {
	StatusValue string `json:"status"`
	Message     string `json:"message,omitempty"`
}

func (s *status) failed() bool {
	return s.StatusValue != "SUCCESS"
}

func (s status) Error() string {
	return fmt.Sprintf("%s: %s", s.StatusValue, s.Message)
}
