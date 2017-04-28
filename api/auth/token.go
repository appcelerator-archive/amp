package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token constants
const (
	TokenIssuer           = "amplifier"
	TokenTypeVerification = "verification"
	TokenTypeLogin        = "login"
	TokenTypePassword     = "password"

	VerificationTokenValidFor = time.Hour
	LoginTokenValidFor        = 24 * time.Hour
	PasswordTokenValidFor     = time.Hour
)

// AuthClaims represents authentication claims
type AuthClaims struct {
	Type        string `json:"Type"`
	AccountName string `json:"AccountName"`
	jwt.StandardClaims
}

// LoginClaims represents login claims
type LoginClaims struct {
	ActiveOrganization string `json:"ActiveOrganization"`
	AuthClaims
}

// CreateVerificationToken creates a verification token for a given user
func CreateVerificationToken(name string) (string, error) {
	claims := AuthClaims{
		Type:        TokenTypeVerification,
		AccountName: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(VerificationTokenValidFor).Unix(),
			Issuer:    TokenIssuer,
		},
	}
	return createToken(claims)
}

// CreateLoginToken creates a login token for a given account
func CreateLoginToken(name string, activeOrganization string) (string, error) {
	claims := LoginClaims{
		ActiveOrganization: activeOrganization,
		AuthClaims: AuthClaims{
			Type:        TokenTypeLogin,
			AccountName: name,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(LoginTokenValidFor).Unix(),
				Issuer:    TokenIssuer,
			},
		},
	}
	return createToken(claims)
}

// CreatePasswordToken creates a password token for a given user name
func CreatePasswordToken(name string) (string, error) {
	claims := AuthClaims{
		Type:        TokenTypePassword,
		AccountName: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(PasswordTokenValidFor).Unix(),
			Issuer:    TokenIssuer,
		},
	}
	return createToken(claims)
}

// createToken creates a token for a given user name
func createToken(claims jwt.Claims) (string, error) {
	// Forge the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("unable to issue token")
	}
	return ss, nil
}

// ValidateToken validates a token and return its claims
func ValidateToken(signedString string, tokenType string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(signedString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	if claims.Type != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}
