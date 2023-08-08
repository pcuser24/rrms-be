package property

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
)

type Service interface {
	CheckOwnership(id uuid.UUID, userId uuid.UUID) (bool, error)
	CreateProperty(data *dto.CreateProperty) (*model.PropertyModel, error)
	GetPropertyById(id uuid.UUID) (*model.PropertyModel, error)
	UpdateProperty(data *dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateProperty(data *dto.CreateProperty) (*model.PropertyModel, error) {
	fmt.Println("service: CreateProperty")
	return s.repo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*model.PropertyModel, error) {
	return s.repo.GetPropertyById(context.Background(), id)
}

func (s *service) UpdateProperty(data *dto.UpdateProperty) error {
	return s.repo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckOwnership(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return s.repo.CheckOwnership(context.Background(), id, userId)
}

func (s *service) DeleteProperty(id uuid.UUID) error {
	return s.repo.DeleteProperty(context.Background(), id)
}
