package stores

import (
	fs "cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"time"
)

// Firestore represents the database connection.
type Firestore struct {
	// The firestore client
	Client *fs.Client
	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now       func() time.Time
	projectID string
}

// CreateFirestoreDB creates a new Firestore connection
func CreateFirestoreDB(ctx context.Context, projectID string) (*Firestore, error) {
	c, clientErr := fs.NewClient(ctx, projectID)
	if clientErr != nil {
		return nil, fmt.Errorf("CreateFirestoreDB: could not create new firestore client: %w", clientErr)
	}
	return &Firestore{
		Client:    c,
		Now:       time.Now,
		projectID: projectID,
	}, nil
}

type CasesByDateService struct {
	db         *Firestore
	collection string
	colRef     *fs.CollectionRef
}

func NewCasesByDateService(db *Firestore, collection string) *CasesByDateService {
	return &CasesByDateService{
		db:         db,
		collection: collection,
		colRef:     db.Client.Collection(collection),
	}
}

func (c *CasesByDateService) Save(ctx context.Context, cases []CaseCount) error {
	batch := c.db.Client.Batch()

	for _, cs := range cases {
		ID := cs.ReportingDate.Format("2006-01-02")
		ref := c.colRef.Doc(ID)
		month := cs.ReportingDate.Month()
		monthStr := fmt.Sprintf("%d", month)
		if month < 10 {
			monthStr = fmt.Sprintf("0%d", month)
		}
		batch.Set(ref, map[string]interface{}{
			"reportingDate": cs.ReportingDate,
			"count":         cs.Count,
			"year":          cs.ReportingDate.Year(),
			"month":         fmt.Sprintf("%d-%s", cs.ReportingDate.Year(), monthStr),
		}, fs.MergeAll)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to save cases: %w", err)
	}

	return nil
}

type CasesCountByDate struct {
	ReportingDate *time.Time
	Count         int32
	Year          int32
	Month         string
}

func (c *CasesByDateService) FindByMonth(ctx context.Context, month string) ([]CasesCountByDate, error) {
	var cases []CasesCountByDate
	iter := c.colRef.Query.Where("month", "==", month).Documents(ctx)

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return cases, fmt.Errorf("FindByMonth() error: %w", err)
		}

		var cs CasesCountByDate
		dataErr := doc.DataTo(&cs)
		if dataErr != nil {
			return cases, fmt.Errorf("FindByMonth: unmarshal error: %w", err)
		}
		cases = append(cases, cs)
	}
	return cases, nil
}

func (c *CasesByDateService) FindByYear(ctx context.Context, year int32) ([]CasesCountByDate, error) {
	var cases []CasesCountByDate
	iter := c.colRef.Query.Where("year", "==", year).Documents(ctx)

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return cases, fmt.Errorf("FindByMonth() error: %w", err)
		}

		var cs CasesCountByDate
		dataErr := doc.DataTo(&cs)
		if dataErr != nil {
			return cases, fmt.Errorf("FindByMonth: unmarshal error: %w", err)
		}
		cases = append(cases, cs)
	}
	return cases, nil
}
