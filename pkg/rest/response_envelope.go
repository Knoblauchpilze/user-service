package rest

import (
	"encoding/json"

	"github.com/google/uuid"
)

type responseEnvelope struct {
	RequestId uuid.UUID       `json:"requestId"`
	Status    string          `json:"status"`
	Details   json.RawMessage `json:"details,omitempty"`
}
