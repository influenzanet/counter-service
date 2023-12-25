package service

import (
	"context"
	"math"

	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/registry"
	"github.com/influenzanet/counter-service/pkg/stats"
	"github.com/influenzanet/counter-service/pkg/types"
)

type CounterService struct {
	definition types.CounterServiceDefinition
	registry   types.RegistryService
	collector  types.CollectorService
	channel    chan []types.Metric
}

func NewCounterService(dbService *db.StudyDBService, instance string, def types.CounterServiceDefinition) *CounterService {

	statSvc := stats.NewStatService(dbService, instance, def)
	registrySvc := registry.NewRegistryService(def.Name)
	channel := make(chan []types.Metric, 3)

	return &CounterService{
		definition: def,
		registry:   registrySvc,
		collector:  statSvc,
		channel:    channel,
	}
}

func (c *CounterService) Start(ctx context.Context) {
	go c.registry.Handle(ctx, c.channel)
	go c.collector.Run(ctx, c.channel)
}

func (c *CounterService) Registry() types.RegistryService {
	return c.registry
}

func (c *CounterService) Name() string {
	return c.definition.Name
}

func (c *CounterService) IsPublic() bool {
	return c.definition.Public
}

func (c *CounterService) IsRoot() bool {
	return c.definition.Root
}

func (c *CounterService) Data() types.CounterData {
	data := c.registry.Read()
	update := int64(math.Round(c.definition.UpdateDelay.Seconds()))
	return types.CounterData{
		UpdateDelay: update,
		Metrics:     data,
	}
}

func (c *CounterService) Definition() types.CounterServiceDefinition {
	return c.definition
}
