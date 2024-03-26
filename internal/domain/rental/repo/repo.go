package repo

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Repo interface {
	CreatePreRental(ctx context.Context, data *dto.CreatePreRental) (*model.PrerentalModel, error)
	GetPreRental(ctx context.Context, id int64) (*model.PrerentalModel, error)
	GetPreRentalContract(ctx context.Context, id int64) (*model.PreRentalContractModel, error)
	UpdatePreRental(ctx context.Context, data *dto.UpdatePreRental, id int64) error
	UpdatePreRentalContract(ctx context.Context, data *dto.UpdatePreRentalContract, id int64) error
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) CreatePreRental(ctx context.Context, data *dto.CreatePreRental) (*model.PrerentalModel, error) {
	prdb, err := r.dao.CreatePreRental(ctx, data.ToCreatePreRentalDB())
	if err != nil {
		return nil, err
	}
	prm := model.ToPreRentalModel(&prdb)

	for _, items := range data.Coaps {
		coapdb, err := r.dao.CreatePreRentalCoap(ctx, items.ToCreatePreRentalCoapDB(prdb.ID))
		if err != nil {
			_ = r.dao.DeletePreRental(ctx, prdb.ID)
			return nil, err
		}
		prm.Coaps = append(prm.Coaps, *model.ToPreRentalCoapModel(&coapdb))
	}

	return prm, nil
}

func (r *repo) GetPreRental(ctx context.Context, id int64) (*model.PrerentalModel, error) {
	prdb, err := r.dao.GetPreRental(ctx, id)
	if err != nil {
		return nil, err
	}
	prm := model.ToPreRentalModel(&prdb)

	coapdb, err := r.dao.GetPreRentalCoapByPreRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, item := range coapdb {
		prm.Coaps = append(prm.Coaps, *model.ToPreRentalCoapModel(&item))
	}

	return prm, nil
}

func (r *repo) GetPreRentalContract(ctx context.Context, id int64) (*model.PreRentalContractModel, error) {
	prcdb, err := r.dao.GetPreRentalContract(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.PreRentalContractModel{
		ID:                   id,
		ContractType:         prcdb.ContractType.CONTRACTTYPE,
		ContractContent:      types.PNStr(prcdb.ContractContent),
		ContractLastUpdateAt: prcdb.ContractLastUpdateAt.Time,
		ContractLastUpdateBy: prcdb.ContractLastUpdateBy.Bytes,
	}, nil
}

func (r *repo) UpdatePreRental(ctx context.Context, data *dto.UpdatePreRental, id int64) error {
	return r.dao.UpdatePreRental(ctx, data.ToUpdatePreRentalDB(id))
}

func (r *repo) UpdatePreRentalContract(ctx context.Context, data *dto.UpdatePreRentalContract, id int64) error {
	return r.dao.UpdatePreRentalContract(ctx, data.ToUpdatePreRentalContractDB(id))
}
