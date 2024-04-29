package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateRental(ctx context.Context, data *dto.CreateRental) (model.RentalModel, error)
	GetRental(ctx context.Context, id int64) (model.RentalModel, error)
	GetRentalSide(ctx context.Context, id int64, userId uuid.UUID) (string, error)
	UpdateRental(ctx context.Context, data *dto.UpdateRental, id int64) error
	// UpdateRentalContract(ctx context.Context, data *dto.UpdateRentalContract, id int64) error
	CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)

	CreateContract(ctx context.Context, data *dto.CreateContract) (*model.ContractModel, error)
	GetContractByID(ctx context.Context, id int64) (*model.ContractModel, error)
	GetContractByRentalID(ctx context.Context, id int64) (*model.ContractModel, error)
	PingRentalContract(ctx context.Context, id int64) (any, error)
	UpdateContract(ctx context.Context, data *dto.UpdateContract) error
	UpdateContractContent(ctx context.Context, data *dto.UpdateContractContent) error

	CreateRentalPayment(ctx context.Context, data *dto.CreateRentalPayment) (model.RentalPayment, error)
	GetRentalPayment(ctx context.Context, id int64) (model.RentalPayment, error)
	GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]model.RentalPayment, error)
	UpdateRentalPayment(ctx context.Context, data *dto.UpdateRentalPayment) error
	PlanRentalPayments(ctx context.Context) ([]int64, error)
	PlanRentalPayment(ctx context.Context, rentalId int64) ([]int64, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) CreateRental(ctx context.Context, data *dto.CreateRental) (model.RentalModel, error) {
	prdb, err := r.dao.CreateRental(ctx, data.ToCreateRentalDB())
	if err != nil {
		return model.RentalModel{}, err
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
		for _, items := range data.Policies {
			policydb, err := r.dao.CreateRentalPolicy(ctx, items.ToCreateRentalPolicyDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Policies = append(prm.Policies, model.RentalPolicy(policydb))
		}
		return nil
	}()
	if err != nil {
		_err := r.dao.DeleteRental(ctx, prdb.ID)
		return model.RentalModel{}, errors.Join(err, _err)
	}

	return prm, nil
}

func (r *repo) GetRental(ctx context.Context, id int64) (model.RentalModel, error) {
	prdb, err := r.dao.GetRental(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	prm := model.ToRentalModel(&prdb)

	coapdb, err := r.dao.GetRentalCoapsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range coapdb {
		prm.Coaps = append(prm.Coaps, model.ToRentalCoapModel(&item))
	}

	minordb, err := r.dao.GetRentalMinorsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range minordb {
		prm.Minors = append(prm.Minors, model.ToRentalMinor(&item))
	}

	petdb, err := r.dao.GetRentalPetsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range petdb {
		prm.Pets = append(prm.Pets, model.ToRentalPet(&item))
	}

	servicedb, err := r.dao.GetRentalServicesByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range servicedb {
		prm.Services = append(prm.Services, model.ToRentalService(&item))
	}

	return prm, nil
}

func (r *repo) GetRentalSide(ctx context.Context, id int64, userId uuid.UUID) (string, error) {
	return r.dao.GetRentalSide(ctx, database.GetRentalSideParams{
		ID:     id,
		UserID: userId,
	})
}

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

func (r *repo) CreateRentalPayment(ctx context.Context, data *dto.CreateRentalPayment) (model.RentalPayment, error) {
	res, err := r.dao.CreateRentalPayment(ctx, data.ToCreateRentalPaymentDB())
	if err != nil {
		return model.RentalPayment{}, err
	}
	return model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayment(ctx context.Context, id int64) (model.RentalPayment, error) {
	res, err := r.dao.GetRentalPayment(ctx, id)
	if err != nil {
		return model.RentalPayment{}, err
	}
	return model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayments(ctx context.Context, ids []int64) ([]model.RentalPayment, error) {
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select("id", "code", "rental_id", "created_at", "updated_at", "expiry_date", "payment_date", "updated_by", "status", "amount", "note")
	ib.From("rental_payments")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	query, args := ib.Build()
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []model.RentalPayment
		i     database.RentalPayment
	)
	for rows.Next() {
		if err := rows.Scan(&i.ID, &i.Code, &i.RentalID, &i.CreatedAt, &i.UpdatedAt, &i.ExpiryDate, &i.PaymentDate, &i.UpdatedBy, &i.Status, &i.Amount, &i.Note); err != nil {
			return nil, err
		}
		items = append(items, model.ToRentalPaymentModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repo) GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]model.RentalPayment, error) {
	res, err := r.dao.GetPaymentsOfRental(ctx, rentalID)
	if err != nil {
		return nil, err
	}

	var (
		rms []model.RentalPayment
		rm  model.RentalPayment
	)
	for i := range res {
		rm = model.ToRentalPaymentModel(&res[i])
		rms = append(rms, rm)
	}
	return rms, nil
}

func (r *repo) UpdateRentalPayment(ctx context.Context, data *dto.UpdateRentalPayment) error {
	return r.dao.UpdateRentalPayment(ctx, data.ToUpdateRentalPaymentDB())
}

func (r *repo) PlanRentalPayments(ctx context.Context) ([]int64, error) {
	return r.dao.PlanRentalPayments(ctx)
}

func (r *repo) PlanRentalPayment(ctx context.Context, rentalId int64) ([]int64, error) {
	return r.dao.PlanRentalPayment(ctx, rentalId)
}
