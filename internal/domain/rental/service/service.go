package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Service interface {
	CreateRental(data *dto.CreateRental, creatorID uuid.UUID) (model.RentalModel, error)
	GetRental(id int64) (model.RentalModel, error)
	UpdateRental(data *dto.UpdateRental, id int64) error
	// PrepareRentalContract(id int64, data *dto.PrepareRentalContract) (*model.RentalContractModel, error)
	GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error)
	CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error)

	CreateContract(data *dto.CreateContract) (*model.ContractModel, error)
	GetRentalContract(id int64) (*model.ContractModel, error)
	PingRentalContract(id int64) (any, error)
	GetContract(id int64) (*model.ContractModel, error)
	UpdateContract(data *dto.UpdateContract) error
	UpdateContractContent(data *dto.UpdateContractContent) error

	CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error)
	GetRentalPayment(id int64) (model.RentalPayment, error)
	GetPaymentsOfRental(id int64) ([]model.RentalPayment, error)
	UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error

	CreateRentalComplaint(data *dto.CreateRentalComplaint) (model.RentalComplaint, error)
	GetRentalComplaint(id int64) (model.RentalComplaint, error)
	GetRentalComplaintsByRentalId(rid int64) ([]model.RentalComplaint, error)
	CreateRentalComplaintReply(data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error)
	GetRentalComplaintReplies(id int64) ([]model.RentalComplaintReply, error)
	UpdateRentalComplaint(data *dto.UpdateRentalComplaint) error
}

type service struct {
	authRepo auth_repo.Repo
	aRepo    application_repo.Repo
	lRepo    listing_repo.Repo
	pRepo    property_repo.Repo
	rRepo    repo.Repo
	uRepo    unit_repo.Repo

	cronEntries []cron.EntryID
}

func NewService(
	rRepo repo.Repo, authRepo auth_repo.Repo, aRepo application_repo.Repo, lRepo listing_repo.Repo, pRepo property_repo.Repo, uRepo unit_repo.Repo,
	c *cron.Cron,
) Service {
	res := &service{
		authRepo:    authRepo,
		aRepo:       aRepo,
		rRepo:       rRepo,
		lRepo:       lRepo,
		pRepo:       pRepo,
		uRepo:       uRepo,
		cronEntries: []cron.EntryID{},
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
		s.rRepo.PlanRentalPayments(context.Background())
	})
	if err != nil {
		return nil, err
	}
	s.cronEntries = append(s.cronEntries, entryID)

	return s.cronEntries, nil
}
