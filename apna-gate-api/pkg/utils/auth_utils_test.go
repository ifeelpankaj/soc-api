package utils

import (
	"errors"
	"testing"
)

func TestGenerateOTPReturnsSixDigits(t *testing.T) {
	otp, err := GenerateOTP()
	if err != nil {
		t.Fatalf("GenerateOTP returned error: %v", err)
	}
	if len(otp) != 6 {
		t.Fatalf("expected 6 digits, got %q", otp)
	}
	for _, r := range otp {
		if r < '0' || r > '9' {
			t.Fatalf("OTP contains non-digit rune %q in %q", r, otp)
		}
	}
}

func TestHashPasswordAndCheckPassword(t *testing.T) {
	hashed, err := HashPassword("Secret#123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashed == "Secret#123" {
		t.Fatal("password hash should not equal the plain password")
	}
	if err := CheckPassword("Secret#123", hashed); err != nil {
		t.Fatalf("expected password to match hash: %v", err)
	}
	if err := CheckPassword("Wrong#123", hashed); err == nil {
		t.Fatal("expected wrong password to fail")
	}
}

func TestHashPasswordRejectsEmptyPassword(t *testing.T) {
	_, err := HashPassword("")
	if !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("expected ErrEmptyPassword, got %v", err)
	}
}
