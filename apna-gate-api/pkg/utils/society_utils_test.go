package utils

import "testing"

func TestGenerateSocietyCodeBuildsDeterministicPrefixAndPinPart(t *testing.T) {
	code, err := GenerateSocietyCode(" Green Valley ", "Pune", "MH", "411045")
	if err != nil {
		t.Fatalf("GenerateSocietyCode() error = %v", err)
	}

	wantLength := societyPrefixLen + pincodeSuffixLen + societyRandomSuffix
	if len(code) != wantLength {
		t.Fatalf("expected %d chars, got %q", wantLength, code)
	}
	if got := code[:societyPrefixLen]; got != "GREEN" {
		t.Fatalf("prefix = %q", got)
	}
	if got := code[societyPrefixLen : societyPrefixLen+pincodeSuffixLen]; got != "1045" {
		t.Fatalf("pin part = %q", got)
	}
	last := code[len(code)-1]
	if last < 'A' || last > 'Z' {
		t.Fatalf("random suffix should be uppercase letter, got %q", last)
	}
}

func TestGenerateSocietyCodePadsSparseInputs(t *testing.T) {
	code, err := GenerateSocietyCode("@", "", "", "7A")
	if err != nil {
		t.Fatalf("GenerateSocietyCode() error = %v", err)
	}

	wantLength := societyPrefixLen + pincodeSuffixLen + societyRandomSuffix
	if len(code) != wantLength {
		t.Fatalf("expected %d chars, got %q", wantLength, code)
	}
	if got := code[:societyPrefixLen]; got != "XXXXX" {
		t.Fatalf("prefix should be padded with X, got %q", got)
	}
	if got := code[societyPrefixLen : societyPrefixLen+pincodeSuffixLen]; got != "0007" {
		t.Fatalf("pin part should be left padded, got %q", got)
	}
}
