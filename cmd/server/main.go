package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vladislav-the-trainer/marathon-planner/internal/api"
)

var version = "dev" // Set via -ldflags during build

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	r.Get("/health", api.HealthCheck)
	r.Post("/api/plan", api.GeneratePlan)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	log.Printf("Starting Marathon Planner v%s on :8080", version)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
