package planner

// SessionType categorizes each workout
type SessionType string

const (
	SessionInterval SessionType = "interval"
	SessionTempo    SessionType = "tempo"
	SessionLongRun  SessionType = "long_run"
	SessionGym      SessionType = "gym"
	SessionRest     SessionType = "rest"
)

// BaseSession represents one workout in the canonical 20-week plan
type BaseSession struct {
	Type        SessionType
	Description string
	DistanceKm  float64 // base distance for intermediate runner
}

// BaseWeek represents one week in the canonical plan
type BaseWeek struct {
	Sessions [4]BaseSession // always 4 sessions: interval, tempo, gym, long run
}

// BasePlan is the canonical 20-week plan derived from the ODT file
// Sessions order: [0]=Interval, [1]=Tempo, [2]=Gym, [3]=Long Run
// Distance calculation rule: rest is counted BETWEEN intervals only, not after the last one
// Formula: N×interval + (N-1)×rest for distance-based rest
// Time-based rest (min/sec): only interval distances counted, no rest distance added
var BasePlan = [20]BaseWeek{
	// Week 1
	{Sessions: [4]BaseSession{
		{SessionInterval, "5 x 400m, rest 20 sec", 2.0}, // time-based rest: 5×400m = 2.0km
		{SessionTempo, "Easy 1.5km + 3km at 10km pace + Easy 1.5km", 6.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +23 sec/km", 12.0},
	}},
	// Week 2
	{Sessions: [4]BaseSession{
		{SessionInterval, "4 x 800m, rest 400m jog", 4.4}, // 4×800m + 3×400m = 3.2km + 1.2km = 4.4km
		{SessionTempo, "Easy 1.5km + 5km at 10km pace + Easy 1.5km", 8.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 14.0},
	}},
	// Week 3
	{Sessions: [4]BaseSession{
		{SessionInterval, "4 x 600m, rest 200m jog", 3.0}, // 4×600m + 3×200m = 2.4km + 0.6km = 3.0km
		{SessionTempo, "Easy 1.5km + 3km at 5km pace + Easy 3km", 7.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 17.0},
	}},
	// Week 4
	{Sessions: [4]BaseSession{
		{SessionInterval, "12 x 200m, rest 200m jog", 4.6}, // 12×200m + 11×200m = 2.4km + 2.2km = 4.6km
		{SessionTempo, "Easy 3km + 5km at 15km pace + Easy 3km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 20.0},
	}},
	// Week 5
	{Sessions: [4]BaseSession{
		{SessionInterval, "3 x 1600m, rest 400m jog", 5.6}, // 3×1600m + 2×400m = 4.8km + 0.8km = 5.6km
		{SessionTempo, "Easy 3km + 3km at 5km pace + Easy 3km", 9.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 21.0},
	}},
	// Week 6
	{Sessions: [4]BaseSession{
		{SessionInterval, "4 x 800m, rest 2 min jog", 3.2}, // time-based rest: 4×800m = 3.2km
		{SessionTempo, "Easy 1.5km + 8km at marathon pace + Easy 1.5km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +28 sec/km", 24.0},
	}},
	// Week 7
	{Sessions: [4]BaseSession{
		{SessionInterval, "Ladder: 1200m+1000m+800m+600m+400m, rest 200m jog", 4.8}, // 5 intervals + 4×200m = 4.0km + 0.8km = 4.8km
		{SessionTempo, "Easy 1.5km + 8km at 10km pace + Easy 1.5km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +28 sec/km", 27.0},
	}},
	// Week 8
	{Sessions: [4]BaseSession{
		{SessionInterval, "5 x 1000m, rest 400m jog", 6.6}, // 5×1000m + 4×400m = 5.0km + 1.6km = 6.6km
		{SessionTempo, "Easy 1.5km + 6.5km at marathon pace + Easy 1.5km", 9.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +37 sec/km", 32.0},
	}},
	// Week 9
	{Sessions: [4]BaseSession{
		{SessionInterval, "3 x 1600m, rest 400m jog", 5.6}, // 3×1600m + 2×400m = 4.8km + 0.8km = 5.6km
		{SessionTempo, "Easy 3km + 5km at 5km pace + Easy 1.5km", 9.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +28 sec/km", 29.0},
	}},
	// Week 10
	{Sessions: [4]BaseSession{
		{SessionInterval, "2 x 1200m + 4 x 800m, rest 2 min jog", 5.6}, // time-based rest: 2×1200m + 4×800m = 2.4km + 3.2km = 5.6km
		{SessionTempo, "Easy 1.5km + 8km at marathon pace + Easy 1.5km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +28 sec/km", 32.0},
	}},
	// Week 11
	{Sessions: [4]BaseSession{
		{SessionInterval, "6 x 800m, rest 90 sec jog", 4.8}, // time-based rest: 6×800m = 4.8km
		{SessionTempo, "Easy 1.5km + 10km at half marathon pace + Easy 1.5km", 13.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +9 sec/km", 21.0},
	}},
	// Week 12
	{Sessions: [4]BaseSession{
		{SessionInterval, "2 sets of 6 x 400m, rest 90 sec inside set, 2:30 between sets", 4.8}, // time-based rest: 12×400m = 4.8km
		{SessionTempo, "Easy 3km + 5km at 5km pace + Easy 1.5km", 9.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 29.0},
	}},
	// Week 13
	{Sessions: [4]BaseSession{
		{SessionInterval, "2 x 1600m + 2 x 800m, rest 60 sec jog", 4.8}, // time-based rest: 2×1600m + 2×800m = 3.2km + 1.6km = 4.8km
		{SessionTempo, "Easy 1.5km + 6.5km at half marathon pace + Easy 1.5km", 9.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 32.0},
	}},
	// Week 14
	{Sessions: [4]BaseSession{
		{SessionInterval, "4 x 1200m, rest 2 min jog", 4.8}, // time-based rest: 4×1200m = 4.8km
		{SessionTempo, "16km at marathon pace", 16.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +12 sec/km", 24.0},
	}},
	// Week 15
	{Sessions: [4]BaseSession{
		{SessionInterval, "1000m+2000m+1000m+2000m, rest 400m jog", 7.2}, // 4 intervals + 3×400m = 6.0km + 1.2km = 7.2km
		{SessionTempo, "Easy 1.5km + 8km at half marathon pace + Easy 1.5km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +19 sec/km", 32.0},
	}},
	// Week 16
	{Sessions: [4]BaseSession{
		{SessionInterval, "3 x 1600m, rest 400m jog", 5.6}, // 3×1600m + 2×400m = 4.8km + 0.8km = 5.6km
		{SessionTempo, "16km at marathon pace", 16.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +9 sec/km", 24.0},
	}},
	// Week 17
	{Sessions: [4]BaseSession{
		{SessionInterval, "10 x 400m, rest 400m jog", 7.6}, // 10×400m + 9×400m = 4.0km + 3.6km = 7.6km
		{SessionTempo, "Warmup 10min + 13km at marathon pace + Cooldown 10min", 13.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace +9 sec/km", 32.0},
	}},
	// Week 18
	{Sessions: [4]BaseSession{
		{SessionInterval, "8 x 800m, rest 90 sec jog", 6.4}, // time-based rest: 8×800m = 6.4km
		{SessionTempo, "Easy 1.5km + 8km at marathon pace + Easy 1.5km", 11.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace", 21.0},
	}},
	// Week 19
	{Sessions: [4]BaseSession{
		{SessionInterval, "5 x 1000m, rest 400m jog", 6.6}, // 5×1000m + 4×400m = 5.0km + 1.6km = 6.6km
		{SessionTempo, "Easy 3km + 5km at 5km pace + Easy 1.5km", 9.5},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "Long run at marathon pace", 16.0},
	}},
	// Week 20 — race week
	{Sessions: [4]BaseSession{
		{SessionInterval, "6 x 400m, rest 400m jog", 4.4}, // 6×400m + 5×400m = 2.4km + 2.0km = 4.4km
		{SessionTempo, "Warmup 10min + 5km at marathon pace + Cooldown 10min", 5.0},
		{SessionGym, "Full body strength training", 0},
		{SessionLongRun, "MARATHON RACE DAY!", 42.2},
	}},
}
