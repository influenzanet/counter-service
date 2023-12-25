package stats

import (
	"time"

	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/types"
)

// StatCollector base interface for a counter collecting stats
type StatCollector interface {
	Fetch(db *db.StudyDBService, instanceID string, studyKey string) (types.Metric, error)
}

func participantCounter(name string, dbService *db.StudyDBService, instanceID string, studyKey string, filter types.StatFilter, options db.ParticipantOptions) (types.Metric, error) {
	counter := types.Metric{Name: name, Type: types.METRIC_TYPE_COUNT}
	count, err := dbService.CountParticipants(instanceID, studyKey, filter, options)
	if err != nil {
		return counter, err
	}
	counter.Value = count
	return counter, nil
}

// ParticipantEnrolledCollector collect count of participants with active status in the study (No time frame)
type ParticipantEnrolledCollector struct {
}

func (u *ParticipantEnrolledCollector) Fetch(dbService *db.StudyDBService, instanceID string, studyKey string) (types.Metric, error) {
	filter := types.StatFilter{}
	return participantCounter("participants_enrolled", dbService, instanceID, studyKey, filter, db.ParticipantOptions{ActiveStatus: true})
}

// ParticipantActiveCollector collect count of participants with active status in the study and at least one survey submitted in the filter time frame
type ParticipantActiveCollector struct {
	SurveyKeys  []string
	ActiveDelay time.Duration
}

func (u *ParticipantActiveCollector) Fetch(dbService *db.StudyDBService, instanceID string, studyKey string) (types.Metric, error) {
	from := time.Now().Add(-u.ActiveDelay)
	filter := types.StatFilter{From: from.Unix()}
	return participantCounter("participants_active", dbService, instanceID, studyKey, filter, db.ParticipantOptions{ActiveStatus: true, SelectOn: db.SurveySubmission, SurveyKeys: u.SurveyKeys})
}
