package service

import (
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"

	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	reminder_service "github.com/user2410/rrms-backend/internal/domain/reminder"

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
}

type service struct {
	applicationRepo application_repo.Repo
	authRepo        auth_repo.Repo
	chatRepo        chat_repo.Repo
	listingRepo     listing_repo.Repo
	propertyRepo    property_repo.Repo
	unitRepo        unit_repo.Repo
	reminderService reminder_service.Service
	miscService     misc_service.Service

	s3Client        s3.S3Client
	imageBucketName string

	feSite string
}

func NewService(
	applicationRepo application_repo.Repo,
	authRepo auth_repo.Repo,
	chatRepo chat_repo.Repo,
	listingRepo listing_repo.Repo,
	propertyRepo property_repo.Repo,
	unitRepo unit_repo.Repo,
	reminderService reminder_service.Service,
	miscService misc_service.Service,
	s3Client s3.S3Client,
	imageBucketName string,
	feSite string,
) Service {
	return &service{
		applicationRepo: applicationRepo,
		authRepo:        authRepo,
		chatRepo:        chatRepo,
		listingRepo:     listingRepo,
		propertyRepo:    propertyRepo,
		unitRepo:        unitRepo,
		reminderService: reminderService,
		miscService:     miscService,
		s3Client:        s3Client,
		imageBucketName: imageBucketName,
		feSite:          feSite,
	}
}
