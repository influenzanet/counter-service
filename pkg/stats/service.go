package stats

import (
	"fmt"
	"time"
	"context"
	"errors"
	"github.com/coneno/logger"
	"github.com/influenzanet/counter-service/pkg/db"
	"github.com/influenzanet/counter-service/pkg/types"
)

var(
	ErrPreflightError = errors.New("Preflight check failed")
)

type StatsService struct {
	dbService  *db.StudyDBService
	instanceID	string
	studyKey string
	delay		time.Duration
	collectors []StatCollector
	filter types.StatFilter
}

func NewStatService(dbService *db.StudyDBService, instanceID string, studyKey string, delay time.Duration) *StatsService {

	collectors := []StatCollector{
		&ParticipantActiveCollector{},
		&ParticipantIntakeCollector{},
	}
	return &StatsService{dbService: dbService, collectors: collectors, instanceID: instanceID, studyKey:studyKey, delay: delay}
}

func (s *StatsService) WithFilter(filter types.StatFilter) {
	s.filter = filter
}

func (s *StatsService) preFlight() bool {
	ok, err := s.dbService.CheckDB(s.instanceID)
	if(err != nil) {
		logger.Error.Printf("Unable to check db existence : %s", err)
		return false
	}
	if(!ok) {
		return false
	}
	ok, err = s.dbService.CheckParticipantCollection(s.instanceID, s.studyKey)
	if(err != nil) {
		logger.Error.Printf("Unable to check db existence : %s", err)
		return false
	}
	if(!ok) {
		return false
	}
	return true
}

func (s *StatsService) Fetch() ([]types.Counter, error) {

	var counters []types.Counter
	
	if(!s.preFlight()) {
		return nil, ErrPreflightError
	}

	counters = make([]types.Counter, 0, len(s.collectors))

	for index, collector := range s.collectors {
		counter, err := collector.Fetch(s.dbService, s.instanceID, s.studyKey, s.filter)
		if err != nil {
			logger.Error.Printf("Error for counter %d : %s", index, err)
		} else {
			counter.Time = time.Now()
			counters = append(counters, counter)
		}
	}
	return counters, nil
}

func (s *StatsService) Run(ctx context.Context, out chan<-[]types.Counter) error {

		if s.delay < time.Second {
			return fmt.Errorf("Service delay must be at least one second")
		}
	
		for {
			logger.Info.Printf("Fetching %s:%s %v", s.instanceID, s.studyKey, s.delay)
			res, err := s.Fetch()
			if err != nil {
				logger.Error.Printf("Error during fetch %s:%s : %s", s.instanceID, s.studyKey, err)
			} else {
				out <- res
			}
			
			select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(s.delay):
					// Ok
			}
		}
		return nil
	
}