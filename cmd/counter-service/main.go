package main

import (
	"context"
	"log"
	"os"

	"github.com/influenzanet/counter-service/internal"
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/server"
	"github.com/influenzanet/counter-service/pkg/service"
	"github.com/influenzanet/counter-service/pkg/types"
	"github.com/influenzanet/counter-service/pkg/version"
)

const FormatDateOnly = "2006-01-02"

func main() {

	// Only show version
	ShowVersion := false

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "version" {
			ShowVersion = true
		}
	}

	log.Printf("%s Version: %s (%s)", version.Name, version.Version, version.Revision)

	if ShowVersion {
		os.Exit(0)
	}

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
