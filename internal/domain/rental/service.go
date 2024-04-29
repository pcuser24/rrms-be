package rental

import (
	"context"
	"errors"
	"time"

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
	CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error)

	CreateContract(data *dto.CreateContract) (*model.ContractModel, error)
	GetRentalContract(id int64) (*model.ContractModel, error)
	PingRentalContract(id int64) (any, error)
	GetContract(id int64) (*model.ContractModel, error)
	UpdateContract(data *dto.UpdateContract) error
	UpdateContractContent(data *dto.UpdateContractContent) error

	SetupCronjob(c *cron.Cron) ([]cron.EntryID, error)

	CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error)
	GetRentalPayment(id int64) (model.RentalPayment, error)
	GetPaymentsOfRental(id int64) ([]model.RentalPayment, error)
	UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error
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
) Service {
	return &service{
		authRepo:    authRepo,
		aRepo:       aRepo,
		rRepo:       rRepo,
		lRepo:       lRepo,
		pRepo:       pRepo,
		uRepo:       uRepo,
		cronEntries: []cron.EntryID{},
	}
}

var ErrInvalidRentalExpired = errors.New("rental expired")

func (s *service) CreateRental(data *dto.CreateRental, userId uuid.UUID) (model.RentalModel, error) {
	expiryDate := data.MoveinDate.AddDate(0, int(data.RentalPeriod), 0)
	today := time.Now().Truncate(24 * time.Hour) // time representing today at 00:00:00
	if expiryDate.Before(today) {
		return model.RentalModel{}, ErrInvalidRentalExpired
	}

	// TODO: validate applicationId, propertyId, unitId
	data.CreatorID = userId
	rental, err := s.rRepo.CreateRental(context.Background(), data)
	if err != nil {
		return model.RentalModel{}, err
	}

	// plan rental payments
	_, err = s.rRepo.PlanRentalPayment(context.Background(), rental.ID)
	if err != nil {
		// TODO: log the error
	}
	// TODO: send notification to user
	return rental, nil
}

func (s *service) GetRental(id int64) (model.RentalModel, error) {
	return s.rRepo.GetRental(context.Background(), id)
}

func (s *service) UpdateRental(data *dto.UpdateRental, id int64) error {
	return s.rRepo.UpdateRental(context.Background(), data, id)
}

func (s *service) CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.rRepo.CheckRentalVisibility(context.Background(), id, userId)
}

func (s *service) CreateContract(data *dto.CreateContract) (*model.ContractModel, error) {
	return s.rRepo.CreateContract(context.Background(), data)
}

func (s *service) GetRentalContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByRentalID(context.Background(), id)
}

func (s *service) PingRentalContract(id int64) (any, error) {
	return s.rRepo.PingRentalContract(context.Background(), id)
}

func (s *service) GetContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByID(context.Background(), id)
}

func (s *service) UpdateContract(data *dto.UpdateContract) error {
	return s.rRepo.UpdateContract(context.Background(), data)
}

func (s *service) UpdateContractContent(data *dto.UpdateContractContent) error {
	return s.rRepo.UpdateContractContent(context.Background(), data)
}

func (s *service) CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error) {
	// TODO: validate rental payment (code)
	return s.rRepo.CreateRentalPayment(context.Background(), data)
}

func (s *service) GetRentalPayment(id int64) (model.RentalPayment, error) {
	return s.rRepo.GetRentalPayment(context.Background(), id)
}

func (s *service) GetPaymentsOfRental(rentalID int64) ([]model.RentalPayment, error) {
	return s.rRepo.GetPaymentsOfRental(context.Background(), rentalID)
}

var ErrInvalidTypeTransition = errors.New("invalid type transition")

func (s *service) UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error {
	rp, err := s.rRepo.GetRentalPayment(context.Background(), id)
	if err != nil {
		return err
	}
	side, err := s.rRepo.GetRentalSide(context.Background(), rp.RentalID, userId)
	if err != nil {
		return err
	}

	var _data dto.UpdateRentalPayment

	if rp.Status != status {
		return ErrInvalidTypeTransition
	}
	switch status {
	case database.RENTALPAYMENTSTATUSPLAN:
		__data := data.(*dto.UpdatePlanRentalPayment)
		if side != "A" {
			return ErrInvalidTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:         id,
			UserID:     userId,
			Status:     __data.Status,
			Amount:     &__data.Amount,
			Discount:   __data.Discount,
			ExpiryDate: __data.ExpiryDate,
		}
	case database.RENTALPAYMENTSTATUSISSUED:
		__data := data.(*dto.UpdateIssuedRentalPayment)
		if side != "B" {
			return ErrInvalidTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:     id,
			UserID: userId,
			Status: __data.Status,
			Note:   __data.Note,
		}
	case database.RENTALPAYMENTSTATUSPENDING:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "B" {
			return ErrInvalidTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:          id,
			UserID:      userId,
			PaymentDate: __data.PaymentDate,
			Status:      database.RENTALPAYMENTSTATUSREQUEST2PAY,
		}
	case database.RENTALPAYMENTSTATUSREQUEST2PAY:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "A" {
			return ErrInvalidTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:          id,
			UserID:      userId,
			PaymentDate: __data.PaymentDate,
			Status:      database.RENTALPAYMENTSTATUSPAID,
		}
	default:
		return ErrInvalidTypeTransition
	}
	return s.rRepo.UpdateRentalPayment(context.Background(), &_data)

}

// SetupCronjob periodically checks for rental payment due date and send reminder to user
func (s *service) SetupCronjob(c *cron.Cron) ([]cron.EntryID, error) {
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
