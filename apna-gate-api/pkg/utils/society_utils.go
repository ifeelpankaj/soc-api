package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
)

const (
	societyCodeCharset  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	societyCodeLength   = 10
	societyPrefixLen    = 5
	pincodeSuffixLen    = 4
	societyRandomSuffix = 1
)

// GenerateSocietyCode generates a PAN-style society code.
//
// Format:
//   - 5 uppercase alphanumeric characters from name + city + state
//   - 4 digits from the end of the pincode
//   - 1 random uppercase letter
//
// Example:
//
//	GREEN1045K
//
// Total length: 10 characters.
//
// Callers should still handle database uniqueness constraint violations
// and retry generation when necessary.
func GenerateSocietyCode(name, city, state, pincode string) (string, error) {
	prefix := buildPrefix(name, city, state)
	pinPart := buildPinPart(pincode)

	randomPart, err := randomString(societyRandomSuffix)
	if err != nil {
		return "", fmt.Errorf("generate society code random suffix: %w", err)
	}

	code := prefix + pinPart + randomPart

	if len(code) != societyCodeLength {
		return "", fmt.Errorf(
			"invalid society code length: got %d, expected %d",
			len(code),
			societyCodeLength,
		)
	}

	return code, nil
}

// buildPrefix extracts the first 5 uppercase alphanumeric characters
// from name + city + state.
//
// If fewer than 5 valid characters are available, the result is
// right-padded with 'X'.
func buildPrefix(name, city, state string) string {
	combined := strings.ToUpper(
		strings.TrimSpace(name) +
			strings.TrimSpace(city) +
			strings.TrimSpace(state),
	)

	var b strings.Builder
	b.Grow(societyPrefixLen)

	for _, r := range combined {
		if isASCIIAlphaNumeric(r) {
			b.WriteRune(r)

			if b.Len() == societyPrefixLen {
				break
			}
		}
	}

	for b.Len() < societyPrefixLen {
		b.WriteByte('X')
	}

	return b.String()
}

// buildPinPart extracts only digits from the pincode and returns
// the last 4 digits.
//
// If fewer than 4 digits are available, the result is left-padded
// with zeroes.
//
// Examples:
//
//	411045 -> 1045
//	123    -> 0123
//	7      -> 0007
//	""     -> 0000
func buildPinPart(pincode string) string {
	var b strings.Builder

	for _, r := range pincode {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	digits := b.String()

	if len(digits) >= pincodeSuffixLen {
		return digits[len(digits)-pincodeSuffixLen:]
	}

	return strings.Repeat("0", pincodeSuffixLen-len(digits)) + digits
}

// randomString generates a cryptographically secure random uppercase
// alphabetic string of the requested length.
func randomString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	result := make([]byte, length)

	for i, value := range randomBytes {
		result[i] = societyCodeCharset[int(value)%len(societyCodeCharset)]
	}

	return string(result), nil
}

func isASCIIAlphaNumeric(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}
