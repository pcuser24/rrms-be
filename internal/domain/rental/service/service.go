package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/utils"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

const (
	MAX_IMAGE_SIZE      = 10 * 1024 * 1024 // 10MB
	UPLOAD_URL_LIFETIME = 5                // 5 minutes
)

type Service interface {
	PreCreateRental(data *dto.PreCreateRental, creatorID uuid.UUID) error
	CreatePreRental(data *dto.CreateRental, creatorID uuid.UUID) (rental_model.PreRental, error)
	// CreateRental(data *dto.CreateRental, creatorID uuid.UUID) (rental_model.RentalModel, error)
	GetPreRentalExtended(id int64, userId uuid.UUID, key string) (dto.GetPreRentalResponse, error)
	GetPreRentalByID(id int64) (rental_model.PreRental, error)
	GetPreRentalsToMe(userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]rental_model.PreRental, error)
	GetManagedPreRentals(userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]rental_model.PreRental, error)
	UpdatePreRentalState(id int64, payload *dto.UpdatePreRental) (newRentalID int64, err error)
	GetRental(id int64) (rental_model.RentalModel, error)
	GetRentalByIds(userId uuid.UUID, ids []int64, fields []string) ([]rental_model.RentalModel, error)
	UpdateRental(data *dto.UpdateRental, id int64) error
	// PrepareRentalContract(id int64, data *dto.PrepareRentalContract) (*model.RentalContractModel, error)
	GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]rental_model.RentalModel, error)
	GetMyRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]rental_model.RentalModel, error)
	CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error)
	CheckPreRentalVisibility(id int64, userId uuid.UUID) (bool, error)
	CreatePreRentalKey(prerental *rental_model.PreRental) (string, error)
	VerifyPreRentalKey(prerental *rental_model.PreRental, key string) error

	CreateContract(data *dto.CreateContract) (*rental_model.ContractModel, error)
	GetRentalContractsOfUser(userId uuid.UUID, query *dto.GetRentalContracts) ([]rental_model.ContractModel, error)
	GetRentalContract(id int64) (*rental_model.ContractModel, error)
	PingRentalContract(id int64) (any, error)
	GetContract(id int64) (*rental_model.ContractModel, error)
	UpdateContract(data *dto.UpdateContract) error
	UpdateContractContent(data *dto.UpdateContractContent) error

	CreateRentalPayment(data *dto.CreateRentalPayment) (rental_model.RentalPayment, error)
	GetRentalPayment(id int64) (rental_model.RentalPayment, error)
	GetPaymentsOfRental(id int64) ([]rental_model.RentalPayment, error)
	GetManagedRentalPayments(uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) ([]rental_model.RentalPayment, error)
	UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error

	PreCreateRentalComplaint(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error
	CreateRentalComplaint(data *dto.CreateRentalComplaint) (rental_model.RentalComplaint, error)
	GetRentalComplaint(id int64) (rental_model.RentalComplaint, error)
	GetRentalComplaintsByRentalId(rid int64, limit, offset int32) ([]rental_model.RentalComplaint, error)
	PreCreateRentalComplaintReply(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error
	CreateRentalComplaintReply(data *dto.CreateRentalComplaintReply) (rental_model.RentalComplaintReply, error)
	GetRentalComplaintsOfUser(userId uuid.UUID, query dto.GetRentalComplaintsOfUserQuery) ([]rental_model.RentalComplaint, error)
	GetRentalComplaintReplies(id int64, limit, offset int32) ([]rental_model.RentalComplaintReply, error)
	UpdateRentalComplaintStatus(data *dto.UpdateRentalComplaintStatus) error

	NotifyCreatePreRental(
		r *rental_model.RentalModel,
		secret string,
	) error
	NotifyUpdatePreRental(
		preRental *rental_model.PreRental,
		rental *rental_model.RentalModel,
		updateData *dto.UpdatePreRental,
	) error
	NotifyUpdatePayments(
		r *rental_model.RentalModel,
		rp *rental_model.RentalPayment, // old rental payment data before update
		u *dto.UpdateRentalPayment,
	) error
	NotifyCreateRentalPayment(
		r *rental_model.RentalModel,
		rp *rental_model.RentalPayment,
	) error
	NotifyCreateContract(
		c *rental_model.ContractModel,
		r *rental_model.RentalModel,
	) error
	NotifyUpdateContract(
		c *rental_model.ContractModel,
		r *rental_model.RentalModel,
		side string,
	) error
	NotifyCreateRentalComplaint(
		c *rental_model.RentalComplaint,
		r *rental_model.RentalModel,
	) error
	NotifyCreateComplaintReply(
		c *rental_model.RentalComplaint,
		cr *rental_model.RentalComplaintReply,
		r *rental_model.RentalModel,
	) error
	NotifyUpdateComplaintStatus(
		c *rental_model.RentalComplaint,
		r *rental_model.RentalModel,
		status database.RENTALCOMPLAINTSTATUS,
		updatedBy uuid.UUID,
	) error
}

type service struct {
	domainRepo repos.DomainRepo

	mService misc_service.Service

	cronEntries []cron.EntryID

	s3Client        s3.S3Client
	imageBucketName string

	asynctaskDistributor asynctask.Distributor

	feSite string
	secret string
}

func NewService(
	domainRepo repos.DomainRepo,
	mService misc_service.Service,
	c *cron.Cron,
	s3Client s3.S3Client, imageBucketName string,
	feSite, secret string,
	asynctaskDistributor asynctask.Distributor,
) Service {
	res := &service{
		domainRepo:           domainRepo,
		cronEntries:          []cron.EntryID{},
		mService:             mService,
		s3Client:             s3Client,
		imageBucketName:      imageBucketName,
		feSite:               feSite,
		secret:               secret,
		asynctaskDistributor: asynctaskDistributor,
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

func (s *service) CreatePreRentalKey(prerental *rental_model.PreRental) (string, error) {
	return utils.CreatePreRentalKey(s.secret, prerental)
}

func (s *service) VerifyPreRentalKey(prerental *rental_model.PreRental, key string) error {
	return utils.VerifyPreRentalKey(prerental, key, s.secret)
}
