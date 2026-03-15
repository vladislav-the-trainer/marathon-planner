# Marathon Training Planner

Web app that generates personalized marathon training plans based on your fitness level, available training days, and target finish time. Built with Go and Alpine.js.

Built as a personal tool for a triathlete who actually uses it.

## Features

- **Personalized Training Plans**: Generate custom 20-week marathon training plans
- **Adaptive Difficulty**: Plans scale for beginner, intermediate, and advanced runners
- **Calendar Integration**: Interactive date picker with automatic session date calculation
- **Smart Scheduling**: Alternating training day patterns optimized for recovery
- **Race Day Ready**: Automatic rest day before marathon, special race day highlighting
- **Progress Tracking**: Color-coded session types, weekly distance totals
- **Responsive Design**: Works on desktop and mobile
- **Filter Past Sessions**: Only shows upcoming training sessions

## Tech Stack

**Backend:**
- Go 1.26.1
- Chi router v5
- SQLite (modernc.org/sqlite - pure Go, no CGo)
- CORS enabled

**Frontend:**
- Alpine.js v3 (CDN, no build step)
- Flatpickr date picker
- Vanilla HTML/CSS/JS
- S3-ready static files

## Getting Started

### Prerequisites

- Go 1.26.1 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/vladislav-the-trainer/marathon-planner.git
cd marathon-planner
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

4. Open your browser:
```
http://localhost:8080
```

## Usage

1. **Select your fitness level**: Beginner, Intermediate, or Advanced
2. **Choose your marathon date**: Use the calendar picker (defaults to 20 weeks from today)
3. **Set your target finish time**: Enter your goal in minutes
4. **Select training days**: Choose 3 or 4 days per week
5. **Generate your plan**: Get a personalized 20-week training schedule

### Training Plan Features

- **Session Types**:
  - **Interval** (Orange): High-intensity speed work
  - **Tempo** (Purple): Sustained threshold running
  - **Long Run** (Green): Weekly distance builder
  - **Gym** (Indigo): Strength and conditioning
  - **Rest** (Gray): Recovery days

- **Smart Scheduling**:
  - Alternating weekly patterns for optimal recovery
  - Automatic rest day before marathon
  - Post-race dates shown as empty tiles
  - Past sessions automatically hidden

## API Endpoints

### GET /health
Health check endpoint
```bash
curl http://localhost:8080/health
```

### POST /api/plan
Generate training plan
```bash
curl -X POST http://localhost:8080/api/plan \
  -H "Content-Type: application/json" \
  -d '{
    "fitness_level": "beginner",
    "weeks_until_race": 20,
    "target_finish_min": 240,
    "training_days_per_week": 4
  }'
```

## Project Structure

```
marathon-planner/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, HTTP server
├── internal/
│   ├── api/
│   │   └── handlers.go          # HTTP handlers
│   └── planner/
│       ├── baseplan.go          # 20-week base plan data
│       └── planner.go           # Plan generation logic
├── web/
│   ├── index.html               # Alpine.js SPA
│   └── static/
│       ├── css/styles.css       # Styling
│       └── js/app.js            # Frontend logic
├── go.mod
└── README.md
```

## Development

### Building for Production

```bash
go build -o marathon-planner cmd/server/main.go
```

### Running Tests

```bash
go test ./...
```

## Deployment

### Frontend (AWS S3)
- Upload `web/` directory to S3 bucket
- Enable static website hosting
- Configure CloudFront for HTTPS

### Backend (AWS EC2/ECS)
- Build binary: `go build -o marathon-planner cmd/server/main.go`
- Run on port 8080
- Configure CORS for S3 frontend origin

## Training Methodology

The base plan follows a structured 20-week progression:
- **Weeks 1-8**: Base building
- **Weeks 9-16**: Peak training with increasing long runs
- **Weeks 17-19**: Taper period
- **Week 20**: Race week

Plans are scaled based on fitness level:
- **Beginner**: 0.8x base distances
- **Intermediate**: 1.0x base distances
- **Advanced**: 1.15x base distances

## Roadmap

- [ ] PDF export functionality
- [ ] ICS calendar export
- [ ] User plan editing in browser
- [ ] SQLite storage for saved plans
- [ ] User authentication
- [ ] Plan sharing functionality

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Built by [Vladislav the Trainer](https://github.com/vladislav-the-trainer)

## Acknowledgments

- Base training plan derived from "Get Prepared for Marathon During 20 Weeks" methodology
- Built for runners, by a runner
