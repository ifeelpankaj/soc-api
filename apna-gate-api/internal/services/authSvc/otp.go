package authsvc

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

func HashOTP(otp string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(otp))
	return hex.EncodeToString(mac.Sum(nil))
}

func CompareOTP(rawOTP, otpHash, secret string) bool {
	expected := HashOTP(rawOTP, secret)
	return hmac.Equal([]byte(expected), []byte(otpHash))
}

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	const otpLength = 6
	const digits = "0123456789"

	otp := make([]byte, otpLength)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
