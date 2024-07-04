package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreatePreRental(ctx context.Context, data *dto.CreateRental) (model.PreRental, error)
	CreateRental(ctx context.Context, data *dto.CreateRental) (model.RentalModel, error)
	GetPreRental(ctx context.Context, id int64) (model.PreRental, error)
	GetPreRentalsToTenant(ctx context.Context, userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]model.PreRental, error)
	GetManagedPreRentals(ctx context.Context, userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]model.PreRental, error)
	MovePreRentalToRental(ctx context.Context, id int64) (model.RentalModel, error)
	RemovePreRental(ctx context.Context, id int64) error
	GetRental(ctx context.Context, id int64) (model.RentalModel, error)
	GetRentalsByIds(ctx context.Context, ids []int64, fields []string) ([]model.RentalModel, error)
	GetRentalSide(ctx context.Context, id int64, userId uuid.UUID) (string, error)
	GetManagedRentals(ctx context.Context, userId uuid.UUID, query *dto.GetRentalsQuery) ([]int64, error)
	GetMyRentals(ctx context.Context, userId uuid.UUID, query *dto.GetRentalsQuery) ([]int64, error)
	UpdateRental(ctx context.Context, data *dto.UpdateRental, id int64) error
	// UpdateRentalContract(ctx context.Context, data *dto.UpdateRentalContract, id int64) error
	CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)
	FilterVisibleRentals(ctx context.Context, userId uuid.UUID, ids []int64) ([]int64, error)
	CheckPreRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)

	CreateContract(ctx context.Context, data *dto.CreateContract) (*model.ContractModel, error)
	GetRentalContractsOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRentalContracts) ([]int64, error)
	GetContractsByIds(ctx context.Context, ids []int64, fields []string) ([]model.ContractModel, error)
	GetContractByID(ctx context.Context, id int64) (*model.ContractModel, error)
	GetContractByRentalID(ctx context.Context, id int64) (*model.ContractModel, error)
	PingRentalContract(ctx context.Context, id int64) (any, error)
	UpdateContract(ctx context.Context, data *dto.UpdateContract) error
	UpdateContractContent(ctx context.Context, data *dto.UpdateContractContent) error

	CreateRentalPayment(ctx context.Context, data *dto.CreateRentalPayment) (model.RentalPayment, error)
	GetRentalPayment(ctx context.Context, id int64) (model.RentalPayment, error)
	GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]model.RentalPayment, error)
	GetManagedRentalPayments(ctx context.Context, uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) ([]model.RentalPayment, error)
	UpdateRentalPayment(ctx context.Context, data *dto.UpdateRentalPayment) error
	PlanRentalPayments(ctx context.Context) ([]int64, error)
	PlanRentalPayment(ctx context.Context, rentalId int64) ([]int64, error)
	UpdateFinePayments(ctx context.Context) error
	UpdateFinePaymentsOfRental(ctx context.Context, rentalId int64) error

	CreateRentalComplaint(ctx context.Context, data *dto.CreateRentalComplaint) (model.RentalComplaint, error)
	GetRentalComplaint(ctx context.Context, id int64) (model.RentalComplaint, error)
	GetRentalComplaintsOfUser(ctx context.Context, userId uuid.UUID, query dto.GetRentalComplaintsOfUserQuery) ([]model.RentalComplaint, error)
	GetRentalComplaintsByRentalId(ctx context.Context, rid int64, limit, offset int32) ([]model.RentalComplaint, error)
	CreateRentalComplaintReply(ctx context.Context, data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error)
	GetRentalComplaintReplies(ctx context.Context, rid int64, limit, offset int32) ([]model.RentalComplaintReply, error)
	UpdateRentalComplaint(ctx context.Context, data *dto.UpdateRentalComplaint) error
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}
