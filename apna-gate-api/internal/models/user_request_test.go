package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUpdateUserRequestUnmarshalDateOfBirth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
		want    string
		wantErr bool
	}{
		{name: "iso date", payload: `{"date_of_birth":"2000-04-04"}`, want: "2000-04-04"},
		{name: "slash date", payload: `{"date_of_birth":"2000/4/4"}`, want: "2000-04-04"},
		{name: "rfc3339 date", payload: `{"date_of_birth":"1990-01-15T00:00:00Z"}`, want: "1990-01-15"},
		{name: "invalid date", payload: `{"date_of_birth":"not-a-date"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateUserRequest
			err := json.Unmarshal([]byte(tt.payload), &req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.DateOfBirth == nil {
				t.Fatalf("expected date_of_birth")
			}
			got := req.DateOfBirth.UTC().Format("2006-01-02")
			if got != tt.want {
				t.Fatalf("got %s want %s", got, tt.want)
			}
		})
	}
}

func TestParseFlexibleDate(t *testing.T) {
	t.Parallel()

	parsed, err := ParseFlexibleDate("2000/4/4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed.Format("2006-01-02") != "2000-04-04" {
		t.Fatalf("unexpected parsed date: %s", parsed.Format(time.RFC3339))
	}
}
