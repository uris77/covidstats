package stores

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mn "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// District is the district representation
type District string

const (
	bz District = "Belize"
	cy District = "Cayo"
	cz District = "Corozal"
	ow District = "Orange Walk"
	sc District = "Stann Creek"
	to District = "Toledo"
)

type godataDistrict struct {
	name District
	code string
}

// Mongo represents a mongo client
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

func (m *Mongo) locationCollection() string {
	return "location"
}

// NewMongoStore creates a new mongo connection
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

// Address represents an address
type Address struct {
	TypeID     string `json:"type_id" bson:"typeId"`
	LocationID string `json:"location_id" bson:"locationId"`
	ParentID   string `json:"parent_location_id" bson:"parentLocationId"`
}

// Case represents a COVID case
type Case struct {
	ReportingDate *time.Time `bson:"dateOfReporting" json:"reportingDate"`
	ResidenceID   string     `bson:"usualPlaceOfResidenceLocationId"`
	District      District   `json:"district"`
	Total         int        `json:"total"`
}

// Location represents a location in Belize
type Location struct {
	ID               string `bson:"_id"`
	ParentLocationID string `bson:"parentLocationId"`
}

// FindConfirmedCases finds confirmed cases for a given date range
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

// FindLocationByID retrieves locations that match the locationID in the case
func (m *Mongo) FindLocationByID(ctx context.Context, ID string) (Location, error) {
	var locs Location
	collection := m.Client.Database(m.Database).Collection(m.locationCollection())
	filter := bson.M{"_id": ID}
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return locs, MongoQueryErr{Reason: "location.FindOne failed", Inner: result.Err()}
	}
	var cs Location
	if err := result.Decode(&cs); err != nil {
		return locs, MongoQueryErr{
			Reason: "error decoding location",
			Inner:  err,
		}
	}

	return locs, nil
}

// AddDistrictToCase adds the district field to the cases
func (m *Mongo) AddDistrictToCase(ctx context.Context, cases []Case) ([]Case, error) {
	var cs []Case
	for _, c := range cases {
		loc, err := m.FindLocationByID(ctx, c.ResidenceID)
		if err != nil {
			return cs, MongoQueryErr{Reason: "error finding location", Inner: err}
		}
		c.District = findDistrictFromCode(loc.ParentLocationID)
		cs = append(cs, c)
	}
	return cs, nil
}

//type CasesByDistrict struct {
//	DateOfReporting *time.Time `json:"dateOfReporting"`
//	District        District   `json:"district"`
//	Total           int        `json:"total"`
//}

//func countCasesByDistrict(cases []Case) []CasesByDistrict {
//	var cs []CasesByDistrict
//	for _, c := range cases {
//		ex := findCaseByDistrictAndDate(cs, c.ReportingDate, c.District)
//		if ex != nil {
//			ex.Total = ex.Total + 1
//		} else {
//			ex.Total = 1
//			cs = append(cs, *ex)
//		}
//	}
//	return cs
//}

//func findCaseByDistrictAndDate(cases []CasesByDistrict, date *time.Time, district District) *CasesByDistrict {
//	for _, c := range cases {
//		if c.District == district && c.DateOfReporting == date {
//			return &c
//		}
//	}
//	return nil
//}

func findDistrictFromCode(code string) District {
	dist := bz
	districtsToCode := []godataDistrict{
		{
			name: bz,
			code: "bfc2bb66-04dc-41d9-aa83-401d11fbcc2e",
		},
		{
			name: cy,
			code: "b7db843c-4954-41da-be28-7324547ff482",
		},
		{
			name: cz,
			code: "e815bc13-206d-4044-beba-4c2b15b61ae3",
		},
		{
			name: ow,
			code: "fde132ed-5ca4-412d-a6ed-afe409be9c65",
		},
		{
			name: sc,
			code: "75a785b4-ac17-47f9-a20e-bcb7dde1850a",
		},
		{
			name: to,
			code: "e047d0d9-114b-4c8a-bb6f-467d69ce2af6",
		},
	}

	for _, c := range districtsToCode {
		if c.code == code {
			dist = c.name
		}
	}
	return dist
}

// CaseByDistrict are the cases reported for a district
type CaseByDistrict struct {
	District string `json:"district"`
}

// CaseCount represents how many cases were reported on a date
type CaseCount struct {
	ReportingDate *time.Time `bson:"_id" json:"reportingDate"`
	Count         int        `bson:"count" json:"count"`
}

// GroupCasesByDate retrieves confirmed cases grouped by the reporting date
func (m *Mongo) GroupCasesByDate(ctx context.Context, outbreakID string, reportingDate time.Time) ([]CaseCount, error) {
	var cases []CaseCount
	collection := m.Client.Database(m.Database).Collection(m.personCollection())
	lastDate := reportingDate.Add(time.Hour * 24)

	matchStage := bson.D{
		{"$match", bson.D{ //nolint:govet
			{"outbreakId", outbreakID}, //nolint:govet
			{"classification", "LNG_REFERENCE_DATA_CATEGORY_CASE_CLASSIFICATION_CONFIRMED"}, //nolint:govet
			{"deleted", false}, //nolint:govet
			{"$and", bson.A{ //nolint:govet
				bson.M{"dateOfReporting": bson.M{"$gte": reportingDate}},
				bson.M{"dateOfReporting": bson.M{"$lt": lastDate}},
			}}},
		},
	}

	groupStage := bson.D{
		{"$group", bson.M{ //nolint:govet
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
