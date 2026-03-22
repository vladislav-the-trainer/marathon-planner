package planner

import "time"

// FitnessLevel represents the user's current running ability
type FitnessLevel string

const (
	FitnessBeginner     FitnessLevel = "beginner"
	FitnessIntermediate FitnessLevel = "intermediate"
	FitnessAdvanced     FitnessLevel = "advanced"
)

// Input holds everything we need to generate a plan
type Input struct {
	FitnessLevel        FitnessLevel
	WeeksUntilRace      int
	TargetFinishMinutes int
	TrainingDaysPerWeek int
	RaceDate            time.Time // optional, used for calendar export later
}

// Plan is the full generated training plan
type Plan struct {
	TotalWeeks int           `json:"total_weeks"`
	RaceDate   string        `json:"race_date"` // YYYY-MM-DD format
	Weeks      []PlannedWeek `json:"weeks"`
}

// PlannedWeek is one week of actual scheduled sessions
type PlannedWeek struct {
	WeekNumber int              `json:"week_number"`
	Sessions   []PlannedSession `json:"sessions"`
}

// PlannedSession is one concrete workout with a date and scaled distance
type PlannedSession struct {
	Date        string      `json:"date"`     // YYYY-MM-DD format
	DayName     string      `json:"day_name"` // "Monday", "Tuesday" etc
	Type        SessionType `json:"type"`
	Description string      `json:"description"`
	DistanceKm  float64     `json:"distance_km"`
}

// fitnessScaleFactor adjusts distances based on fitness level
// beginner does less volume, advanced does more
func fitnessScaleFactor(level FitnessLevel) float64 {
	switch level {
	case FitnessBeginner:
		return 0.8 // 80% of base distances
	case FitnessIntermediate:
		return 1.0 // exactly as in the ODT plan
	case FitnessAdvanced:
		return 1.15 // 115% of base distances
	default:
		return 1.0
	}
}

// weekScaleFactor handles plans shorter or longer than 20 weeks
// maps user's week number to the corresponding base plan week
func mapToBaseWeek(userWeek, totalUserWeeks int) int {
	// Last week is always the race week (week 20, index 19)
	if userWeek == totalUserWeeks {
		return 19
	}

	// simple proportional mapping onto first 19 weeks
	// e.g. if user has 10 weeks, week 5 maps to base week ~9
	baseWeek := int(float64(userWeek-1) / float64(totalUserWeeks-1) * 19)
	if baseWeek < 0 {
		return 0
	}
	if baseWeek > 18 {
		return 18
	}
	return baseWeek
}

// dayOfWeekNames maps day indices to names
var dayOfWeekNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// assignedDayIndices returns the day indices (0=Monday, 6=Sunday) for training days
func assignedDayIndices(daysPerWeek int, weekNumber int) []int {
	if daysPerWeek == 3 {
		return []int{1, 3, 6} // Tuesday, Thursday, Sunday
	}
	// alternate pattern based on odd/even week
	if weekNumber%2 == 1 { // odd weeks
		return []int{1, 3, 4, 6} // Tuesday, Thursday, Friday, Sunday
	}
	// even weeks
	return []int{0, 2, 3, 5} // Monday, Wednesday, Thursday, Saturday
}

// sessionOrder maps days to session types for 4-day week
// Tuesday=interval, Thursday=tempo, Saturday=gym, Sunday=long run
var sessionOrder = []SessionType{
	SessionInterval,
	SessionTempo,
	SessionGym,
	SessionLongRun,
}

// Generate creates a full training plan based on user input with calculated dates
func Generate(input Input) Plan {
	scaleFactor := fitnessScaleFactor(input.FitnessLevel)
	today := time.Now()

	// Calculate the Monday of the race week
	raceDayOfWeek := int(input.RaceDate.Weekday())
	if raceDayOfWeek == 0 {
		raceDayOfWeek = 7 // Sunday = 7 for easier calculation
	}
	raceWeekMonday := input.RaceDate.AddDate(0, 0, -(raceDayOfWeek - 1))

	plan := Plan{
		TotalWeeks: input.WeeksUntilRace,
		RaceDate:   input.RaceDate.Format("2006-01-02"),
		Weeks:      make([]PlannedWeek, input.WeeksUntilRace),
	}

	for userWeek := 0; userWeek < input.WeeksUntilRace; userWeek++ {
		baseWeekIdx := mapToBaseWeek(userWeek+1, input.WeeksUntilRace)
		baseWeek := BasePlan[baseWeekIdx]
		trainingDayIndices := assignedDayIndices(input.TrainingDaysPerWeek, userWeek+1)

		// Calculate additional scaling for extended plans (>20 weeks)
		extendedScaleFactor := 1.0
		if input.WeeksUntilRace > 20 {
			weeksBeforeBasePlan := input.WeeksUntilRace - 20
			if userWeek < weeksBeforeBasePlan {
				reductionWeeks := weeksBeforeBasePlan - userWeek
				extendedScaleFactor = 1.0 - (float64(reductionWeeks) * 0.05)
				if extendedScaleFactor < 0.5 {
					extendedScaleFactor = 0.5
				}
			}
		}

		// Calculate the Monday of this training week
		weeksBack := input.WeeksUntilRace - (userWeek + 1)
		thisWeekMonday := raceWeekMonday.AddDate(0, 0, -weeksBack*7)

		plannedWeek := PlannedWeek{
			WeekNumber: userWeek + 1,
			Sessions:   make([]PlannedSession, 0, 7), // All 7 days
		}

		// Create a map of training day indices to session data
		trainingDayMap := make(map[int]int) // dayIndex -> sessionIndex
		for i, dayIdx := range trainingDayIndices {
			if i < len(sessionOrder) {
				trainingDayMap[dayIdx] = i
			}
		}

		// Generate sessions for all 7 days of the week
		for dayIdx := 0; dayIdx < 7; dayIdx++ {
			sessionDate := thisWeekMonday.AddDate(0, 0, dayIdx)
			dateStr := sessionDate.Format("2006-01-02")
			dayName := dayOfWeekNames[dayIdx]

			// Skip dates in the past or today
			if sessionDate.Before(today) || sessionDate.Format("2006-01-02") == today.Format("2006-01-02") {
				continue
			}

			// Check if this is the marathon race day
			if dateStr == plan.RaceDate {
				plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
					Date:        dateStr,
					DayName:     dayName,
					Type:        SessionLongRun,
					Description: "🏁 MARATHON RACE DAY 🏁",
					DistanceKm:  42.2,
				})
				continue
			}

			// Check if this is the day before the marathon (force rest day)
			dayBeforeRace := input.RaceDate.AddDate(0, 0, -1)
			if dateStr == dayBeforeRace.Format("2006-01-02") {
				plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
					Date:        dateStr,
					DayName:     dayName,
					Type:        SessionRest,
					Description: "Rest day before marathon",
					DistanceKm:  0,
				})
				continue
			}

			// Check if this is after the marathon (post-race empty tile)
			if sessionDate.After(input.RaceDate) {
				plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
					Date:        dateStr,
					DayName:     dayName,
					Type:        "post-race",
					Description: "",
					DistanceKm:  0,
				})
				continue
			}

			// Check if this is a training day
			if sessionIdx, isTrainingDay := trainingDayMap[dayIdx]; isTrainingDay {
				baseSession := baseWeek.Sessions[sessionIdx]
				scaledDistance := baseSession.DistanceKm * scaleFactor * extendedScaleFactor

				if baseSession.Type == SessionGym {
					scaledDistance = 0
				}

				plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
					Date:        dateStr,
					DayName:     dayName,
					Type:        baseSession.Type,
					Description: baseSession.Description,
					DistanceKm:  roundToHalf(scaledDistance),
				})
			} else {
				// Rest day
				plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
					Date:        dateStr,
					DayName:     dayName,
					Type:        SessionRest,
					Description: "Rest day",
					DistanceKm:  0,
				})
			}
		}

		plan.Weeks[userWeek] = plannedWeek
	}

	return plan
}

// roundToHalf rounds distance to nearest 0.5km — more natural for runners
func roundToHalf(km float64) float64 {
	return float64(int(km*2+0.5)) / 2
}
