package application

import (
	"context"
	"fmt"
	"github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	repo2 "github.com/user2410/rrms-backend/internal/domain/application/repo"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Service interface {
	CreateApplication(data *dto.CreateApplicationDto) (*model.ApplicationModel, error)
	GetApplicationById(id int64) (*model.ApplicationModel, error)
	GetApplicationsByUserId(uid uuid.UUID) ([]model.ApplicationModel, error)
	GetApplicationsToUser(uid uuid.UUID) ([]model.ApplicationModel, error)
	UpdateApplicationStatus(aid int64, status database.APPLICATIONSTATUS) error
	DeleteApplication(id int64) error
}

type service struct {
	repo            repo2.Repo
	taskDistributor asynctask.TaskDistributor
}

func NewService(
	repo repo2.Repo,
	taskDistributor asynctask.TaskDistributor,
) Service {
	return &service{
		repo:            repo,
		taskDistributor: taskDistributor,
	}
}

func (s *service) CreateApplication(data *dto.CreateApplicationDto) (*model.ApplicationModel, error) {
	a, err := s.repo.CreateApplication(context.Background(), data)
	if err != nil {
		return nil, err
	}
	if err = s.taskDistributor.DistributeTaskSendEmailOnNewApplication(
		context.Background(),
		&dto.TaskSendEmailOnNewApplicationPayload{
			Username:      a.FullName,
			ApplicationId: a.ID,
			ListingId:     a.ListingID,
		},
	); err != nil {
		log.Errorf("failed to distribute DistributeTaskSendNewApplicationEmail task: %v", err)
	}
	return a, nil
}

func (s *service) GetApplicationById(id int64) (*model.ApplicationModel, error) {
	return s.repo.GetApplicationById(context.Background(), id)
}

func (s *service) DeleteApplication(id int64) error {
	return s.repo.DeleteApplication(context.Background(), id)
}

var ErrInvalidStatusTransition = fmt.Errorf("invalid status transition")

func (s *service) UpdateApplicationStatus(aid int64, status database.APPLICATIONSTATUS) error {
	a, err := s.repo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return err
	}

	switch status {
	case database.APPLICATIONSTATUSCONDITIONALLYAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING && status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSREJECTED:
		if status != database.APPLICATIONSTATUSPENDING && status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	}

	return s.repo.UpdateApplicationStatus(context.Background(), aid, status)
}

func (s *service) GetApplicationsByUserId(uid uuid.UUID) ([]model.ApplicationModel, error) {
	return s.repo.GetApplicationsByUserId(context.Background(), uid)
}

func (s *service) GetApplicationsToUser(uid uuid.UUID) ([]model.ApplicationModel, error) {
	return s.repo.GetApplicationsToUser(context.Background(), uid)
}
