package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"strconv"
	"context"
	"github.com/coneno/logger"
	"github.com/influenzanet/counter-service/internal"
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/types"
	"github.com/influenzanet/counter-service/pkg/server"
	"github.com/influenzanet/counter-service/pkg/service"
)

const FormatDateOnly = "2006-01-02"

func getStudies() []string {
	instances := make([]string, 0)
	for _, i := range strings.Split(os.Getenv("STUDIES"), ",") {
		instance := strings.TrimSpace(i)
		if(instance == "") {
			continue
		}
		instances = append(instances, instance)
	}
	return instances
}

func getDateFromEnv(name string) time.Time {
	d := os.Getenv(name)
	t, err := time.Parse(FormatDateOnly, d)
	if(err != nil) {
		panic(fmt.Sprintf("env %s must be a valid date", name))
	}
	return t
}

func getIntFromEnv(name string, defaultValue int) int {
	d := os.Getenv(name)
	if(d == "") {
		return defaultValue
	}
	i, err := strconv.Atoi(d)
	if(err != nil) {
		logger.Error.Printf("Invalid numeric value provided for env %s", name)
		return defaultValue
	}
	return i
}

func main() {

	instance := os.Getenv("INSTANCE_ID")

	studies := getStudies()
	if(len(studies) == 0) {
		panic("At least one study must be provided")
	}

	logger.Info.Printf("Handling studies : %v", studies)

	influenzanetStudy := os.Getenv("INFLUENZANET_STUDY")

	from := getDateFromEnv("FROM_DATE")

	delayMinutes := getIntFromEnv("UPDATE_DELAY", 60 * 5)

	port := getIntFromEnv("PORT", 0)

	config := internal.GetStudyDBConfig()

	dbService := db.NewStudyDBService(config)

	delay := time.Duration(delayMinutes) * time.Minute

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	registries := make(map[string]types.RegistryService)

	filter := types.StatFilter{From: from.Unix(),}

	for _, study := range(studies) {
		svc := service.NewCounterService(dbService, instance, study, delay, filter)
		registries[study] = svc.Registry()
		svc.Start(ctx)
	}

	meta := &types.Meta{Studies: studies, InfluenzanetStudy: influenzanetStudy, From: from}

	httpServer := server.NewHTTPServer(port, registries, meta)
	httpServer.Start()
}
