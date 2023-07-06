package registry

import (

	"context"
	"sync"
	"github.com/coneno/logger"
	"github.com/influenzanet/counter-service/pkg/types"
)

type RegistryService struct {
	studyKey	string
	counters	map[string]types.Counter
	mu       sync.Mutex
}

func NewRegistryService(studyKey string) *RegistryService {
	counters := make(map[string]types.Counter, 0)
	return &RegistryService{studyKey: studyKey, counters:counters}
}

func (r *RegistryService) Handle(ctx context.Context, input <-chan[]types.Counter ) error {
	for{
		select {
		case <-ctx.Done():
			return ctx.Err()

		case result := <-input:
			logger.Info.Printf("Update study '%s'", r.studyKey)
			r.updateCounters(result)
		}	
	}
}

func (r *RegistryService) updateCounters(result []types.Counter) {
	r.mu.Lock()
	for _, res := range result {
		r.counters[res.Name] = res
	}
	r.mu.Unlock()
}

func (r *RegistryService) Read() []types.Counter {
	r.mu.Lock()
	cc := make([]types.Counter,0, len(r.counters))
	for _, c := range r.counters {
		cc = append(cc, c)
	}
	r.mu.Unlock()
	return cc
} 