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
	TotalWeeks int
	Weeks      []PlannedWeek
}

// PlannedWeek is one week of actual scheduled sessions
type PlannedWeek struct {
	WeekNumber int
	Sessions   []PlannedSession
}

// PlannedSession is one concrete workout with a date and scaled distance
type PlannedSession struct {
	Day         string // "Tuesday", "Thursday" etc
	Type        SessionType
	Description string
	DistanceKm  float64
}

// fitnessScaleFactor adjusts distances based on fitness level
// beginner does less volume, advanced does more
func fitnessScaleFactor(level FitnessLevel) float64 {
	switch level {
	case FitnessBeginner:
		return 0.8 // 80% of base distances
	case FitnessAdvanced:
		return 1.15 // 115% of base distances
	default: // intermediate
		return 1.0 // exactly as in the ODT plan
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
	baseWeek := int(float64(userWeek) / float64(totalUserWeeks-1) * 19)
	if baseWeek < 0 {
		return 0
	}
	if baseWeek > 18 {
		return 18
	}
	return baseWeek
}

// assignedDays returns the training days for a given number of days per week
// sensible defaults for marathon training
func assignedDays(daysPerWeek int, weekNumber int) []string {
	if daysPerWeek == 3 {
		return []string{"Tuesday", "Thursday", "Sunday"}
	}
	// alternate pattern based on odd/even week
	if weekNumber%2 == 1 { // odd weeks
		return []string{"Tuesday", "Thursday", "Friday", "Sunday"}
	}
	// even weeks
	return []string{"Monday", "Wednesday", "Thursday", "Saturday"}
}

// sessionOrder maps days to session types for 4-day week
// Tuesday=interval, Thursday=tempo, Saturday=gym, Sunday=long run
var sessionOrder = []SessionType{
	SessionInterval,
	SessionTempo,
	SessionGym,
	SessionLongRun,
}

// Generate creates a full training plan based on user input
func Generate(input Input) Plan {
	scaleFactor := fitnessScaleFactor(input.FitnessLevel)

	plan := Plan{
		TotalWeeks: input.WeeksUntilRace,
		Weeks:      make([]PlannedWeek, input.WeeksUntilRace),
	}

	for userWeek := 0; userWeek < input.WeeksUntilRace; userWeek++ {
		baseWeekIdx := mapToBaseWeek(userWeek+1, input.WeeksUntilRace)
		baseWeek := BasePlan[baseWeekIdx]
		days := assignedDays(input.TrainingDaysPerWeek, userWeek+1)

		// Calculate additional scaling for extended plans (>20 weeks)
		extendedScaleFactor := 1.0
		if input.WeeksUntilRace > 20 {
			weeksBeforeBasePlan := input.WeeksUntilRace - 20
			if userWeek < weeksBeforeBasePlan {
				// Each week before the base plan starts reduces distance by 5%
				reductionWeeks := weeksBeforeBasePlan - userWeek
				extendedScaleFactor = 1.0 - (float64(reductionWeeks) * 0.05)
				// Don't reduce below 50% of base
				if extendedScaleFactor < 0.5 {
					extendedScaleFactor = 0.5
				}
			}
		}

		plannedWeek := PlannedWeek{
			WeekNumber: userWeek + 1,
			Sessions:   make([]PlannedSession, 0, len(days)),
		}

		for i, day := range days {
			if i >= len(sessionOrder) {
				break
			}
			baseSession := baseWeek.Sessions[i]
			scaledDistance := baseSession.DistanceKm * scaleFactor * extendedScaleFactor

			// gym sessions have no distance — don't scale them
			if baseSession.Type == SessionGym {
				scaledDistance = 0
			}

			plannedWeek.Sessions = append(plannedWeek.Sessions, PlannedSession{
				Day:         day,
				Type:        baseSession.Type,
				Description: baseSession.Description,
				DistanceKm:  roundToHalf(scaledDistance),
			})
		}

		plan.Weeks[userWeek] = plannedWeek
	}

	return plan
}

// roundToHalf rounds distance to nearest 0.5km — more natural for runners
func roundToHalf(km float64) float64 {
	return float64(int(km*2+0.5)) / 2
}
