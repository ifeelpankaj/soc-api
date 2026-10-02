package jobs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-server/internal/models"
)

var visitorReportHeaders = []string{
	"entry_id", "visitor_name", "visitor_phone", "visitor_email", "flat_number", "block", "floor",
	"source", "purpose", "status", "expected_at", "expected_checkout_at", "checked_in_at", "checked_out_at",
	"auto_closed_at", "vehicle_number", "vehicle_type", "companions_count", "companion_details", "approver_name", "guard_name",
	"rejection_reason", "notes", "created_at", "updated_at",
}

func GenerateVisitorReportCSV(rows []models.MonthlyVisitorReportRow) ([]byte, error) {
	// Display timezone is independent of the configured report scheduling zone.
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, fmt.Errorf("load report display timezone: %w", err)
	}
	var output bytes.Buffer
	// A UTF-8 BOM makes names and addresses display correctly in desktop Excel.
	output.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&output)
	if err := writer.Write(visitorReportHeaders); err != nil {
		return nil, err
	}
	for _, row := range rows {
		record := []string{
			strconv.FormatInt(row.ID, 10), safeCSV(row.VisitorName), safeCSV(stringValue(row.VisitorPhone)),
			safeCSV(stringValue(row.VisitorEmail)), safeCSV(stringValue(row.FlatNumber)), safeCSV(stringValue(row.Block)),
			safeCSV(stringValue(row.Floor)), safeCSV(row.Source), safeCSV(row.Purpose), safeCSV(row.Status),
			formatReportTime(row.ExpectedAt, loc), formatReportTime(row.ExpectedCheckoutAt, loc),
			formatReportTime(row.CheckedInAt, loc), formatReportTime(row.CheckedOutAt, loc), formatReportTime(row.AutoClosedAt, loc),
			safeCSV(stringValue(row.VehicleNumber)), safeCSV(stringValue(row.VehicleType)), strconv.FormatInt(int64(row.CompanionsCount), 10),
			safeCSV(row.CompanionDetails),
			safeCSV(stringValue(row.ApproverName)), safeCSV(stringValue(row.GuardName)), safeCSV(stringValue(row.RejectionReason)),
			safeCSV(stringValue(row.Notes)), formatReportTime(&row.CreatedAt, loc), formatReportTime(&row.UpdatedAt, loc),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func safeCSV(value string) string {
	trimmed := strings.TrimLeft(value, " \t\r\n")
	if trimmed == "" {
		return value
	}
	switch trimmed[0] {
	case '=', '+', '-', '@':
		return "'" + value
	default:
		return value
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func formatReportTime(value *time.Time, loc *time.Location) string {
	if value == nil || value.IsZero() {
		return ""
	}
	local := value.In(loc)
	months := [...]string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sept", "oct", "nov", "dec"}
	return fmt.Sprintf("%d %s %d %s IST", local.Day(), months[local.Month()-1], local.Year(), local.Format("15:04"))
}
