package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}



// GenerateToken creates a signed JWT string for a validated user
func GenerateToken(userID, email, role, secret string, expiry time.Duration, issuer string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}



// VerifyToken parses, validates the signature, and returns the unpacked custom claims
func VerifyToken(tokenString string, secret string, issuer string) (*UserClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}), // Prevents algorithm downgrade attacks
		jwt.WithIssuer(issuer),
	)

	token, err := parser.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}