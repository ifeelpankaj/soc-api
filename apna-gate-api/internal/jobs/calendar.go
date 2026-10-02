package jobs

import "time"

// ReportMonth returns the canonical first day of timestamp's calendar month in loc.
// The returned time is a date value represented at midnight UTC for PostgreSQL DATE.
func ReportMonth(timestamp time.Time, loc *time.Location) time.Time {
	local := timestamp.In(loc)
	return time.Date(local.Year(), local.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// MonthBounds converts a canonical report month into its local calendar bounds in UTC.
func MonthBounds(reportMonth time.Time, loc *time.Location) (time.Time, time.Time) {
	startLocal := time.Date(reportMonth.Year(), reportMonth.Month(), 1, 0, 0, 0, 0, loc)
	return startLocal.UTC(), startLocal.AddDate(0, 1, 0).UTC()
}
