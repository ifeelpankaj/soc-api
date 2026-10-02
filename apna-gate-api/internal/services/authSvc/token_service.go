package authsvc

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	SessionVersion int64  `json:"session_version"`
	UserID         int64  `json:"user_id"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Role           string `json:"role"`
	EmailVerified  bool   `json:"email_verified"`
	PhoneVerified  bool   `json:"phone_verified"`
	TokenType      string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateToken(
	userID int64,
	email string,
	phoneNumber string,
	role string,
	emailVerified bool,
	phoneVerified bool,
	tokenType string,
	secret string,
	issuer string,
	expiry time.Duration,
	versions ...int64,
) (string, error) {
	now := time.Now()
	var version int64
	if len(versions) > 0 {
		version = versions[0]
	}

	claims := Claims{
		SessionVersion: version,
		UserID:         userID,
		Email:          email,
		PhoneNumber:    phoneNumber,
		Role:           role,
		EmailVerified:  emailVerified,
		PhoneVerified:  phoneVerified,
		TokenType:      tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string, issuer string, expectedTokenType string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken.WithCause(err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != issuer {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedTokenType {
		return nil, ErrWrongTokenType
	}

	return claims, nil
}
