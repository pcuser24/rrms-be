package application

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
)

type Service interface {
	CreateApplication(data *dto.CreateApplicationDto) (*model.ApplicationModel, error)
	GetApplicationById(id int64) (*model.ApplicationModel, error)
	DeleteApplication(id int64) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateApplication(data *dto.CreateApplicationDto) (*model.ApplicationModel, error) {
	return s.repo.CreateApplication(context.Background(), data)
}

func (s *service) GetApplicationById(id int64) (*model.ApplicationModel, error) {
	return s.repo.GetApplicationById(context.Background(), id)
}

func (s *service) DeleteApplication(id int64) error {
	return s.repo.DeleteApplication(context.Background(), id)
}
