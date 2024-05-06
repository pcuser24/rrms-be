package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreateContract(ctx context.Context, data *dto.CreateContract) (*model.ContractModel, error) {
	prdb, err := r.dao.CreateContract(ctx, data.ToCreateContractDB())
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) GetContractByRentalID(ctx context.Context, id int64) (*model.ContractModel, error) {
	prdb, err := r.dao.GetContractByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) PingRentalContract(ctx context.Context, id int64) (any, error) {
	res, err := r.dao.PingContractByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	return struct {
		ID        int64                   `json:"id"`
		RentalID  int64                   `json:"rentalId"`
		Status    database.CONTRACTSTATUS `json:"status"`
		UpdatedBy uuid.UUID               `json:"updatedBy"`
		UpdatedAt time.Time               `json:"updatedAt"`
	}{
		ID:        res.ID,
		RentalID:  res.RentalID,
		Status:    res.Status,
		UpdatedBy: res.UpdatedBy,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *repo) GetContractByID(ctx context.Context, id int64) (*model.ContractModel, error) {
	prdb, err := r.dao.GetContractByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) UpdateContract(ctx context.Context, data *dto.UpdateContract) error {
	return r.dao.UpdateContract(ctx, data.ToUpdateContractDB())
}

func (r *repo) UpdateContractContent(ctx context.Context, data *dto.UpdateContractContent) error {
	return r.dao.UpdateContractContent(ctx, data.ToUpdateContractContentDB())
}
