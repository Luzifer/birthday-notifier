package dateutil

import "time"

// IsToday uses ProjectToNextBirthday to get the next birthday and
// compares it to TodayStartOfDay
func IsToday(t time.Time) bool {
	return ProjectToNextBirthday(t).
		Equal(TodayStartOfDay())
}

// ProjectToNextBirthday takes a birth date and projects it to the
// next birthday being today or later
func ProjectToNextBirthday(t time.Time) time.Time {
	projected := time.Date(time.Now().Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	if projected.Before(TodayStartOfDay()) {
		projected = time.Date(time.Now().Year()+1, t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	}
	return projected
}

// TodayStartOfDay gets the start of the current day
func TodayStartOfDay() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
}
