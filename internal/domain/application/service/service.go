package service

import (
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	reminder_service "github.com/user2410/rrms-backend/internal/domain/reminder"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
)

type Service interface {
	PreCreateApplication(data *dto.PreCreateApplication, creatorID uuid.UUID) error
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

	SendNotificationOnNewApplication(am *model.ApplicationModel) error
	SendNotificationOnUpdateApplication(am *model.ApplicationModel, status database.APPLICATIONSTATUS) error
}

type service struct {
	domainRepo      repos.DomainRepo
	reminderService reminder_service.Service
	miscService     misc_service.Service

	s3Client        s3.S3Client
	imageBucketName string

	asynctaskDistributor asynctask.Distributor
	feSite               string
}

func NewService(
	domainRepo repos.DomainRepo,
	reminderService reminder_service.Service,
	miscService misc_service.Service,
	s3Client s3.S3Client,
	imageBucketName string,
	asynctaskDistributor asynctask.Distributor,
	feSite string,
) Service {
	return &service{
		domainRepo:           domainRepo,
		reminderService:      reminderService,
		miscService:          miscService,
		s3Client:             s3Client,
		imageBucketName:      imageBucketName,
		asynctaskDistributor: asynctaskDistributor,
		feSite:               feSite,
	}
}
