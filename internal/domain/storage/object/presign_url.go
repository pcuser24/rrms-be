package object

import (
	"net/http"
	"time"
)

type PresignURL struct {
	URL          string        `json:"url"`
	Method       string        `json:"method"`
	SignedHeader http.Header   `json:"signedHeader"`
	LifeTime     time.Duration `json:"lifeTime"`
}
