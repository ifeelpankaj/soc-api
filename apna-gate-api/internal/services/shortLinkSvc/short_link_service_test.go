package shortlinksvc

import (
	"encoding/json"
	"testing"
	"time"

	"go-server/internal/models"
)

func TestResponseUsesCanonicalLinkPath(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour)
	svc := NewShortLinkService(nil, "http://192.168.1.3:3000/")

	response := svc.Response(&models.ShortLink{
		ShortCode: "Ab3Xy9",
		ExpiresAt: &expiresAt,
	})

	if response == nil {
		t.Fatal("response is nil")
	}
	if response.URL != "http://192.168.1.3:3000/link/Ab3Xy9" {
		t.Fatalf("url = %q, want canonical /link URL", response.URL)
	}
	if response.ShortCode != "Ab3Xy9" {
		t.Fatalf("short code = %q", response.ShortCode)
	}
	if response.ExpiresAt != &expiresAt {
		t.Fatal("expires_at was not preserved")
	}

	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	if !containsJSONField(payload, "code") {
		t.Fatalf("response JSON = %s, want code field", payload)
	}
	if containsJSONField(payload, "short_code") {
		t.Fatalf("response JSON = %s, did not expect short_code field", payload)
	}
}

func containsJSONField(payload []byte, field string) bool {
	var body map[string]any
	if err := json.Unmarshal(payload, &body); err != nil {
		return false
	}
	_, ok := body[field]
	return ok
}
