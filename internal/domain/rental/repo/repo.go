package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateRental(ctx context.Context, data *dto.CreateRental) (*model.RentalModel, error)
	GetRental(ctx context.Context, id int64) (*model.RentalModel, error)
	// GetRentalContract(ctx context.Context, id int64) (*model.RentalContractModel, error)
	UpdateRental(ctx context.Context, data *dto.UpdateRental, id int64) error
	// UpdateRentalContract(ctx context.Context, data *dto.UpdateRentalContract, id int64) error
	CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)

	CreateContract(ctx context.Context, data *dto.CreateContract) (*model.ContractModel, error)
	GetContractByID(ctx context.Context, id int64) (*model.ContractModel, error)
	GetContractByRentalID(ctx context.Context, id int64) (*model.ContractModel, error)
	PingRentalContract(ctx context.Context, id int64) (any, error)
	UpdateContract(ctx context.Context, data *dto.UpdateContract) error
	UpdateContractContent(ctx context.Context, data *dto.UpdateContractContent) error
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) CreateRental(ctx context.Context, data *dto.CreateRental) (*model.RentalModel, error) {
	prdb, err := r.dao.CreateRental(ctx, data.ToCreateRentalDB())
	if err != nil {
		return nil, err
	}
	prm := model.ToRentalModel(&prdb)

	err = func() error {
		for _, items := range data.Coaps {
			coapdb, err := r.dao.CreateRentalCoap(ctx, items.ToCreateRentalCoapDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Coaps = append(prm.Coaps, model.ToRentalCoapModel(&coapdb))
		}
		for _, items := range data.Minors {
			minordb, err := r.dao.CreateRentalMinor(ctx, items.ToCreateRentalMinorDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Minors = append(prm.Minors, model.ToRentalMinor(&minordb))
		}
		for _, items := range data.Pets {
			petdb, err := r.dao.CreateRentalPet(ctx, items.ToCreateRentalPetDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Pets = append(prm.Pets, model.ToRentalPet(&petdb))
		}
		for _, items := range data.Services {
			servicedb, err := r.dao.CreateRentalService(ctx, items.ToCreateRentalServiceDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Services = append(prm.Services, model.ToRentalService(&servicedb))
		}
		return nil
	}()
	if err != nil {
		_err := r.dao.DeleteRental(ctx, prdb.ID)
		return nil, errors.Join(err, _err)
	}

	return prm, nil
}

func (r *repo) GetRental(ctx context.Context, id int64) (*model.RentalModel, error) {
	prdb, err := r.dao.GetRental(ctx, id)
	if err != nil {
		return nil, err
	}
	prm := model.ToRentalModel(&prdb)

	coapdb, err := r.dao.GetRentalCoapsByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, item := range coapdb {
		prm.Coaps = append(prm.Coaps, model.ToRentalCoapModel(&item))
	}

	minordb, err := r.dao.GetRentalMinorsByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, item := range minordb {
		prm.Minors = append(prm.Minors, model.ToRentalMinor(&item))
	}

	petdb, err := r.dao.GetRentalPetsByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, item := range petdb {
		prm.Pets = append(prm.Pets, model.ToRentalPet(&item))
	}

	servicedb, err := r.dao.GetRentalServicesByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, item := range servicedb {
		prm.Services = append(prm.Services, model.ToRentalService(&item))
	}

	return prm, nil
}

// func (r *repo) GetRentalContract(ctx context.Context, id int64) (*model.RentalContractModel, error) {
// 	prcdb, err := r.dao.GetRentalContract(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model.RentalContractModel{
// 		ID:                   id,
// 		ContractType:         prcdb.ContractType.CONTRACTTYPE,
// 		ContractContent:      types.PNStr(prcdb.ContractContent),
// 		ContractLastUpdateAt: prcdb.ContractLastUpdateAt.Time,
// 		ContractLastUpdateBy: prcdb.ContractLastUpdateBy.Bytes,
// 	}, nil
// }

func (r *repo) UpdateRental(ctx context.Context, data *dto.UpdateRental, id int64) error {
	return r.dao.UpdateRental(ctx, data.ToUpdateRentalDB(id))
}

// func (r *repo) UpdateRentalContract(ctx context.Context, data *dto.UpdateRentalContract, id int64) error {
// 	return r.dao.UpdateRentalContract(ctx, data.ToUpdateRentalContractDB(id))
// }

func (r *repo) CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error) {
	return r.dao.CheckRentalVisibility(ctx, database.CheckRentalVisibilityParams{
		ID: id,
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
	})
}

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
