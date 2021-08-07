package stores

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMongo_FindConfirmedCases(t *testing.T) {
	database := os.Getenv("MONGO_DB")
	uri := os.Getenv("MONGO_URI")
	outbreakID := os.Getenv("OUTBREAK_ID")
	isoLayout := "2006-01-02"

	store, err := NewMongoStore(uri, database)
	if err != nil {
		t.Fatalf("mongo connection failed %v", err)
	}
	ctx := context.Background()

	// Connect
	if err := store.Connect(ctx); err != nil {
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

	//t.Logf("cases: %v", cases)
}

func TestMongo_GroupCasesByDate(t *testing.T) {
	database := os.Getenv("MONGO_DB")
	uri := os.Getenv("MONGO_URI")
	outbreakID := os.Getenv("OUTBREAK_ID")
	isoLayout := "2006-01-02"

	store, err := NewMongoStore(uri, database)
	if err != nil {
		t.Fatalf("mongo connection failed %v", err)
	}
	ctx := context.Background()

	// Connect
	if err := store.Connect(ctx); err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	defer store.Disconnect(ctx) //nolint:errcheck

	reportingDate, _ := time.Parse(isoLayout, "2021-08-06")
	cases, err := store.GroupCasesByDate(ctx, outbreakID, reportingDate)
	if err != nil {
		t.Fatalf("grouping cases failed: %v", err)
	}
	t.Logf("cases: %v", cases)
}
