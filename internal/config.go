package internal

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/influenzanet/counter-service/pkg/types"
	"github.com/influenzanet/go-utils/pkg/configs"
)

const ENV_USE_NO_CURSOR_TIMEOUT = "USE_NO_CURSOR_TIMEOUT"

func LoadConfig() types.ServiceConfig {

	dbConfig := configs.GetMongoDBConfig("STUDY_")

	InstanceID := configs.RequireEnv("INSTANCE_ID")

	PlatformCode := os.Getenv("PLATFORM")

	port := configs.GetEnvInt("PORT", 0)

	// Shorthand configuration for influenzanet counters
	InfluenzanetCounter := os.Getenv("INFLUENZANET_COUNTER")

	definitions := make([]types.CounterServiceDefinition, 0, 1)

	if InfluenzanetCounter != "" {
		ifnDef, err := GetInfluenzanetStudy(InfluenzanetCounter)
		if err != nil {
			log.Fatalf("Error reading INFLUENZANET_COUNTER: %s", err)
		}
		definitions = append(definitions, ifnDef)
	}

	extraCountersEnv := os.Getenv("EXTRA_COUNTERS")
	if extraCountersEnv != "" {
		dd, err := parseDefinitions([]byte(extraCountersEnv))
		if err != nil {
			log.Fatalf("Error reading EXTRA_COUNTERS in %s: %s", extraCountersEnv, err)
		}
		definitions = append(definitions, dd...)
	}

	extraCountersFile := os.Getenv("EXTRA_COUNTERS_FILE")

	if extraCountersFile != "" {
		dd, err := GetDefinitionsFromFile(extraCountersFile)
		if err != nil {
			log.Fatalf("Error reading EXTRA_COUNTERS_FILE in %s: %s", extraCountersFile, err)
		}
		definitions = append(definitions, dd...)
	}

	metaAuthKey := os.Getenv("META_AUTH_KEY")

	return types.ServiceConfig{
		StudyDBConfig:  dbConfig,
		Port:           port,
		MetaAuthKey:    metaAuthKey,
		InstanceID:     InstanceID,
		StatDefinition: definitions,
		Platform:       PlatformCode,
	}

}

func GetInfluenzanetStudy(envString string) (types.CounterServiceDefinition, error) {
	def := types.CounterServiceDefinition{
		Name:                     "influenzanet",
		Root:                     true,
		Public:                   true,
		ActiveParticipantSurveys: []string{"intake", "weekly", "vaccination"},
		ActiveParticipantDelay:   types.Duration{Duration: time.Hour * 546 * 24}, // Default is 18 month to count participant as active,
		UpdateDelay:              types.Duration{Duration: time.Hour * 24},
	}

	if strings.Contains(envString, "{") {
		err := json.Unmarshal([]byte(envString), &def)
		if err != nil {
			return def, nil
		}
	} else {
		// If provided value is only a string, then it's considered to be the StudyKey
		def.StudyKey = envString
	}

	return def, nil
}

func GetDefinitionsFromFile(extraStudiesFile string) ([]types.CounterServiceDefinition, error) {
	b, err := os.ReadFile(extraStudiesFile)
	if err != nil {
		return nil, err
	}
	return parseDefinitions(b)
}

func parseDefinitions(data []byte) ([]types.CounterServiceDefinition, error) {
	defs := make([]types.CounterServiceDefinition, 0)
	err := json.Unmarshal(data, &defs)
	if err != nil {
		return nil, err
	}
	return defs, nil
}
