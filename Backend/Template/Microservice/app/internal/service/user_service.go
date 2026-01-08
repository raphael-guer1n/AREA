package service

import (
	"github.com/raphael-guer1n/AREA/Template/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.List()
}

func (s *UserService) CreateUser(email, firstName, lastName string) (domain.User, error) {
	return s.repo.Create(email, firstName, lastName)
}
