package service

import (
	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/application/repo"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"

	reminder_service "github.com/user2410/rrms-backend/internal/domain/reminder"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
)

type Service interface {
	CreateApplication(data *dto.CreateApplication) (*model.ApplicationModel, error)
	GetApplicationById(id int64) (*model.ApplicationModel, error)
	GetApplicationByIds(ids []int64, fields []string, userId uuid.UUID) ([]model.ApplicationModel, error)
	GetApplicationsByUserId(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error)
	GetApplicationsToUser(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error)
	UpdateApplicationStatus(aid int64, userId uuid.UUID, data *dto.UpdateApplicationStatus) error
	CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error)
	CheckApplicationUpdatability(aid int64, uid uuid.UUID) (bool, error)
	CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error)
	GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error)
	GetRentalByApplicationId(aid int64) (rental_model.RentalModel, error)
}

type service struct {
	aRepo           repo.Repo
	cRepo           chat_repo.Repo
	lRepo           listing_repo.Repo
	pRepo           property_repo.Repo
	rService        reminder_service.Service
	taskDistributor application_asynctask.TaskDistributor
}

func NewService(
	aRepo repo.Repo,
	cRepo chat_repo.Repo,
	lRepo listing_repo.Repo,
	pRepo property_repo.Repo,
	rService reminder_service.Service,
	taskDistributor application_asynctask.TaskDistributor,
) Service {
	return &service{
		aRepo:           aRepo,
		cRepo:           cRepo,
		lRepo:           lRepo,
		pRepo:           pRepo,
		rService:        rService,
		taskDistributor: taskDistributor,
	}
}
