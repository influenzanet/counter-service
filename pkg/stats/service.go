package stats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coneno/logger"
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/types"
)

var (
	ErrPreflightError = errors.New("preflight check failed")
)

type StatsService struct {
	dbService   *db.StudyDBService
	instanceID  string
	studyKey    string
	updateDelay time.Duration
	collectors  []StatCollector
}

func NewStatService(dbService *db.StudyDBService, instanceID string, def types.CounterServiceDefinition) *StatsService {

	collectors := []StatCollector{
		&ParticipantEnrolledCollector{},
		&ParticipantActiveCollector{SurveyKeys: def.ActiveParticipantSurveys, ActiveDelay: def.ActiveParticipantDelay.Duration},
	}
	return &StatsService{dbService: dbService, collectors: collectors, instanceID: instanceID, studyKey: def.StudyKey, updateDelay: def.UpdateDelay.Duration}
}

func (s *StatsService) preFlight() bool {
	ok, err := s.dbService.CheckDB(s.instanceID)
	if err != nil {
		logger.Error.Printf("Unable to check db existence : %s", err)
		return false
	}
	if !ok {
		return false
	}
	ok, err = s.dbService.CheckParticipantCollection(s.instanceID, s.studyKey)
	if err != nil {
		logger.Error.Printf("Unable to check db existence : %s", err)
		return false
	}
	if !ok {
		return false
	}
	return true
}

func (s *StatsService) Fetch() ([]types.Metric, error) {

	var metrics []types.Metric

	if !s.preFlight() {
		return nil, ErrPreflightError
	}

	metrics = make([]types.Metric, 0, len(s.collectors))

	for index, collector := range s.collectors {
		metric, err := collector.Fetch(s.dbService, s.instanceID, s.studyKey)
		if err != nil {
			logger.Error.Printf("Error for counter %d : %s", index, err)
		} else {
			metric.Time = time.Now()
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (s *StatsService) Run(ctx context.Context, out chan<- []types.Metric) error {

	if s.updateDelay < time.Second {
		return fmt.Errorf("service delay must be at least one second")
	}

	for {
		logger.Info.Printf("Fetching %s:%s %v", s.instanceID, s.studyKey, s.updateDelay)
		res, err := s.Fetch()
		if err != nil {
			logger.Error.Printf("Error during fetch %s:%s : %s", s.instanceID, s.studyKey, err)
		} else {
			out <- res
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.updateDelay):
			// Ok
		}
	}
	return nil

}
