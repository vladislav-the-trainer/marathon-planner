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
        },

        downloadPDF() {
            // Set the document title for the PDF filename
            const raceDate = this.formatMarathonDate(this.input.marathon_date).replace(/ /g, '-');
            const originalTitle = document.title;
            document.title = `Marathon-Training-Plan-${raceDate}`;
            
            // Trigger print dialog
            window.print();
            
            // Restore original title
            setTimeout(() => {
                document.title = originalTitle;
            }, 100);
        },

        exportPlan() {
            // Create export data with all plan information
            const now = new Date();
            const exportData = {
                version: '1.0',
                exported_at: now.toISOString(),
                input: {
                    fitness_level: this.input.fitness_level,
                    marathon_date: this.input.marathon_date,
                    target_finish_min: this.input.target_finish_min,
                    training_days_per_week: this.input.training_days_per_week
                },
                plan: this.plan
            };

            // Create and download JSON file
            const dataStr = JSON.stringify(exportData, null, 2);
            const blob = new Blob([dataStr], { type: 'application/json' });
            const url = URL.createObjectURL(blob);

            // Format race date for filename
            const raceDate = this.formatMarathonDate(this.input.marathon_date).replace(/ /g, '-');

            // Format current date and time for filename (YYYY-MM-DD_HH-MM-SS)
            const exportDate = now.getFullYear() + '-' +
                String(now.getMonth() + 1).padStart(2, '0') + '-' +
                String(now.getDate()).padStart(2, '0');
            const exportTime = String(now.getHours()).padStart(2, '0') + '-' +
                String(now.getMinutes()).padStart(2, '0') + '-' +
                String(now.getSeconds()).padStart(2, '0');

            const filename = `Marathon-Plan-${raceDate}_${exportDate}_${exportTime}.json`;

            const a = document.createElement('a');
            a.href = url;
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        },

        triggerImport() {
            document.getElementById('importFile').click();
        },

        importPlan(event) {
            const file = event.target.files[0];
            if (!file) return;

            // Validate file type
            if (!file.name.endsWith('.json')) {
                this.error = 'Please select a valid JSON file (.json)';
                return;
            }

            // Validate file size (max 1MB)
            if (file.size > 1024 * 1024) {
                this.error = 'File is too large. Maximum size is 1MB.';
                return;
            }

            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    const data = JSON.parse(e.target.result);

                    // Validate imported data structure
                    const validationResult = this.validateImportedData(data);
                    if (!validationResult.valid) {
                        this.error = validationResult.error;
                        return;
                    }

                    // Apply imported data
                    this.input.fitness_level = data.input.fitness_level;
                    this.input.marathon_date = data.input.marathon_date;
                    this.input.target_finish_min = data.input.target_finish_min;
                    this.input.training_days_per_week = data.input.training_days_per_week;
                    this.plan = data.plan;
                    this.error = '';

                } catch (err) {
                    this.error = 'Invalid file format. Please select a valid plan file.';
                }
            };

            reader.onerror = () => {
                this.error = 'Failed to read file. Please try again.';
            };

            reader.readAsText(file);

            // Reset file input so same file can be imported again
            event.target.value = '';
        },

        validateImportedData(data) {
            // Check required top-level fields
            if (!data || typeof data !== 'object') {
                return { valid: false, error: 'Invalid file structure.' };
            }

            if (!data.input || !data.plan) {
                return { valid: false, error: 'Missing required plan data.' };
            }

            // Validate input fields
            const input = data.input;

            // Validate fitness_level
            const validLevels = ['beginner', 'intermediate', 'advanced'];
            if (!validLevels.includes(input.fitness_level)) {
                return { valid: false, error: 'Invalid fitness level.' };
            }

            // Validate marathon_date format (YYYY-MM-DD)
            const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
            if (!input.marathon_date || !dateRegex.test(input.marathon_date)) {
                return { valid: false, error: 'Invalid marathon date format.' };
            }

            // Validate date is in the future (with small buffer for timezone)
            const raceDate = new Date(input.marathon_date);
            const today = new Date();
            today.setHours(0, 0, 0, 0);
            if (raceDate < today) {
                return { valid: false, error: 'Marathon date cannot be in the past.' };
            }

            // Validate target_finish_min (120-480 minutes = 2h to 8h)
            if (typeof input.target_finish_min !== 'number' ||
                input.target_finish_min < 120 ||
                input.target_finish_min > 480) {
                return { valid: false, error: 'Target time must be between 2h and 8h (120-480 minutes).' };
            }

            // Validate training_days_per_week (3-6)
            if (![3, 4, 5, 6].includes(input.training_days_per_week)) {
                return { valid: false, error: 'Training days must be 3, 4, 5, or 6.' };
            }

            // Validate plan structure
            const plan = data.plan;
            if (!plan.total_weeks || !Array.isArray(plan.weeks)) {
                return { valid: false, error: 'Invalid plan structure.' };
            }

            // Validate total_weeks (reasonable range: 1-52 weeks)
            if (plan.total_weeks < 1 || plan.total_weeks > 52) {
                return { valid: false, error: 'Invalid total weeks (must be 1-52).' };
            }

            // Validate weeks array length matches total_weeks
            if (plan.weeks.length !== plan.total_weeks) {
                return { valid: false, error: 'Weeks count mismatch.' };
            }

            // Validate each week structure
            for (let i = 0; i < plan.weeks.length; i++) {
                const week = plan.weeks[i];
                if (!week.week_number || !Array.isArray(week.sessions)) {
                    return { valid: false, error: `Invalid week ${i + 1} structure.` };
                }

                // Validate each session
                for (const session of week.sessions) {
                    if (typeof session.distance_km !== 'number' || session.distance_km < 0) {
                        return { valid: false, error: 'Invalid session distance.' };
                    }
                    if (session.distance_km > 100) {
                        return { valid: false, error: 'Unrealistic session distance (max 100km).' };
                    }
                }
            }

            return { valid: true };
        },

        get paceTable() {
            if (!this.input.target_finish_min) return [];

            const marathonDistance = 42.195;
            const marathonTimeMin = this.input.target_finish_min;

            // Distances to calculate paces for
            const distances = [
                { km: 42.195, label: 'Marathon (42.2 km)' },
                { km: 21.0975, label: 'Half Marathon (21.1 km)' },
                { km: 10, label: '10K' },
                { km: 5, label: '5K' },
                { km: 1, label: '1K' }
            ];

            // Riegel formula: T2 = T1 × (D2/D1)^1.06
            const fatigueFactor = 1.06;

            return distances.map(d => {
                // Calculate time using Riegel formula
                const timeMin = marathonTimeMin * Math.pow(d.km / marathonDistance, fatigueFactor);

                // Calculate pace (min/km)
                const paceMin = timeMin / d.km;
                const paceMinInt = Math.floor(paceMin);
                const paceSec = Math.round((paceMin - paceMinInt) * 60);
                const paceStr = `${paceMinInt}:${paceSec.toString().padStart(2, '0')}`;

                // Calculate time string
                const hours = Math.floor(timeMin / 60);
                const mins = Math.floor(timeMin % 60);
                const secs = Math.round((timeMin - Math.floor(timeMin)) * 60);
                let timeStr = '';
                if (hours > 0) {
                    timeStr = `${hours}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
                } else {
                    timeStr = `${mins}:${secs.toString().padStart(2, '0')}`;
                }

                return {
                    distance: d.label,
                    pace: `${paceStr} /km`,
                    time: timeStr
                };
            });
        }
    };
}
