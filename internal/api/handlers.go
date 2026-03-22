package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vladislav-the-trainer/marathon-planner/internal/planner"
)

type QuestionnaireRequest struct {
	FitnessLevel        string `json:"fitness_level"`
	RaceDate            string `json:"race_date"` // ISO format: YYYY-MM-DD
	TargetFinishMin     int    `json:"target_finish_min"`
	TrainingDaysPerWeek int    `json:"training_days_per_week"`
}

func GeneratePlan(w http.ResponseWriter, r *http.Request) {
	var req QuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse race date
	raceDate, err := time.Parse("2006-01-02", req.RaceDate)
	if err != nil {
		http.Error(w, "Invalid race_date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Validate race date is in the future
	if raceDate.Before(time.Now()) {
		http.Error(w, "race_date must be in the future", http.StatusBadRequest)
		return
	}

	// Calculate weeks until race
	weeksUntilRace := int(time.Until(raceDate).Hours() / 24 / 7)
	if weeksUntilRace < 4 {
		http.Error(w, "race_date must be at least 4 weeks in the future", http.StatusBadRequest)
		return
	}

	input := planner.Input{
		FitnessLevel:        planner.FitnessLevel(req.FitnessLevel),
		WeeksUntilRace:      weeksUntilRace,
		RaceDate:            raceDate,
		TargetFinishMinutes: req.TargetFinishMin,
		TrainingDaysPerWeek: req.TrainingDaysPerWeek,
	}

	plan := planner.Generate(input)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(plan); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
