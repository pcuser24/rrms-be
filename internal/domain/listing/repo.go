package listing

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(ctx context.Context, data *dto.UpdateListing) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	AddListingPolicies(ctx context.Context, lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error)
	AddListingUnits(ctx context.Context, lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error)
	DeleteListingPolicies(ctx context.Context, lid uuid.UUID, ids []int64) error
	DeleteListingUnits(ctx context.Context, lid uuid.UUID, ids []uuid.UUID) error
	CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
}

type repo struct {
	dao db.DAO
}

func NewRepo(d db.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d db.DAO) (interface{}, error) {
		var lm *model.ListingModel
		res, err := d.CreateListing(ctx, *data.ToCreateListingDB())
		if err != nil {
			return nil, err
		}
		lm = model.ToListingModel(&res)

		lm.Policies, err = r.AddListingPolicies(ctx, res.ID, data.Policies)
		if err != nil {
			return nil, err
		}

		lm.Units, err = r.AddListingUnits(ctx, res.ID, data.Units)
		if err != nil {
			return nil, err
		}

		return lm, nil
	})
	if err != nil {
		return nil, err
	}

	l := res.(*model.ListingModel)

	return l, nil
}

func (r *repo) AddListingPolicies(ctx context.Context, lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error) {
	var res []model.ListingPolicyModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("listing_policy")
	ib.Cols("listing_id", "policy_id", "note")
	for _, i := range items {
		ib.Values(lid, i.PolicyID, types.StrN(i.Note))
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()

	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.ListingPolicyModel, error) {
		defer rows.Close()
		var items []model.ListingPolicyModel
		for rows.Next() {
			var i db.ListingPolicy
			if err := rows.Scan(
				&i.ListingID,
				&i.PolicyID,
				&i.Note,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToListingPolicyModel(&i))
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return items, nil
	}()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) AddListingUnits(ctx context.Context, lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error) {
	var res []model.ListingUnitModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("listing_unit")
	ib.Cols("listing_id", "unit_id")
	for _, i := range items {
		ib.Values(lid, i.UnitID)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()

	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.ListingUnitModel, error) {
		defer rows.Close()
		var items []model.ListingUnitModel
		for rows.Next() {
			var i db.ListingUnit
			if err := rows.Scan(
				&i.ListingID,
				&i.UnitID,
			); err != nil {
				return nil, err
			}
			items = append(items, model.ListingUnitModel(i))
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return items, nil
	}()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) GetListingByID(ctx context.Context, id uuid.UUID) (*model.ListingModel, error) {
	resDB, err := r.dao.GetListingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := model.ToListingModel(&resDB)

	p, err := r.dao.GetListingPolicies(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, i := range p {
		res.Policies = append(res.Policies, *model.ToListingPolicyModel(&i))
	}

	u, err := r.dao.GetListingUnits(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, i := range u {
		res.Units = append(res.Units, model.ListingUnitModel(i))
	}

	return res, nil
}

func (r *repo) UpdateListing(ctx context.Context, data *dto.UpdateListing) error {
	return r.dao.UpdateListing(ctx, *data.ToUpdateListingDB())
}

func (r *repo) DeleteListing(ctx context.Context, lid uuid.UUID) error {
	return r.dao.DeleteListing(ctx, lid)
}

func (r *repo) bulkDelete(ctx context.Context, uid uuid.UUID, ids []interface{}, table_name, info_id_field string) error {
	if len(ids) == 0 {
		return nil
	}

	ib := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	ib.DeleteFrom(table_name)
	ib.Where(
		ib.Equal("listing_id", uid),
		ib.In(info_id_field, ids...),
	)
	sql, args := ib.Build()
	_, err := r.dao.ExecContext(ctx, sql, args...)
	return err
}

func (r *repo) DeleteListingPolicies(ctx context.Context, lid uuid.UUID, ids []int64) error {
	ids_i := make([]interface{}, len(ids))
	for i, v := range ids {
		ids_i[i] = v
	}
	return r.bulkDelete(ctx, lid, ids_i, "listing_policy", "policy_id")
}

func (r *repo) DeleteListingUnits(ctx context.Context, lid uuid.UUID, ids []uuid.UUID) error {
	ids_i := make([]interface{}, len(ids))
	for i, v := range ids {
		ids_i[i] = v
	}
	return r.bulkDelete(ctx, lid, ids_i, "listing_unit", "unit_id")
}

func (r *repo) CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckListingOwnership(ctx, db.CheckListingOwnershipParams{
		ID:        lid,
		CreatorID: uid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckValidUnitForListing(ctx, db.CheckValidUnitForListingParams{
		ID:   uid,
		ID_2: lid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}
