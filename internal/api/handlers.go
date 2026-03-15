package api

import (
	"encoding/json"
	"net/http"

	"github.com/vladislav-the-trainer/marathon-planner/internal/planner"
)

type QuestionnaireRequest struct {
	FitnessLevel        string `json:"fitness_level"`
	WeeksUntilRace      int    `json:"weeks_until_race"`
	TargetFinishMin     int    `json:"target_finish_min"`
	TrainingDaysPerWeek int    `json:"training_days_per_week"`
}

func GeneratePlan(w http.ResponseWriter, r *http.Request) {
	var req QuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate weeks - minimum 4 weeks, no hard maximum (extended plans supported)
	if req.WeeksUntilRace < 4 {
		http.Error(w, "weeks_until_race must be at least 4 weeks", http.StatusBadRequest)
		return
	}

	input := planner.Input{
		FitnessLevel:        planner.FitnessLevel(req.FitnessLevel),
		WeeksUntilRace:      req.WeeksUntilRace,
		TargetFinishMinutes: req.TargetFinishMin,
		TrainingDaysPerWeek: req.TrainingDaysPerWeek,
	}

	plan := planner.Generate(input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
