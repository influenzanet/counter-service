package stats

import (
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/types"
)

type StatCollector interface {
	Fetch(db *db.StudyDBService, instanceID string, studyKey string, filter types.StatFilter) (types.Counter, error)
}


func participantCounter(name string, dbService *db.StudyDBService, instanceID string, studyKey string, filter types.StatFilter, options db.ParticipantOptions) (types.Counter, error) {
	counter := types.Counter{Name: name, Type: types.COUNTER_TYPE_COUNT}
	count, err := dbService.CountParticipants(instanceID, studyKey, filter, options)
	if err != nil {
		return counter, err
	}
	counter.Value = count
	return counter, nil
}

// ParticipantActiveCollector collect count of participants with active status in the study (No time frame)
type ParticipantActiveCollector struct {
}

func (u *ParticipantActiveCollector) Fetch(dbService *db.StudyDBService, instanceID string,  studyKey string, filter types.StatFilter) (types.Counter, error) {
	return participantCounter("participants_active", dbService, instanceID, studyKey, filter, db.ParticipantOptions{ActiveStatus: true})
}

// ParticipantIntakeCollector collect count of participants with active status in the study and an intake survey submitted in the filter time frame
type ParticipantIntakeCollector struct {
}

func (u *ParticipantIntakeCollector) Fetch(dbService *db.StudyDBService, instanceID string,  studyKey string, filter types.StatFilter) (types.Counter, error) {
	return participantCounter("participants_intake", dbService, instanceID, studyKey, filter, db.ParticipantOptions{ActiveStatus: true, SelectOn: db.SurveySubmission, SurveyKey: "intake"})
}
