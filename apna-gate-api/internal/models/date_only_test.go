package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateOnlyMarshalJSON(t *testing.T) {
	t.Parallel()

	dob := time.Date(1990, time.January, 15, 0, 0, 0, 0, time.UTC)
	response := UserResponse{
		ID:          1,
		FullName:    "Test User",
		DateOfBirth: NewDateOnly(&dob),
		Timezone:    "Asia/Kolkata",
		Language:    "en",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	got, ok := decoded["date_of_birth"].(string)
	if !ok {
		t.Fatalf("expected date_of_birth string, got %#v", decoded["date_of_birth"])
	}
	if got != "1990-01-15" {
		t.Fatalf("got %q want %q", got, "1990-01-15")
	}
}
