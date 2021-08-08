package main

import (
	"context"
	"covidstats"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	logger.Infof("GOOGLE_APPLICATION_CREDENTIALS %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	server, err := covidstats.NewServer(ctx, logger, projectID)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	server.Start(port)
}
