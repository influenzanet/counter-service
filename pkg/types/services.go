package types

import (
	"context"
)

// RegistryService holds metrics values after fetch, used as cache between two update
// Metrics values are sent to registry through channel
type RegistryService interface {
	Handle(ctx context.Context, input <-chan []Metric) error
	Read() []Metric
}

// Collector a collector service compute metrics and send them to channel
type CollectorService interface {
	Run(ctx context.Context, out chan<- []Metric) error
}

// CounterService manages data collection of metrics from collector
type CounterService interface {
	Start(ctx context.Context)
	Registry() RegistryService
	Name() string
	IsPublic() bool
	IsRoot() bool
	Data() CounterData
	Definition() CounterServiceDefinition
}

// CounterData holds output data for a CounterService
type CounterData struct {
	UpdateDelay int64    `json:"period"` // Number of second between 2 updates
	Metrics     []Metric `json:"data"`
}
