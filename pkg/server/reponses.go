package server

import "github.com/influenzanet/counter-service/pkg/types"

// RootResponse is the JSON response sent at the root of the service
type RootResponse struct {
	Platform string                       `json:"platform"`
	Extra    []string                     `json:"extra,omitempty"`
	Counters map[string]types.CounterData `json:"counters"`
}
