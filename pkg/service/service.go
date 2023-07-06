package service

import(
	"github.com/influenzanet/counter-service/pkg/types"
	"github.com/influenzanet/counter-service/pkg/stats"
	"github.com/influenzanet/counter-service/pkg/registry"
	"github.com/influenzanet/counter-service/pkg/db"
	"time"
	"context"
)

type CounterService struct {
	registry types.RegistryService
	collector types.CollectorService
	channel chan []types.Counter
}

func NewCounterService(dbService *db.StudyDBService, instance string, studyKey string, delay time.Duration, filter types.StatFilter) *CounterService {

	statSvc := stats.NewStatService(dbService, instance, studyKey, delay)
	statSvc.WithFilter(filter)
	registrySvc := registry.NewRegistryService(studyKey)
	channel := make(chan []types.Counter, 3)

	return &CounterService{
		registry: registrySvc,
		collector: statSvc,
		channel: channel,
	}
}


func (c *CounterService) Start(ctx context.Context) {
	go c.registry.Handle(ctx, c.channel)
	go c.collector.Run(ctx, c.channel)
}

func (c *CounterService) Registry() types.RegistryService {
	return c.registry
}
