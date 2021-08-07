package stores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mn "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	Database   string
	URI        string
	Client     *mn.Client
	Connect    func(context.Context) error
	Disconnect func(context.Context) error
}

func (m *Mongo) personCollection() string {
	return "person"
}

func NewMongoStore(uri, database string) (Mongo, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mn.NewClient(clientOpts)
	if err != nil {
		return Mongo{}, MongoConnectionErr{
			Reason: "failed to create mongo client",
			Inner:  err,
		}
	}
	return Mongo{
		Database:   database,
		URI:        uri,
		Client:     client,
		Connect:    client.Connect,
		Disconnect: client.Disconnect,
	}, nil
}

type Case struct {
	ReportingDate *time.Time `bson:"dateOfReporting" json:"reportingDate"`
}

func (m *Mongo) FindConfirmedCases(ctx context.Context, outbreakID string, reportingDate time.Time, endDate *time.Time) ([]Case, error) {
	collection := m.Client.Database(m.Database).Collection(m.personCollection())
	lastDate := endDate
	if lastDate == nil {
		l := reportingDate.Add(time.Hour * 24)
		lastDate = &l
	}
	filter := bson.M{
		"outbreakId":     outbreakID,
		"classification": "LNG_REFERENCE_DATA_CATEGORY_CASE_CLASSIFICATION_CONFIRMED",
		"deleted":        false,
		"$and": bson.A{
			bson.M{"dateOfReporting": bson.M{"$gte": reportingDate}},
			bson.M{"dateOfReporting": bson.M{"$lt": lastDate}},
		},
	}

	var cases []Case
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return cases, MongoQueryErr{
			Reason: fmt.Sprintf("failed to retrieve cases for outbreak %s on reporting date %v", outbreakID, reportingDate),
			Inner:  err,
		}
	}
	if err := cursor.All(ctx, &cases); err != nil {
		return cases, MongoQueryErr{
			Reason: fmt.Sprintf("error executing query for outbreak %s on reporting date %v", outbreakID, reportingDate),
			Inner:  err,
		}
	}

	return cases, nil
}

type CaseCount struct {
	ReportingDate *time.Time `bson:"_id" json:"reportingDate"`
	Count         int32      `bson:"count" json:"count"`
}

func (m *Mongo) GroupCasesByDate(ctx context.Context, outbreakID string, reportingDate time.Time) ([]CaseCount, error) {
	var cases []CaseCount
	collection := m.Client.Database(m.Database).Collection(m.personCollection())
	lastDate := reportingDate.Add(time.Hour * 24)

	matchStage := bson.D{
		{"$match", bson.D{ //nolint:govet
			{"outbreakId", outbreakID},
			{"classification", "LNG_REFERENCE_DATA_CATEGORY_CASE_CLASSIFICATION_CONFIRMED"},
			{"deleted", false},
			{"$and", bson.A{
				bson.M{"dateOfReporting": bson.M{"$gte": reportingDate}},
				bson.M{"dateOfReporting": bson.M{"$lt": lastDate}},
			}}},
		},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":   "$dateOfReporting",
			"count": bson.M{"$sum": 1},
		}},
	}
	cursor, err := collection.Aggregate(ctx, mn.Pipeline{matchStage, groupStage})
	if err != nil {
		return cases, MongoQueryErr{
			Reason: fmt.Sprintf("failed to retrieve cases for outbreak %s on reporting date %v", outbreakID, reportingDate),
			Inner:  err,
		}
	}
	if err := cursor.All(ctx, &cases); err != nil {
		return cases, MongoQueryErr{
			Reason: fmt.Sprintf("error executing query for outbreak %s on reporting date %v", outbreakID, reportingDate),
			Inner:  err,
		}
	}

	return cases, nil
}
