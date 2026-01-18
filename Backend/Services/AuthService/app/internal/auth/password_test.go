package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.True(t, len(hash) > 0)
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "samePassword"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "bcrypt should generate different salts each time")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""

	hash, err := HashPassword(password)

	// bcrypt allows empty passwords
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPassword_CorrectPassword(t *testing.T) {
	password := "correctPassword"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	result := CheckPassword(password, hash)

	assert.True(t, result)
}

func TestCheckPassword_IncorrectPassword(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	result := CheckPassword(wrongPassword, hash)

	assert.False(t, result)
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	password := ""

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	result := CheckPassword(password, hash)

	assert.True(t, result)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	password := "password123"
	invalidHash := "not-a-valid-bcrypt-hash"

	result := CheckPassword(password, invalidHash)

	assert.False(t, result)
}

func TestCheckPassword_CaseSensitive(t *testing.T) {
	password := "Password123"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact match", "Password123", true},
		{"lowercase", "password123", false},
		{"uppercase", "PASSWORD123", false},
		{"different case", "pAssWord123", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckPassword(tc.input, hash)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHashAndCheckPassword_RoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		password string
	}{
		{"simple password", "password"},
		{"complex password", "P@ssw0rd!#$%"},
		{"long password", "this-is-a-very-long-password-with-many-characters-1234567890"},
		{"special characters", "pass!@#$%^&*()_+-=[]{}|;:',.<>?/"},
		{"unicode", "пароль密码🔒"},
		{"numbers only", "123456789"},
		{"single character", "a"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword(tc.password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)

			result := CheckPassword(tc.password, hash)
			assert.True(t, result, "Should verify the original password")

			wrongResult := CheckPassword(tc.password+"wrong", hash)
			assert.False(t, wrongResult, "Should not verify wrong password")
		})
	}
}

func TestHashPassword_ConsistentFormat(t *testing.T) {
	password := "testPassword"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	// bcrypt hashes start with "$2a$" or "$2b$" followed by cost
	assert.Regexp(t, `^\$2[ab]\$\d{2}\$.+`, hash)
}
