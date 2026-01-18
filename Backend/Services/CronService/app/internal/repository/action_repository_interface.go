package repository

import "github.com/raphael-guer1n/AREA/CronService/internal/domain"

type ActionRepositoryInterface interface {
	Create(action *domain.Action) error
	GetByActionID(actionID int) (*domain.Action, error)
	GetAll() ([]*domain.Action, error)
	Update(action *domain.Action) error
	Delete(actionID int) error
}
