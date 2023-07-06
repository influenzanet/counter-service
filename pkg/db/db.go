package db

import (
	"context"
	"time"
	"github.com/influenzanet/counter-service/pkg/types"
	"github.com/influenzanet/study-service/pkg/dbs/studydb"
	models "github.com/influenzanet/study-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudyDBService struct {
	*studydb.StudyDBService
	timeout int
}

func NewStudyDBService(configs models.DBConfig) *StudyDBService {

	return &StudyDBService{
		StudyDBService: studydb.NewStudyDBService(configs),
		timeout:       configs.Timeout,
	}
}

func (dbService *StudyDBService) dbName(instanceID string) string {
	return dbService.DBNamePrefix + instanceID + "_studyDB"
}

func (dbService *StudyDBService) participantName(studyKey string) string {
	return studyKey + "_participants"
}

func inArray(s string, a []string) bool {
	for _, n := range a {
		if(n == s) {
			return true
		}
	}
	return false
}

func (dbService *StudyDBService) CheckDB(instanceID string) (bool, error) {
	ctx, cancel := dbService.getContext()
	defer cancel()
	dbName := dbService.dbName(instanceID)
	nn, err := dbService.DBClient.ListDatabaseNames(ctx, bson.D{{Key: "name", Value: dbName}})
	if(err != nil) {
		return false, err
	}
	return inArray(dbName, nn), nil
}

func (dbService *StudyDBService) CheckCollection(instanceID string, collection string) (bool, error) {
	dbName := dbService.dbName(instanceID)
	db := dbService.DBClient.Database(dbName)
	ctx, cancel := dbService.getContext()
	defer cancel()
	nn, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: collection}})
	if(err != nil) {
		return false, err
	}
	return inArray(collection, nn), nil
}

func (dbService *StudyDBService) CheckParticipantCollection(instanceID string, studyKey string) (bool, error) {
	coll := dbService.participantName(studyKey)
	return dbService.CheckCollection(instanceID, coll)
}

func (dbService *StudyDBService) collectionRefParticipants(instanceID string, studyKey string) *mongo.Collection {
	dbName := dbService.dbName(instanceID)
	coll := dbService.participantName(studyKey)
	return dbService.DBClient.Database(dbName).Collection(coll)
}

// DB utils
func (dbService *StudyDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
}

func filterField(field string, filter types.StatFilter) interface{} {

	if filter.From == 0 && filter.Until == 0 {
		return nil
	}

	var criteria interface{}

	if filter.From > 0 && filter.Until == 0 {
		criteria = bson.D{{"$gt", filter.From}}
	}

	if filter.Until > 0 && filter.From == 0 {
		criteria = bson.D{{"$lt", filter.Until}}
	}

	if filter.Until > 0 && filter.From > 0 {
		criteria = bson.M{"$and": bson.A{
			bson.D{{"$gt", filter.From}},
			bson.D{{"$lt", filter.Until}},
		},
		}
	}
	return bson.D{{field, criteria}}
}

type ParticipantSelector int 

const (
	NoSelection ParticipantSelector = 0
	StudyEntry  ParticipantSelector = 1
	SurveySubmission ParticipantSelector = 2
)


type ParticipantOptions struct {
	ActiveStatus           bool
	SelectOn			ParticipantSelector
	SurveyKey	string
}

func combineCriteria(cc []interface{}) interface{} {
	if len(cc) == 0 {
		return bson.D{}
	}
	if len(cc) == 1 {
		return cc[0]
	}
	a := make(bson.A, 0, len(cc))
	for _, c := range cc {
		a = append(a, c)
	}
	return bson.M{"$and": a}
}

func (svc *StudyDBService) CountParticipants(instanceID string, studyKey string, filter types.StatFilter, opts ParticipantOptions) (int64, error) {
	ctx, cancel := svc.getContext()
	defer cancel()

	users := svc.collectionRefParticipants(instanceID, studyKey)

	criteria := make([]interface{}, 0, 1)

	if opts.ActiveStatus {
		criteria = append(criteria, bson.D{{"studyStatus", "active"}})
	}

	if opts.SelectOn == SurveySubmission {
		field := "lastSubmission." + opts.SurveyKey
		criteria = append(criteria, filterField(field, filter))
	}

	if opts.SelectOn == StudyEntry {
		criteria = append(criteria, filterField("enteredAt", filter))
	}

	cc := combineCriteria(criteria)

	count, err := users.CountDocuments(ctx, cc)

	if err != nil {
		return 0, err
	}
	return count, nil
}

