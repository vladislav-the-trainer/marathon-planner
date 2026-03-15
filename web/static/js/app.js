function plannerApp() {
    return {
        input: {
            fitness_level: 'Beginner',
            marathon_date: '',
            target_finish_min: 240,
            training_days_per_week: 4
        },
        plan: null,
        loading: false,
        error: '',

        // API endpoint - TODO: change to production URL before deployment
        apiUrl: 'http://localhost:8080/api/plan',

        init() {
            // Initialize flatpickr calendar
            const today = new Date();
            const minDate = new Date(today);
            minDate.setDate(minDate.getDate() + 1);

            const defaultDate = new Date(today);
            defaultDate.setDate(defaultDate.getDate() + (20 * 7) + 1); // 20 weeks + 1 day

            this.input.marathon_date = defaultDate.toISOString().split('T')[0];

            flatpickr('#marathon_date', {
                minDate: minDate,
                defaultDate: defaultDate,
                dateFormat: 'Y-m-d',
                onChange: (selectedDates, dateStr) => {
                    this.input.marathon_date = dateStr;
                }
            });
        },

        get minDate() {
            const today = new Date();
            today.setDate(today.getDate() + 1);
            return today.toISOString().split('T')[0];
        },

        get weeksUntilRace() {
            if (!this.input.marathon_date) return 0;
            const today = new Date();
            const raceDate = new Date(this.input.marathon_date);
            const diffTime = raceDate - today;
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            return Math.floor(diffDays / 7);
        },

        formatTime(minutes) {
            if (!minutes) return '';
            const hours = Math.floor(minutes / 60);
            const mins = minutes % 60;
            return `${hours}h ${mins}m`;
        },

        formatMarathonDate(dateStr) {
            if (!dateStr) return '';
            const date = new Date(dateStr);
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            const day = date.getDate();
            const month = months[date.getMonth()];
            const year = date.getFullYear();
            return `${day} ${month} ${year}`;
        },

        formatSessionType(type) {
            const types = {
                'interval': 'Interval',
                'tempo': 'Tempo',
                'long_run': 'Long Run',
                'gym': 'Gym',
                'rest': 'Rest'
            };
            return types[type] || type;
        },

        weekTotal(week) {
            if (!week.Sessions) return 0;
            return week.Sessions.reduce((sum, session) => {
                return sum + (session.DistanceKm || 0);
            }, 0).toFixed(1);
        },

        getFullWeek(week) {
            const allDays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
            const sessionMap = {};
            const raceDate = new Date(this.input.marathon_date + 'T12:00:00');
            const today = new Date();
            today.setHours(0, 0, 0, 0);

            // Map existing sessions by day
            if (week.Sessions) {
                week.Sessions.forEach(session => {
                    sessionMap[session.Day] = session;
                });
            }

            // Create full week with rest days
            const fullWeek = allDays.map((day, index) => {
                const sessionDate = this.getSessionDateObject(week.WeekNumber, index);

                // Normalize dates to YYYY-MM-DD format for comparison
                const sessionDateStr = sessionDate.getFullYear() + '-' +
                    String(sessionDate.getMonth() + 1).padStart(2, '0') + '-' +
                    String(sessionDate.getDate()).padStart(2, '0');
                const raceDateStr = raceDate.getFullYear() + '-' +
                    String(raceDate.getMonth() + 1).padStart(2, '0') + '-' +
                    String(raceDate.getDate()).padStart(2, '0');

                // If this is the race day, replace with marathon event
                if (sessionDateStr === raceDateStr) {
                    return {
                        Day: day,
                        Type: 'long_run',
                        Description: '🏁 MARATHON RACE DAY 🏁',
                        DistanceKm: 42.2,
                        _dayIndex: index
                    };
                }

                // Ensure at least 1 rest day before marathon
                const dayBeforeRace = new Date(raceDate);
                dayBeforeRace.setDate(raceDate.getDate() - 1);
                const dayBeforeRaceStr = dayBeforeRace.getFullYear() + '-' +
                    String(dayBeforeRace.getMonth() + 1).padStart(2, '0') + '-' +
                    String(dayBeforeRace.getDate()).padStart(2, '0');

                if (sessionDateStr === dayBeforeRaceStr) {
                    return {
                        Day: day,
                        Type: 'rest',
                        Description: 'Rest day before marathon',
                        DistanceKm: 0,
                        _dayIndex: index
                    };
                }

                if (sessionMap[day]) {
                    return {
                        ...sessionMap[day],
                        _dayIndex: index
                    };
                } else {
                    return {
                        Day: day,
                        Type: 'rest',
                        Description: 'Rest day',
                        DistanceKm: 0,
                        _dayIndex: index
                    };
                }
            });

            // Filter out sessions in the past and today, but mark sessions after marathon as post-race
            return fullWeek.map((session) => {
                const sessionDate = this.getSessionDateObject(week.WeekNumber, session._dayIndex);
                const sessionDateStr = sessionDate.getFullYear() + '-' +
                    String(sessionDate.getMonth() + 1).padStart(2, '0') + '-' +
                    String(sessionDate.getDate()).padStart(2, '0');
                const raceDateStr = raceDate.getFullYear() + '-' +
                    String(raceDate.getMonth() + 1).padStart(2, '0') + '-' +
                    String(raceDate.getDate()).padStart(2, '0');

                // Mark sessions after the marathon as post-race (empty tiles)
                if (sessionDateStr > raceDateStr) {
                    return {
                        ...session,
                        Type: 'post-race',
                        Description: '',
                        DistanceKm: 0
                    };
                }

                return session;
            }).filter((session) => {
                const sessionDate = this.getSessionDateObject(week.WeekNumber, session._dayIndex);
                return sessionDate > today;
            });
        },

        getSessionDateObject(weekNumber, dayIndex) {
            if (!this.input.marathon_date) return new Date();

            const raceDate = new Date(this.input.marathon_date + 'T00:00:00');
            const totalWeeks = this.plan?.TotalWeeks || 20;

            // Get the day of week for race date (0=Sunday, 1=Monday, etc.)
            const raceDayOfWeek = raceDate.getDay();

            // Calculate the Monday of the race week
            const raceWeekMonday = new Date(raceDate);
            const daysFromMonday = raceDayOfWeek === 0 ? 6 : raceDayOfWeek - 1; // Sunday is 6 days from Monday
            raceWeekMonday.setDate(raceDate.getDate() - daysFromMonday);

            // Calculate weeks back from race week
            const weeksBack = totalWeeks - weekNumber;

            // Calculate the Monday of this week
            const thisWeekMonday = new Date(raceWeekMonday);
            thisWeekMonday.setDate(raceWeekMonday.getDate() - (weeksBack * 7));

            // Add the day index to get the actual session date
            const sessionDate = new Date(thisWeekMonday);
            sessionDate.setDate(thisWeekMonday.getDate() + dayIndex);
            sessionDate.setHours(0, 0, 0, 0);

            return sessionDate;
        },

        getSessionDate(weekNumber, dayIndex) {
            const sessionDate = this.getSessionDateObject(weekNumber, dayIndex);

            // Format as "7 Apr"
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            const day = sessionDate.getDate();
            const month = months[sessionDate.getMonth()];
            return `${day} ${month}`;
        },

        async generatePlan() {
            this.loading = true;
            this.error = '';

            try {
                const requestData = {
                    fitness_level: this.input.fitness_level,
                    weeks_until_race: this.weeksUntilRace,
                    target_finish_min: this.input.target_finish_min,
                    training_days_per_week: this.input.training_days_per_week
                };

                const response = await fetch(this.apiUrl, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(requestData)
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText || 'Failed to generate plan');
                }

                this.plan = await response.json();
            } catch (err) {
                this.error = err.message;
                console.error('Error generating plan:', err);
            } finally {
                this.loading = false;
            }
        },

        reset() {
            this.plan = null;
            this.error = '';
        }
    };
}
