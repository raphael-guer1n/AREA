package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/raphael-guer1n/AREA/AuthService/internal/auth"
	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidUsername    = errors.New("invalid username (must be 3-20 alphanumeric characters)")
	ErrInvalidPassword    = errors.New("invalid password (must be at least 6 characters)")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	repo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register creates a new user with validation
func (s *AuthService) Register(email, username, password string) (*domain.User, string, error) {
	// Validate email format
	if !isValidEmail(email) {
		return nil, "", ErrInvalidEmail
	}

	// Validate username (3-20 alphanumeric characters)
	if !isValidUsername(username) {
		return nil, "", ErrInvalidUsername
	}

	// Validate password (minimum 6 characters)
	if len(password) < 6 {
		return nil, "", ErrInvalidPassword
	}

	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("error checking email: %w", err)
	}
	if existingUser != nil {
		return nil, "", ErrEmailAlreadyExists
	}

	// Check if username already exists
	existingUser, err = s.repo.FindByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("error checking username: %w", err)
	}
	if existingUser != nil {
		return nil, "", ErrUsernameExists
	}

	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user, err := s.repo.Create(email, username, passwordHash)
	if err != nil {
		return nil, "", fmt.Errorf("error creating user: %w", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(emailOrUsername, password string) (*domain.User, string, error) {
	// Find user by email or username
	user, err := s.repo.FindByEmailOrUsername(emailOrUsername)
	if err != nil {
		return nil, "", fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	return user, token, nil
}

// LoginWithOAuth creates or retrieves a user from OAuth2 user info and returns a JWT
func (s *AuthService) LoginWithOAuth(provider string, userInfo *oauth2.UserInfo) (*domain.User, string, error) {
	if userInfo == nil {
		return nil, "", fmt.Errorf("user info is required")
	}

	email := strings.TrimSpace(userInfo.Email)
	if email == "" {
		base := sanitizeIdentifier(userInfo.ID)
		if base == "" {
			randomID, err := generateRandomString(12)
			if err != nil {
				return nil, "", fmt.Errorf("failed to generate identifier: %w", err)
			}
			base = randomID
		}
		email = fmt.Sprintf("%s@%s.oauth", base, provider)
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("error finding user by email: %w", err)
	}

	if user == nil {
		username, err := s.generateUniqueUsername(userInfo, provider)
		if err != nil {
			return nil, "", err
		}

		randomPassword, err := generateRandomString(32)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate password: %w", err)
		}

		passwordHash, err := auth.HashPassword(randomPassword)
		if err != nil {
			return nil, "", fmt.Errorf("error hashing password: %w", err)
		}

		user, err = s.repo.Create(email, username, passwordHash)
		if err != nil {
			return nil, "", fmt.Errorf("error creating user: %w", err)
		}
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	return user, token, nil
}

// GetUserByID retrieves a user by ID (for /auth/me endpoint)
func (s *AuthService) GetUserByID(id int) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Helper functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

func (s *AuthService) generateUniqueUsername(userInfo *oauth2.UserInfo, provider string) (string, error) {
	base := sanitizeIdentifier(userInfo.Username)
	if base == "" {
		base = sanitizeIdentifier(userInfo.Name)
	}
	if base == "" {
		base = sanitizeIdentifier(userInfo.Email)
	}
	if base == "" {
		base = sanitizeIdentifier(fmt.Sprintf("%s_%s", provider, userInfo.ID))
	}
	if base == "" {
		base = provider
	}
	base = strings.ToLower(base)
	if len(base) < 3 {
		base = fmt.Sprintf("%s_user", provider)
	}
	if len(base) > 20 {
		base = base[:20]
	}

	candidate := base
	suffix := 1

	for {
		existing, err := s.repo.FindByUsername(candidate)
		if err != nil {
			return "", fmt.Errorf("error checking username availability: %w", err)
		}
		if existing == nil {
			return candidate, nil
		}

		suffixStr := fmt.Sprintf("%d", suffix)
		maxBaseLen := 20 - len(suffixStr)
		if maxBaseLen < 3 {
			maxBaseLen = 3
		}
		truncated := candidate
		if len(truncated) > maxBaseLen {
			truncated = truncated[:maxBaseLen]
		}
		candidate = fmt.Sprintf("%s%s", truncated, suffixStr)
		suffix++
	}
}

func sanitizeIdentifier(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var builder strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '-' || r == ' ' || r == '_' || r == '.':
			builder.WriteRune('_')
		}
	}

	result := strings.Trim(builder.String(), "_")
	if len(result) > 20 {
		result = result[:20]
	}
	return result
}

func generateRandomString(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = alphabet[int(bytes[i])%len(alphabet)]
	}

	return string(bytes), nil
}
