package repo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

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

func (r *repo) CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error) {
	return r.dao.CheckRentalVisibility(ctx, database.CheckRentalVisibilityParams{
		ID: id,
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
	})
}
