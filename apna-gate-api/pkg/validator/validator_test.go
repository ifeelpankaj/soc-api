package validator

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateStructUsesJSONFieldNamesAndCustomValidators(t *testing.T) {
	type registerRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Phone    string `json:"phone" validate:"required,phone_intl"`
		Password string `json:"password" validate:"required,strong_password"`
		OTP      string `json:"otp" validate:"required,otp_code"`
		Role     string `json:"role" validate:"required,user_role"`
		Name     string `json:"name" validate:"required,alphanumeric_space"`
		Username string `json:"username" validate:"required,no_whitespace"`
	}

	valid := registerRequest{
		Email:    "owner@example.com",
		Phone:    "+911234567890",
		Password: "Strong#123",
		OTP:      "123456",
		Role:     "owner",
		Name:     "Apna Gate 1",
		Username: "apnagate",
	}
	if errs := ValidateStruct(valid); errs != nil {
		t.Fatalf("expected valid request, got %v", errs)
	}
	if !Validate(valid) {
		t.Fatal("Validate returned false for a valid struct")
	}

	invalid := registerRequest{
		Email:    "bad-email",
		Phone:    "123",
		Password: "weak",
		OTP:      "12x456",
		Role:     "moderator",
		Name:     "Apna@Gate",
		Username: "apna gate",
	}
	errs := ValidateStruct(invalid)
	if len(errs) != 7 {
		t.Fatalf("expected 7 validation errors, got %d: %v", len(errs), errs)
	}

	messages := errs.ToStringMap()
	expected := map[string]string{
		"email":    "email must be a valid email address",
		"phone":    "phone must be a valid phone number (e.g., +1234567890)",
		"password": "password must contain at least 8 characters with uppercase, lowercase, number, and special character",
		"otp":      "otp must be a valid 6-digit OTP code",
		"role":     "role must be one of: user, admin, moderator",
		"name":     "name must contain only letters, numbers, and spaces",
		"username": "username must not contain whitespace",
	}
	for field, want := range expected {
		if got := messages[field]; got != want {
			t.Fatalf("%s message mismatch: got %q want %q", field, got, want)
		}
	}
	if got := errs.ToMap()["email"]; got != expected["email"] {
		t.Fatalf("ToMap email mismatch: got %v", got)
	}
	if got := errs.Error(); !strings.Contains(got, "email: email must be a valid email address") {
		t.Fatalf("Error() did not include formatted messages: %q", got)
	}
}

func TestValidateStructHandlesNilAndInvalidInput(t *testing.T) {
	nilErrs := ValidateStruct(nil)
	if len(nilErrs) != 1 || nilErrs[0].Field != "request" {
		t.Fatalf("expected request error for nil input, got %#v", nilErrs)
	}

	invalidErrs := ValidateStruct(42)
	if len(invalidErrs) != 1 || invalidErrs[0].Field != "validation" {
		t.Fatalf("expected invalid validation error, got %#v", invalidErrs)
	}

	if got := (ValidationErrors{}).Error(); got != "" {
		t.Fatalf("empty ValidationErrors Error() = %q", got)
	}
}

func TestValidateVarAndSanitizers(t *testing.T) {
	if err := ValidateVar("admin", "user_role"); err != nil {
		t.Fatalf("expected role to pass validation: %v", err)
	}
	if err := ValidateVar("root", "user_role"); err == nil {
		t.Fatal("expected invalid role to fail validation")
	}
	if GetValidator() == nil {
		t.Fatal("expected validator instance")
	}

	if !IsValidEmail(" OWNER@Example.COM ") {
		t.Fatal("expected normalized email to validate")
	}
	if IsValidEmail("owner@example") {
		t.Fatal("expected malformed email to fail")
	}
	if !IsValidPhone("+919876543210") || IsValidPhone("09876543210") {
		t.Fatal("phone validation did not match expected E.164 behavior")
	}
	if got := SanitizeEmail(" OWNER@Example.COM "); got != "owner@example.com" {
		t.Fatalf("SanitizeEmail = %q", got)
	}
	if got := SanitizePhone(" +91 (98765)-43210 "); got != "+919876543210" {
		t.Fatalf("SanitizePhone = %q", got)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  string
	}{
		{name: "valid", password: "Strong#123"},
		{name: "too short", password: "S#1a", wantErr: "at least 8"},
		{name: "missing uppercase", password: "strong#123", wantErr: "uppercase"},
		{name: "missing lowercase", password: "STRONG#123", wantErr: "lowercase"},
		{name: "missing number", password: "StrongPass#", wantErr: "number"},
		{name: "missing special", password: "Strong123", wantErr: "special"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected password to pass, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}

	if !errors.Is(ErrInvalidEmailFormat, ErrInvalidEmailFormat) {
		t.Fatal("expected package errors to be usable as sentinels")
	}
}

func TestFormatErrorMessagesForBuiltInTags(t *testing.T) {
	type builtInRequest struct {
		RequiredWithout string `json:"required_without" validate:"required_without=OtherWithout"`
		RequiredWith    string `json:"required_with" validate:"required_with=OtherWith"`
		Email           string `json:"email" validate:"email"`
		MinString       string `json:"min_string" validate:"min=3"`
		MaxString       string `json:"max_string" validate:"max=2"`
		LenString       string `json:"len_string" validate:"len=4"`
		MinNumber       int    `json:"min_number" validate:"min=5"`
		MaxNumber       int    `json:"max_number" validate:"max=5"`
		LenNumber       int    `json:"len_number" validate:"len=4"`
		OneOf           string `json:"one_of" validate:"oneof=admin staff"`
		GTE             int    `json:"gte" validate:"gte=10"`
		LTE             int    `json:"lte" validate:"lte=10"`
		GT              int    `json:"gt" validate:"gt=10"`
		LT              int    `json:"lt" validate:"lt=10"`
		Confirm         string `json:"confirm" validate:"eqfield=Password"`
		NotSame         string `json:"not_same" validate:"nefield=Password"`
		Alpha           string `json:"alpha" validate:"alpha"`
		Alphanum        string `json:"alphanum" validate:"alphanum"`
		Numeric         string `json:"numeric" validate:"numeric"`
		URL             string `json:"url" validate:"url"`
		URI             string `json:"uri" validate:"uri"`
		Lowercase       string `json:"lowercase" validate:"lowercase"`
		Uppercase       string `json:"uppercase" validate:"uppercase"`
		Contains        string `json:"contains" validate:"contains=gate"`
		ContainsAny     string `json:"contains_any" validate:"containsany=@#"`
		Excludes        string `json:"excludes" validate:"excludes=bad"`
		StartsWith      string `json:"starts_with" validate:"startswith=APT"`
		EndsWith        string `json:"ends_with" validate:"endswith=END"`
		Password        string `json:"password"`
		OtherWithout    string `json:"other_without"`
		OtherWith       string `json:"other_with"`
	}

	errs := ValidateStruct(builtInRequest{
		RequiredWithout: "",
		RequiredWith:    "",
		Email:           "bad",
		MinString:       "ab",
		MaxString:       "abc",
		LenString:       "abc",
		MinNumber:       4,
		MaxNumber:       6,
		LenNumber:       3,
		OneOf:           "owner",
		GTE:             9,
		LTE:             11,
		GT:              10,
		LT:              10,
		Confirm:         "different",
		NotSame:         "secret",
		Alpha:           "abc1",
		Alphanum:        "abc!",
		Numeric:         "12a",
		URL:             "not-url",
		URI:             "%zz",
		Lowercase:       "LOWER",
		Uppercase:       "upper",
		Contains:        "door",
		ContainsAny:     "plain",
		Excludes:        "bad value",
		StartsWith:      "UNIT-1",
		EndsWith:        "BEGIN",
		Password:        "secret",
		OtherWithout:    "",
		OtherWith:       "present",
	})

	messages := errs.ToStringMap()
	expected := map[string]string{
		"required_without": "required_without is required when OtherWithout is not provided",
		"required_with":    "required_with is required when OtherWith is provided",
		"email":            "email must be a valid email address",
		"min_string":       "min_string must be at least 3 characters long",
		"max_string":       "max_string must not exceed 2 characters",
		"len_string":       "len_string must be exactly 4 characters long",
		"min_number":       "min_number must be at least 5",
		"max_number":       "max_number must not exceed 5",
		"len_number":       "len_number must be exactly 4",
		"one_of":           "one_of must be one of: admin, staff",
		"gte":              "gte must be greater than or equal to 10",
		"lte":              "lte must be less than or equal to 10",
		"gt":               "gt must be greater than 10",
		"lt":               "lt must be less than 10",
		"confirm":          "confirm must be equal to Password",
		"not_same":         "not_same must not be equal to Password",
		"alpha":            "alpha must contain only alphabetic characters",
		"alphanum":         "alphanum must contain only alphanumeric characters",
		"numeric":          "numeric must contain only numbers",
		"url":              "url must be a valid URL",
		"uri":              "uri must be a valid URI",
		"lowercase":        "lowercase must be in lowercase",
		"uppercase":        "uppercase must be in uppercase",
		"contains":         "contains must contain 'gate'",
		"contains_any":     "contains_any must contain at least one of: @#",
		"excludes":         "excludes must not contain 'bad'",
		"starts_with":      "starts_with must start with 'APT'",
		"ends_with":        "ends_with must end with 'END'",
	}
	for field, want := range expected {
		if got := messages[field]; got != want {
			t.Fatalf("%s message mismatch: got %q want %q", field, got, want)
		}
	}
}
