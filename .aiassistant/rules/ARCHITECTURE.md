---
apply: always
---

# Marathon Planner — AI Assistant Context

## Project Overview
A Go web app that generates personalized marathon training plans based on
user fitness level, available training days, and target finish time.
The base 20-week plan is derived from a structured training document
(get-prepared-for-marathon-during-20-weeks.odt).

Built as a personal tool for a triathlete who actually uses it.

## Technical Stack
- Language: Go 1.26.1
- Router: Chi v5 (github.com/go-chi/chi/v5)
- CORS: github.com/go-chi/cors
- Database: SQLite via modernc.org/sqlite (pure Go, no CGo)
- Frontend: Alpine.js v3 + vanilla HTML/CSS/JS (no build step, S3-ready)
- Date Picker: Flatpickr (CDN)
- Deployment target: Frontend on S3 + CloudFront, Backend on AWS EC2/ECS

## Project Structure
```
marathon-planner/
├── cmd/
│   └── server/
│       └── main.go              ← entry point, HTTP server, route registration
├── internal/
│   ├── api/
│   │   └── handlers.go          ← HTTP handlers, request/response structs
│   └── planner/
│       ├── baseplan.go          ← canonical 20-week base plan data
│       └── planner.go           ← plan generation logic, scaling, day assignment
├── web/
│   ├── index.html               ← Alpine.js SPA (questionnaire + plan display)
│   └── static/
│       ├── css/
│       │   └── styles.css       ← responsive styles with color-coded sessions
│       └── js/
│           └── app.js           ← Alpine.js app logic, API integration, calendar logic
├── .aiassistant/
│   └── rules/
│       └── ARCHITECTURE.md      ← this file
├── go.mod
├── go.sum
└── README.md
```

## Module Path
github.com/vladislav-the-trainer/marathon-planner

## Current State
### Implemented Features
- ✅ Project scaffolded with correct directory structure
- ✅ Go module initialized
- ✅ Chi router wired in main.go with CORS support
- ✅ Health check endpoint: GET /health
- ✅ Plan generation endpoint: POST /api/plan
- ✅ Base 20-week plan encoded in baseplan.go
- ✅ Plan generation logic in planner.go (scaling by fitness level, week mapping)
- ✅ Day assignment alternating pattern (odd/even weeks) implemented
- ✅ Alpine.js frontend (questionnaire form + plan display)
- ✅ Color-coded session cards (interval/tempo/gym/long run/race day)
- ✅ Static file serving from web/ directory
- ✅ Calendar date picker integration (Flatpickr)
- ✅ Dynamic date calculation for all training sessions
- ✅ Marathon race day highlighting with special styling
- ✅ Automatic rest day insertion before marathon
- ✅ Filter out past/today sessions from display
- ✅ Post-race date handling (show empty tiles after marathon)
- ✅ Default values for form fields (beginner, 4 days/week, 20 weeks + 1 day from today)
- ✅ Compressed plan warning for <20 weeks

### Not Yet Implemented
- ⏳ SQLite storage layer (not needed until "save my plan" feature is built)
- ⏳ PDF export
- ⏳ ICS calendar export
- ⏳ User plan editing in browser
- ⏳ S3 deployment (frontend)
- ⏳ EC2/ECS deployment (backend)

## Key Design Decisions

### Calendar Integration
- **Date Picker**: Flatpickr library for user-friendly calendar UI
- **Default Date**: Current date + 20 weeks + 1 day (standard marathon training period)
- **Date Display Format**: Session cards show "DayName, DD MMM" (e.g., "Monday, 15 Mar")
- **Marathon Date Display**: Full format "DD MMM YYYY" in plan summary

### Date Calculation Logic
- Week alignment: All weeks start on Monday and end on Sunday
- Race date can fall on any day of the week
- Session dates calculated by:
  1. Finding the Monday of the race week
  2. Going back N weeks to find the Monday of current training week
  3. Adding day index (0-6) to get specific session date
- **Important**: Session date calculation uses the marathon date as anchor point

### Marathon Day Handling
- Marathon replaces any scheduled training on race date
- Day before marathon forced to be rest day (regardless of plan)
- Marathon tile styling: Deep blue background (#2424a7), red border, white text
- Marathon description: "🏁 MARATHON RACE DAY 🏁"
- Distance: 42.2 km

### Post-Race Day Handling
- All days after marathon show as empty tiles
- Styled like rest days (gray, low opacity)
- Only display the date, no type/description/distance
- Session type: 'post-race'

### Past Date Filtering
- Sessions in the past (< today) are hidden
- Sessions for today are also hidden
- Only shows sessions from tomorrow onwards
- Preserves correct date calculation using `_dayIndex` property

### Day Assignment Pattern
Training days alternate by week — this is intentional training design:
- Odd weeks:  Tuesday, Thursday, Friday, Sunday
- Even weeks: Monday, Wednesday, Thursday, Saturday

Session to day mapping:
- Day 1 (Tue or Mon): Interval — hardest quality session, done fresh
- Day 2 (Thu or Wed): Tempo — second quality session
- Day 3 (Fri or Thu): Gym — full body strength, acts as active buffer
- Day 4 (Sun or Sat): Long Run — preceded by one gym day + one full rest day

Recovery logic:
- Interval → 1 rest day → Tempo (legs partially recovered)
- Tempo → next day → Gym (strength work, not cardio, legs not additionally fatigued)
- Gym → 1 full rest day → Long Run (legs fresh enough for maximum distance effort)
- The Gym session is a deliberate buffer between Tempo and Long Run
- Pattern: Interval ... Tempo → Gym → Rest → Long Run

Note: This is NOT two full rest days before the long run.
It is one active day (Gym) + one complete rest day.

### Session Order (applies to both odd and even weeks)
Index 0 = Interval (first day of week)
Index 1 = Tempo (second day)
Index 2 = Gym (third day)
Index 3 = Long Run (last day, Sunday or Saturday)

### Fitness Level Scaling
- beginner:     0.8x base distances
- intermediate: 1.0x (exact ODT plan distances)
- advanced:     1.15x base distances
  Gym sessions are never scaled (distance = 0).

### Week Scaling
User may have fewer or more than 20 weeks until race.
mapToBaseWeek() maps user's week proportionally onto the 20-week base plan.
Formula: baseWeek = int(userWeek / totalUserWeeks * 20)

### Distance Rounding
All distances rounded to nearest 0.5km — more natural for runners.

## API Endpoints

### GET /health
Returns 200 OK with body "OK"

### POST /api/plan
Request example:
```json
{
  "fitness_level": "beginner",
  "weeks_until_race": 20,
  "target_finish_min": 240,
  "training_days_per_week": 4
}
```

Valid values:
- `fitness_level`: "beginner", "intermediate", or "advanced"
- `weeks_until_race`: 4-52 (integer)
- `target_finish_min`: target finish time in minutes (integer)
- `training_days_per_week`: 3 or 4

Validation: weeks_until_race must be between 4 and 52.

Response: Full Plan struct serialized as JSON.

## Code Conventions
- Idiomatic Go — no unnecessary abstractions
- internal/ packages are private to this module by Go convention
- Error handling: always check and handle errors explicitly
- Comments on non-obvious logic, especially training domain concepts
- Keep business logic (planner/) completely separate from HTTP (api/)
- Frontend: Use Alpine.js directives for reactivity, minimize DOM manipulation

## Frontend Architecture
- **Framework**: Alpine.js v3 (CDN, no build step)
- **Date Picker**: Flatpickr (CDN)
- **Structure**: Single-page app with reactive state
- **API Integration**: Fetch API calling Go backend at /api/plan
- **Styling**: Custom CSS with CSS variables, responsive grid layout

### Session Colors
- **Interval**: Orange (#f59e0b) with light orange background
- **Tempo**: Purple (#8b5cf6) with light purple background
- **Long Run**: Green (#10b981) with light green background
- **Gym**: Indigo (#6366f1) with light indigo background
- **Rest**: Gray (#f3f4f6) with 70% opacity
- **Race Day**: Deep blue (#2424a7) with red border and white text
- **Post-Race**: Same as rest days (gray, 70% opacity)

### Form Defaults
- **Fitness Level**: Beginner
- **Marathon Date**: Today + 20 weeks + 1 day
- **Target Finish Time**: 240 minutes (4 hours)
- **Training Days Per Week**: 4 days

### User Experience Features
- Warning banner for plans <20 weeks (displayed on plan view page)
- Marathon date displayed prominently in plan summary
- Week total distance calculated and displayed
- Responsive grid layout (7 columns on desktop, 2 on mobile)
- Sessions filtered to show only future dates
- Color-coded session types for quick visual reference

## Known Issues / Technical Debt
None currently identified.

## Upcoming Tasks (in priority order)
1. ⏳ PDF export functionality
2. ⏳ ICS calendar export
3. ⏳ User plan editing in browser
4. ⏳ SQLite storage — only needed when "save my plan" is added
5. ⏳ S3 deployment for frontend
6. ⏳ EC2/ECS deployment for backend
7. ⏳ Environment-based API URL configuration

## Deployment Notes
### Frontend (S3)
- Upload web/ directory to S3 bucket
- Enable static website hosting
- Update apiUrl in app.js to point to backend endpoint
- Configure CloudFront for HTTPS and caching

### Backend (EC2/ECS)
- Build Go binary: `go build -o marathon-planner cmd/server/main.go`
- Run on port 8080 (or configure via environment)
- Ensure CORS configured for S3 frontend origin
- Consider using Docker for ECS deployment
