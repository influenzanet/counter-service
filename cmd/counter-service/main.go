package main

import (
	"context"
	"log"

	"github.com/influenzanet/counter-service/internal"
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/server"
	"github.com/influenzanet/counter-service/pkg/service"
	"github.com/influenzanet/counter-service/pkg/types"
)

const FormatDateOnly = "2006-01-02"

func main() {

	config := internal.LoadConfig()

	dbService := db.NewStudyDBService(config.StudyDBConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	counters := make(map[string]types.CounterService)

	for _, def := range config.StatDefinition {
		svc := service.NewCounterService(dbService, config.InstanceID, def)

		_, found := counters[def.Name]
		if found {
			log.Fatalf("Counter '%s' is already defined, counter name must be unique", def.Name)
		}

		counters[def.Name] = svc
		svc.Start(ctx)
	}

	httpServer := server.NewHTTPServer(config, counters)
	httpServer.Start()
}
