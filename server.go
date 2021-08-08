package covidstats

import (
	"context"
	"covidstats/stores"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server encapsulates the covidstats backend
type Server struct {
	GCPProjectID    string
	FirestoreClient *stores.Firestore
	casesService    *stores.CasesByDateService
	router          *mux.Router
	logger          *logrus.Logger
}

// NewServer instantiates new server
func NewServer(ctx context.Context, logger *logrus.Logger, gcpProjectID string) (*Server, error) {
	firestoreClient, err := stores.CreateFirestoreDB(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate firestore client: %w", err)
	}
	svc := stores.NewCasesByDateService(firestoreClient, "covid_cases_stats")
	return &Server{
		GCPProjectID:    gcpProjectID,
		FirestoreClient: firestoreClient,
		casesService:    svc,
		router:          mux.NewRouter().PathPrefix("/api").Subrouter(),
		logger:          logger,
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

// registerHandlers registers handlers
func (s *Server) registerHandlers() {
	h := NewChain(enableCors())
	s.router.HandleFunc("/byYear/{year:[0-9]+}", h.Then(s.HandleFindYearStats)).
		Methods(http.MethodOptions, http.MethodGet)
}

// Start boots up the server
func (s *Server) Start(port string) {
	s.registerHandlers()
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}

	// Run server in a goroutine so that it does not block.
	go func() {
		s.logger.Infof("Starting server on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	wait := time.Duration(30)
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	//Does not block if no connections, but will otherwise wait
	// until the timeout deadline
	_ = srv.Shutdown(ctx)
	s.logger.Info("shutting down")
	os.Exit(0)
}
