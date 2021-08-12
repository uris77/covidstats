package stores

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const isoLayout string = "2006-01-02"

func TestMongo_FindConfirmedCases(t *testing.T) {
	database := os.Getenv("MONGO_DB")
	uri := os.Getenv("MONGO_URI")
	outbreakID := os.Getenv("OUTBREAK_ID")

	store, err := NewMongoStore(uri, database)
	if err != nil {
		t.Fatalf("mongo connection failed %v", err)
	}
	ctx := context.Background()

	// Connect
	if err := store.Connect(ctx); err != nil { //nolint:govet
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	defer store.Disconnect(ctx) //nolint:errcheck

	reportingDate, _ := time.Parse(isoLayout, "2021-06-01")
	lastDate, _ := time.Parse(isoLayout, "2021-09-01")
	cases, err := store.FindConfirmedCases(ctx, outbreakID, reportingDate, &lastDate)
	if err != nil {
		t.Fatalf("finding cases failed: %v", err)
	}

	file, _ := json.MarshalIndent(cases, "", "  ")
	_ = ioutil.WriteFile("cases.json", file, 0600)
}

func TestMongo_GroupCasesByDate(t *testing.T) {
	database := os.Getenv("MONGO_DB")
	uri := os.Getenv("MONGO_URI")
	outbreakID := os.Getenv("OUTBREAK_ID")

	store, err := NewMongoStore(uri, database)
	if err != nil {
		t.Fatalf("mongo connection failed %v", err)
	}
	ctx := context.Background()

	// Connect
	if err := store.Connect(ctx); err != nil { //nolint:govet
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	defer store.Disconnect(ctx) //nolint:errcheck

	reportingDate, _ := time.Parse(isoLayout, "2021-08-09")
	cases, err := store.GroupCasesByDate(ctx, outbreakID, reportingDate)
	if err != nil {
		t.Fatalf("grouping cases failed: %v", err)
	}
	t.Logf("cases: %v", cases)
}

func TestMongo_Upload(t *testing.T) {
	database := os.Getenv("MONGO_DB")
	uri := os.Getenv("MONGO_URI")
	outbreakID := os.Getenv("OUTBREAK_ID")

	store, err := NewMongoStore(uri, database)
	if err != nil {
		t.Fatalf("mongo connection failed %v", err)
	}
	ctx := context.Background()

	// Connect
	if err := store.Connect(ctx); err != nil { //nolint:govet
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	defer store.Disconnect(ctx) //nolint:errcheck

	reportingDate, _ := time.Parse(isoLayout, "2021-08-11")
	cases, err := store.GroupCasesByDate(ctx, outbreakID, reportingDate)
	if err != nil {
		t.Fatalf("grouping cases failed: %v", err)
	}

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
