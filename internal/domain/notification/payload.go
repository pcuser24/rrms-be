package notification

import (
	"encoding/json"
)

type Notification struct {
	UserId  useridType      `json:"userId"`
	Payload json.RawMessage `json:"payload"`
}
