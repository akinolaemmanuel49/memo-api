package helpers

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateTokens generates both the detailed token and refresh token
func GenerateTokens(jwtSecret string, userID string) (signedAccessToken string, signedRefreshToken string, err error) {
	// generate access token with a lifetime of 24 hours
	accessClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	// generate refresh token with a lifetime of 1 week
	refreshClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// ValidateToken validates the provided JWT.
// An error is returned if the token is invalid or expired.
func ValidateToken(jwtSecret string, signedToken string) (*jwt.StandardClaims, error) {
	// attempt to parse token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	// extract claims from token
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
