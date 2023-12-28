package types

import (
	"github.com/influenzanet/go-utils/pkg/configs"
)

// ServiceConfig holds parsed configuration from env & files
type ServiceConfig struct {
	StudyDBConfig  configs.DBConfig
	Port           int    // Http Port to listen to
	MetaAuthKey    string // Auth key to access meta.json from server (will not be used if empty)
	InstanceID     string // Name of this instance
	Platform       string // Code of the platform
	StatDefinition []CounterServiceDefinition
}

// Definition for a counter service (provides a list of counter values)
type CounterServiceDefinition struct {
	StudyKey                 string   `json:"studykey"`       // Name of study to use (not shown in output)
	ActiveParticipantSurveys []string `json:"active_surveys"` // List of survey to use for to count active participants
	Root                     bool     `json:"root"`           // If true counter values are shown in root of the service
	Public                   bool     `json:"public"`         // If true lister in studies available in root of the service
	Name                     string   `json:"name"`           // Name of the counter service shown in service output
	ActiveParticipantDelay   Duration `json:"active_delay"`   // Delay to count participant as active
	UpdateDelay              Duration `json:"update_delay"`   // Update period of the counter (1/frequency)
}
