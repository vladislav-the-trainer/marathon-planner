function plannerApp() {
    return {
        input: {
            fitness_level: 'intermediate',
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
            // Use plan.race_date if available, otherwise use input
            const raceDateStr = this.plan?.race_date || dateStr;
            const date = new Date(raceDateStr);
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
                'jogging': 'Jogging',
                'rest': 'Rest'
            };
            return types[type] || type;
        },

        weekTotal(week) {
            if (!week.sessions) return 0;
            return week.sessions.reduce((sum, session) => {
                return sum + (session.distance_km || 0);
            }, 0).toFixed(1);
        },

        formatSessionDate(dateStr) {
            if (!dateStr) return '';
            const date = new Date(dateStr);
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            const day = date.getDate();
            const month = months[date.getMonth()];
            return `${day} ${month}`;
        },

        async generatePlan() {
            this.loading = true;
            this.error = '';

            try {
                const requestData = {
                    fitness_level: this.input.fitness_level,
                    race_date: this.input.marathon_date,
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
                    this.error = errorText || 'Failed to generate plan';
                    return;
                }

                this.plan = await response.json();
            } catch (err) {
                this.error = err.message || 'An unexpected error occurred';
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
