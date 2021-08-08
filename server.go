package covidstats

import (
	"context"
	"covidstats/stores"
	"fmt"
	"net/http"
)

// Server encapsulates the covidstats backend
type Server struct {
	GCPProjectID    string
	FirestoreClient *stores.Firestore
	casesService    *stores.CasesByDateService
}

// NewServer instantiates new server
func NewServer(ctx context.Context, gcpProjectID string) (*Server, error) {
	firestoreClient, err := stores.CreateFirestoreDB(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate firestore client: %w", err)
	}
	svc := stores.NewCasesByDateService(firestoreClient, "covid_cases_stats")
	return &Server{
		GCPProjectID:    gcpProjectID,
		FirestoreClient: firestoreClient,
		casesService:    svc,
	}, nil
}

func enableCors() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Referer, Connection, X-POE-Authorization")
			w.Header().Set("responseType", "*")
			f(w, r)
		}
	}
}

// RegisterHandlers registers handlers
func (s *Server) RegisterHandlers() {
	h := Chainz(s.HandleFindYearStats, enableCors())
	http.HandleFunc("/byYear/{year:[0-9]+}", h)
}
