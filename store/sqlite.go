package store

import (
	"database/sql"
	"os"
	"path"
	"time"

	"github.com/lukasmwerner/timecard/util"
	_ "modernc.org/sqlite"
)

const (
	ClockIn  = 0
	ClockOut = 1
)

type TimeEntry struct {
	Time        time.Time
	Kind        int
	Description string
}

type DB struct {
	db *sql.DB
}

func Open() (*DB, error) {

	homedir, _ := os.UserHomeDir()
	storeLoc := path.Join(homedir, ".config", "timecard")
	util.EnsureDirExists(storeLoc)

	db, err := sql.Open("sqlite", path.Join(storeLoc, "timesheet.db"))

	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS timesheet (time TEXT, kind INTEGER, description TEXT);`)

	return &DB{
		db: db,
	}, nil

}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) PunchIn(description string) error {
	_, err := db.db.Exec(`INSERT INTO timesheet (time, kind, description) VALUES (datetime('now', 'localtime'), ?, ?)`, ClockIn, description)
	return err
}

func (db *DB) PunchOut(description string) error {
	_, err := db.db.Exec(`INSERT INTO timesheet (time, kind, description) VALUES (datetime('now', 'localtime'), ?, ?)`, ClockOut, description)
	return err
}

func (db *DB) getLastEntry() (*TimeEntry, error) {
	query := `
        SELECT time, kind, description
        FROM timesheet
        ORDER BY time DESC
        LIMIT 1
    `

	var timeStr string
	entry := &TimeEntry{}

	err := db.db.QueryRow(query).Scan(&timeStr, &entry.Kind, &entry.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entry.Time, err = time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (db *DB) Status() (bool, *time.Time, error) {
	entry, err := db.getLastEntry()
	if err != nil {
		return false, nil, err
	}

	if entry == nil {
		return false, nil, nil
	}

	isClockedIn := entry.Kind == ClockIn
	if isClockedIn {
		return true, &entry.Time, nil
	}
	return false, nil, nil
}

// DailyHours represents the total hours worked for a specific day
type DailyHours struct {
	Date         time.Time
	Hours        float64
	ClockInCount int     // Number of clock-ins
	TotalBreaks  float64 // Total break time in hours
}

func (db *DB) Report() ([]DailyHours, float64, error) {
	// Query to get all clock-in and clock-out events from the last 7 days
	query := `
        WITH RECURSIVE dates(date) AS (
            SELECT date('now', 'start of day', '-6 days')
            UNION ALL
            SELECT date(date, '+1 day')
            FROM dates
            WHERE date < date('now', 'start of day')
        )
        SELECT 
            dates.date,
            timesheet.time,
            timesheet.kind,
            timesheet.description
        FROM dates
        LEFT JOIN timesheet ON date(timesheet.time) = dates.date
        WHERE timesheet.time >= datetime('now', '-7 days')
        ORDER BY timesheet.time ASC;
    `

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Map to store daily hours and statistics
	dailyStats := make(map[string]*DailyHours)
	var lastClockIn time.Time
	var lastBreakStart time.Time
	var weeklyTotal float64

	// Initialize the map with all dates
	startDate := time.Now().AddDate(0, 0, -6)
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		dailyStats[dateStr] = &DailyHours{
			Date:         date,
			Hours:        0,
			ClockInCount: 0,
			TotalBreaks:  0,
		}
	}

	for rows.Next() {
		var dateStr, timeStr string
		var kind int
		var description string
		if err := rows.Scan(&dateStr, &timeStr, &kind, &description); err != nil {
			return nil, 0, err
		}

		eventTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return nil, 0, err
		}

		currentDate := eventTime.Format("2006-01-02")
		stats := dailyStats[currentDate]

		if kind == 0 { // Clock-in
			lastClockIn = eventTime
			stats.ClockInCount++
			if !lastBreakStart.IsZero() {
				// Calculate break time
				breakDuration := eventTime.Sub(lastBreakStart).Hours()
				stats.TotalBreaks += breakDuration
				lastBreakStart = time.Time{}
			}
		} else if kind == 1 && !lastClockIn.IsZero() { // Clock-out
			duration := eventTime.Sub(lastClockIn).Hours()
			stats.Hours += duration
			weeklyTotal += duration
			lastBreakStart = eventTime
			lastClockIn = time.Time{} // Reset lastClockIn
		}
	}

	// Convert map to sorted slice
	var result []DailyHours
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		if stats, exists := dailyStats[dateStr]; exists {
			result = append(result, *stats)
		}
	}

	return result, weeklyTotal, nil
}
