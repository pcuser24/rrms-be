package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

const (
	MAX_IMAGE_SIZE      = 10 * 1024 * 1024 // 10MB
	UPLOAD_URL_LIFETIME = 5                // 5 minutes
)

type Service interface {
	PreCreateRental(data *dto.PreCreateRental, creatorID uuid.UUID) error
	CreateRental(data *dto.CreateRental, creatorID uuid.UUID) (model.RentalModel, error)
	GetRental(id int64) (model.RentalModel, error)
	UpdateRental(data *dto.UpdateRental, id int64) error
	// PrepareRentalContract(id int64, data *dto.PrepareRentalContract) (*model.RentalContractModel, error)
	GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error)
	GetMyRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error)
	CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error)

	CreateContract(data *dto.CreateContract) (*model.ContractModel, error)
	GetRentalContractsOfUser(userId uuid.UUID, query *dto.GetRentalContracts) ([]model.ContractModel, error)
	GetRentalContract(id int64) (*model.ContractModel, error)
	PingRentalContract(id int64) (any, error)
	GetContract(id int64) (*model.ContractModel, error)
	UpdateContract(data *dto.UpdateContract) error
	UpdateContractContent(data *dto.UpdateContractContent) error

	CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error)
	GetRentalPayment(id int64) (model.RentalPayment, error)
	GetPaymentsOfRental(id int64) ([]model.RentalPayment, error)
	GetManagedRentalPayments(uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) ([]dto.GetManagedRentalPaymentsItem, error)
	UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error

	PreCreateRentalComplaint(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error
	CreateRentalComplaint(data *dto.CreateRentalComplaint) (model.RentalComplaint, error)
	GetRentalComplaint(id int64) (model.RentalComplaint, error)
	GetRentalComplaintsByRentalId(rid int64) ([]model.RentalComplaint, error)
	PreCreateRentalComplaintReply(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error
	CreateRentalComplaintReply(data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error)
	GetRentalComplaintsOfUser(userId uuid.UUID, query dto.GetRentalComplaintsOfUserQuery) ([]model.RentalComplaint, error)
	GetRentalComplaintReplies(id int64, limit, offset int32) ([]model.RentalComplaintReply, error)
	UpdateRentalComplaint(data *dto.UpdateRentalComplaint) error
}

type service struct {
	domainRepo repos.DomainRepo

	mService misc_service.Service

	cronEntries []cron.EntryID

	s3Client        s3.S3Client
	imageBucketName string
}

func NewService(
	domainRepo repos.DomainRepo,
	mService misc_service.Service,
	c *cron.Cron,
	s3Client s3.S3Client, imageBucketName string,
) Service {
	res := &service{
		domainRepo:      domainRepo,
		cronEntries:     []cron.EntryID{},
		mService:        mService,
		s3Client:        s3Client,
		imageBucketName: imageBucketName,
	}
	res.setupCronjob(c)
	return res
}

// SetupCronjob periodically checks for rental payment due date and send reminder to user
func (s *service) setupCronjob(c *cron.Cron) ([]cron.EntryID, error) {
	var (
		entryID cron.EntryID
		err     error
	)
	entryID, err = c.AddFunc("@daily", func() {
		// TODO: log any error
		// plan rental payments
		s.domainRepo.RentalRepo.PlanRentalPayments(context.Background())
		// update fine payments
		s.domainRepo.RentalRepo.UpdateFinePayments(context.Background())
	})
	if err != nil {
		return nil, err
	}
	s.cronEntries = append(s.cronEntries, entryID)

	return s.cronEntries, nil
}
