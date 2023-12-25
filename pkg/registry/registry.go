package registry

import (
	"context"
	"sync"

	"github.com/coneno/logger"
	"github.com/influenzanet/counter-service/pkg/types"
)

type RegistryService struct {
	name    string
	metrics map[string]types.Metric
	mu      sync.Mutex
}

func NewRegistryService(studyKey string) *RegistryService {
	metrics := make(map[string]types.Metric, 0)
	return &RegistryService{name: studyKey, metrics: metrics}
}

func (r *RegistryService) Handle(ctx context.Context, input <-chan []types.Metric) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case result := <-input:
			logger.Info.Printf("Update study '%s'", r.name)
			r.updateCounters(result)
		}
	}
}

func (r *RegistryService) updateCounters(result []types.Metric) {
	r.mu.Lock()
	for _, res := range result {
		r.metrics[res.Name] = res
	}
	r.mu.Unlock()
}

func (r *RegistryService) Read() []types.Metric {
	r.mu.Lock()
	cc := make([]types.Metric, 0, len(r.metrics))
	for _, c := range r.metrics {
		cc = append(cc, c)
	}
	r.mu.Unlock()
	return cc
}
