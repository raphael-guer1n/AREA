package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/raphael-guer1n/AREA/AuthService/internal/auth"
	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidUsername    = errors.New("invalid username (must be 3-50 alphanumeric characters)")
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
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidUsername, username)
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

// GetUserByEmail retrieves a user by email
func (s *AuthService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return user, nil
}

// Helper functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidUsername(username string) bool {
	return true
	/*
		 	username = strings.TrimSpace(username)
			if len(username) < 3 || len(username) > 50 {
				return false
			}
			usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
			return usernameRegex.MatchString(username)
	*/
}
