package logger

import (
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	value := Sanitize("Bearer private-token password=private-password https://user:private-password@host\n" + strings.Repeat("界", 3000))
	if strings.Contains(value, "private") || strings.Contains(value, "\n") {
		t.Fatal("credential/control character leaked")
	}
	if len([]rune(value)) != 2048 {
		t.Fatal("message not bounded")
	}
}
