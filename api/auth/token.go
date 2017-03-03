package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// Token types
const (
	TokenTypeVerification = "verification"
	TokenTypeLogin        = "login"
	TokenTypePassword     = "password"
)

// AccountClaims represents account claims
type AccountClaims struct {
	AccountName        string `json:"AccountName"`
	ActiveOrganization string `json:"ActiveOrganization"`
	Type               string `json:"Type"`
	jwt.StandardClaims
}

// CreateVerificationToken creates a verification token for a given user name
func CreateVerificationToken(name string, validFor time.Duration) (string, error) {
	return createToken(name, "", TokenTypeVerification, validFor)
}

// CreateLoginToken creates a login token for a given user name
func CreateLoginToken(name string, activeOrganization string, validFor time.Duration) (string, error) {
	return createToken(name, activeOrganization, TokenTypeLogin, validFor)
}

// CreatePasswordToken creates a password token for a given user name
func CreatePasswordToken(name string, validFor time.Duration) (string, error) {
	return createToken(name, "", TokenTypePassword, validFor)
}

// createToken creates a token for a given user name
func createToken(account string, activeOrganization string, tokenType string, validFor time.Duration) (string, error) {
	// Forge the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccountClaims{
		account,
		activeOrganization,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validFor).Unix(),
			Issuer:    os.Args[0],
		},
	})
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("unable to issue token")
	}
	return ss, nil
}

// ValidateToken validates a token and return its claims
func ValidateToken(signedString string, tokenType string) (*AccountClaims, error) {
	token, err := jwt.ParseWithClaims(signedString, &AccountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*AccountClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	if claims.Type != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}
