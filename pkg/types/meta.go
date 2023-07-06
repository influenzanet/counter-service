package types

import(
	"time"
)

type Meta struct {
	Studies []string `json:"studies"` // List of available studies with counters
	InfluenzanetStudy string `json:"influenzanet"` // Name of the study for influenzanet counter
	From    time.Time `json:"from"`
}