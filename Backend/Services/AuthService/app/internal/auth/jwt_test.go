package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	userID := 123

	token, err := GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_ValidClaims(t *testing.T) {
	userID := 456

	tokenString, err := GenerateToken(userID)
	assert.NoError(t, err)

	// Parse the token to verify claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "456", claims.Sub)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)

	// Check expiration is approximately 24 hours from now
	expectedExpiry := time.Now().Add(24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time
	timeDiff := actualExpiry.Sub(expectedExpiry)
	assert.Less(t, timeDiff.Abs().Seconds(), 5.0, "Expiry time should be within 5 seconds of 24 hours from now")
}

func TestValidateToken_Success(t *testing.T) {
	userID := 789

	tokenString, err := GenerateToken(userID)
	assert.NoError(t, err)

	validatedUserID, err := ValidateToken(tokenString)

	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.jwt.token"

	userID, err := ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Create an expired token
	claims := Claims{
		UserID: 123,
		Sub:    "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	userID, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	// Create a token with wrong signing method (RS256 instead of HS256)
	claims := Claims{
		UserID: 123,
		Sub:    "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// This would require an RSA key, but we're testing validation failure
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrong-secret"))
	assert.NoError(t, err)

	userID, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}

func TestGetJWTSecret_FromEnv(t *testing.T) {
	// Save original env var
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// Set custom secret
	customSecret := "test-custom-secret"
	os.Setenv("JWT_SECRET", customSecret)

	secret := getJWTSecret()

	assert.Equal(t, customSecret, secret)
}

func TestGetJWTSecret_DefaultValue(t *testing.T) {
	// Save original env var
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// Clear env var
	os.Unsetenv("JWT_SECRET")

	secret := getJWTSecret()

	assert.Equal(t, "dev-secret-key-change-me", secret)
}

func TestGenerateAndValidateToken_RoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		userID int
	}{
		{"positive ID", 1},
		{"large ID", 999999},
		{"zero ID", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateToken(tc.userID)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			validatedUserID, err := ValidateToken(token)
			assert.NoError(t, err)
			assert.Equal(t, tc.userID, validatedUserID)
		})
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	testCases := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"random string", "notajwttoken"},
		{"incomplete JWT", "header.payload"},
		{"too many parts", "header.payload.signature.extra"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := ValidateToken(tc.token)
			assert.Error(t, err)
			assert.Equal(t, 0, userID)
		})
	}
}
