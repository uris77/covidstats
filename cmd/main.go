package main

import (
	"context"
	"covidstats"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")

	server, err := covidstats.NewServer(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	server.RegisterHandlers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
