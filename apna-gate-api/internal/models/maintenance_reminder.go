package models

import "fmt"

const MaintenanceReminderType = "maintenance_payment_reminder"

// Absolute dates remain accurate when a failed notification is retried later.
func MaintenanceReminderText(month, bill, due string, outstanding int64, milestone string) (string, string) {
	title := "Maintenance payment reminder"
	if milestone == "after_7" {
		title = "Overdue maintenance reminder"
	}
	return title, fmt.Sprintf("Your %s maintenance bill %s has ₹%d.%02d outstanding. Due date: %s. Open the bill for details.", month, bill, outstanding/100, outstanding%100, due)
}
