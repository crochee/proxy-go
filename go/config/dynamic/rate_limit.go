package dynamic

import "time"

type RateLimit struct {
	Every time.Duration `json:"every"`
	Burst int           `json:"burst"`
	Mode  int           `json:"mode"`
}
