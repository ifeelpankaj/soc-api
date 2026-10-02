package jobs

import (
	"bytes"
	"encoding/csv"
	"testing"
	"time"

	"go-server/internal/models"
)

func TestGenerateVisitorReportCSVIsExcelSafeAndHandlesEmptyReports(t *testing.T) {
	empty, err := GenerateVisitorReportCSV(nil)
	if err != nil {
		t.Fatal(err)
	}
	emptyRecords := parseReportCSV(t, empty)
	if len(emptyRecords) != 1 {
		t.Fatalf("empty report has %d records, want header only", len(emptyRecords))
	}

	createdAt := time.Date(2026, time.August, 31, 20, 0, 0, 0, time.UTC)
	phone := "+911234567890"
	content, err := GenerateVisitorReportCSV([]models.MonthlyVisitorReportRow{{
		ID: 7, VisitorName: "=HYPERLINK(\"bad\")", VisitorPhone: &phone,
		Source: "guard_entry", Purpose: "guest", Status: "checked_out", CompanionDetails: "@malicious",
		CreatedAt: createdAt, UpdatedAt: createdAt,
	}})
	if err != nil {
		t.Fatal(err)
	}
	records := parseReportCSV(t, content)
	if len(records) != 2 {
		t.Fatalf("report has %d records, want 2", len(records))
	}
	if records[1][1][0] != '\'' || records[1][2][0] != '\'' {
		t.Fatalf("spreadsheet formula fields were not escaped: %#v", records[1][:3])
	}
	if records[1][18][0] != '\'' {
		t.Fatalf("companion details were not formula escaped: %q", records[1][18])
	}
	if records[1][23] != "1 sept 2026 01:30 IST" {
		t.Fatalf("created_at = %q, want timezone-formatted timestamp", records[1][23])
	}
}

func parseReportCSV(t *testing.T, content []byte) [][]string {
	t.Helper()
	if !bytes.HasPrefix(content, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatal("CSV is missing UTF-8 BOM")
	}
	records, err := csv.NewReader(bytes.NewReader(content[3:])).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestAllReportDatesUseReadableIST(t *testing.T) {
	for _, tc := range []struct {
		value time.Time
		want  string
	}{
		{time.Date(2026, 9, 10, 12, 15, 59, 0, time.UTC), "10 sept 2026 17:45 IST"},
		{time.Date(2026, 12, 31, 20, 0, 0, 0, time.UTC), "1 jan 2027 01:30 IST"},
	} {
		v := tc.value
		content, err := GenerateVisitorReportCSV([]models.MonthlyVisitorReportRow{{ExpectedAt: &v, ExpectedCheckoutAt: &v, CheckedInAt: &v, CheckedOutAt: &v, AutoClosedAt: &v, CreatedAt: v, UpdatedAt: v}})
		if err != nil {
			t.Fatal(err)
		}
		record := parseReportCSV(t, content)[1]
		for _, index := range []int{10, 11, 12, 13, 14, 23, 24} {
			if record[index] != tc.want {
				t.Errorf("%s=%q want %q", visitorReportHeaders[index], record[index], tc.want)
			}
		}
	}
	zero := time.Time{}
	if formatReportTime(nil, nil) != "" || formatReportTime(&zero, nil) != "" {
		t.Fatal("missing dates must be blank")
	}
}
