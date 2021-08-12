package stores

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"
)

func TestCasesByDateService_Save(t *testing.T) {
	// Read the contents from a file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		t.Fatalf("failed to read data.json: %v", err)
	}

	type rawDataID struct {
		Date string `json:"$date"`
	}

	type rawData struct {
		ID    rawDataID `json:"_id"`
		Count int       `json:"count"`
	}

	var raw []rawData
	var cases []CaseCount

	if err := json.Unmarshal(data, &raw); err != nil { //nolint:govet
		t.Fatalf("failed to unmarshall cases: %v", err)
	}

	for _, r := range raw {
		d, err := time.Parse("2006-01-02", r.ID.Date[0:10]) //nolint:govet
		if err != nil {
			t.Fatalf("parsing date failed: %v", err)
		}
		cases = append(cases, CaseCount{
			ReportingDate: &d,
			Count:         r.Count,
		})
	}

	ctx := context.Background()
	// create Firestore Client
	fsClient, err := CreateFirestoreDB(ctx, "epi-belize")
	if err != nil {
		t.Fatalf("failed to create firestore db: %v", err)
	}

	// create Service
	svc := NewCasesByDateService(fsClient, "covid_cases_stats")

	// Save data
	if err := svc.Save(ctx, cases); err != nil {
		t.Fatalf("saving cases failed: %v", err)
	}
}

func TestCasesByDateService_FindByMonth(t *testing.T) {
	ctx := context.Background()
	// create Firestore Client
	fsClient, err := CreateFirestoreDB(ctx, "epi-belize")
	if err != nil {
		t.Fatalf("failed to create firestore db: %v", err)
	}
	svc := NewCasesByDateService(fsClient, "covid_cases_stats")
	cases, err := svc.FindByMonth(ctx, "2021-08")
	if err != nil {
		t.Fatalf("FindByMonth failed: %v", err)
	}
	t.Logf("cases: %v", len(cases))
}

func TestCasesByDateService_FindByYear(t *testing.T) {
	ctx := context.Background()
	// create Firestore Client
	fsClient, err := CreateFirestoreDB(ctx, "epi-belize")
	if err != nil {
		t.Fatalf("failed to create firestore db: %v", err)
	}

	// create Service
	svc := NewCasesByDateService(fsClient, "covid_cases_stats")

	cases, err := svc.FindByYear(ctx, 2021)
	if err != nil {
		t.Fatalf("FindByYear failed: %v", err)
	}
	t.Logf("cases: %v", cases)
}
