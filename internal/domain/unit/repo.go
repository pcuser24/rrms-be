package unit

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (*model.UnitModel, error)
	GetUnitsOfProperty(ctx context.Context, id uuid.UUID) ([]model.UnitModel, error)
	UpdateUnit(ctx context.Context, data *dto.UpdateUnit) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
	CheckUnitOwnership(ctx context.Context, uid uuid.UUID, userId uuid.UUID) (bool, error)
	AddUnitAmenities(ctx context.Context, uid uuid.UUID, items []dto.CreateUnitAmenity) ([]model.UnitAmenityModel, error)
	AddUnitMedium(ctx context.Context, uid uuid.UUID, items []dto.CreateUnitMedia) ([]model.UnitMediaModel, error)
	GetAllAmenities(ctx context.Context) ([]model.UAmenity, error)
	DeleteUnitAmenities(ctx context.Context, uid uuid.UUID, ids []int64) error
	DeleteUnitMedium(ctx context.Context, uid uuid.UUID, ids []int64) error
}

type repo struct {
	dao db.DAO
}

func NewRepo(d db.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d db.DAO) (interface{}, error) {
		var um *model.UnitModel
		res, err := d.CreateUnit(ctx, *data.ToCreateUnitDB())
		if err != nil {
			return nil, err
		}
		um = model.ToUnitModel(&res)

		um.Amenities, err = r.AddUnitAmenities(ctx, res.ID, data.Amenities)
		if err != nil {
			return nil, err
		}

		um.Medium, err = r.AddUnitMedium(ctx, res.ID, data.Medium)
		if err != nil {
			return nil, err
		}

		return um, nil
	})
	if err != nil {
		return nil, err
	}

	u := res.(*model.UnitModel)
	return u, nil
}

func (r *repo) GetUnitById(ctx context.Context, id uuid.UUID) (*model.UnitModel, error) {
	u, err := r.dao.GetUnitById(ctx, id)
	if err != nil {
		return nil, err
	}

	um := model.ToUnitModel(&u)

	a, err := r.dao.GetUnitAmenities(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, adb := range a {
		um.Amenities = append(um.Amenities, *model.ToUnitAmenityModel(&adb))
	}

	m, err := r.dao.GetUnitMedia(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, mdb := range m {
		um.Medium = append(um.Medium, model.UnitMediaModel(mdb))
	}

	return um, nil
}

func (r *repo) GetUnitsOfProperty(ctx context.Context, id uuid.UUID) ([]model.UnitModel, error) {
	resDb, err := r.dao.GetUnitsOfProperty(ctx, id)
	if err != nil {
		return nil, err
	}

	var res []model.UnitModel
	for _, i := range resDb {
		um := *model.ToUnitModel(&i)
		a, err := r.dao.GetUnitAmenities(ctx, i.ID)
		if err != nil {
			return nil, err
		}
		for _, adb := range a {
			um.Amenities = append(um.Amenities, *model.ToUnitAmenityModel(&adb))
		}
		m, err := r.dao.GetUnitMedia(ctx, i.ID)
		if err != nil {
			return nil, err
		}
		for _, mdb := range m {
			um.Medium = append(um.Medium, model.UnitMediaModel(mdb))
		}
		res = append(res, um)
	}
	return res, nil
}

func (r *repo) UpdateUnit(ctx context.Context, data *dto.UpdateUnit) error {
	return r.dao.UpdateUnit(ctx, *data.ToUpdateUnitDB())
}

func (r *repo) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteUnit(ctx, id)
}

func (r *repo) AddUnitAmenities(ctx context.Context, uid uuid.UUID, items []dto.CreateUnitAmenity) ([]model.UnitAmenityModel, error) {
	var res []model.UnitAmenityModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("unit_amenity")
	ib.Cols("unit_id", "amenity_id", "description")
	for _, i := range items {
		ib.Values(uid, i.AmenityID, types.StrN(i.Description))
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	fmt.Println("amenity sql:", sql)
	fmt.Println("amenity args:", args)

	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.UnitAmenityModel, error) {
		defer rows.Close()
		var items []model.UnitAmenityModel
		for rows.Next() {
			var i db.UnitAmenity
			if err := rows.Scan(
				&i.UnitID,
				&i.AmenityID,
				&i.Description,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToUnitAmenityModel(&i))
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

func (r *repo) AddUnitMedium(ctx context.Context, uid uuid.UUID, items []dto.CreateUnitMedia) ([]model.UnitMediaModel, error) {
	var res []model.UnitMediaModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("unit_media")
	ib.Cols("unit_id", "url", "type")
	for _, media := range items {
		ib.Values(uid, media.Url, media.Type)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	fmt.Println("medium sql:", sql)
	fmt.Println("medium args:", args)
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.UnitMediaModel, error) {
		defer rows.Close()
		var items []model.UnitMediaModel
		for rows.Next() {
			var i db.UnitMedium
			if err := rows.Scan(
				&i.ID,
				&i.UnitID,
				&i.Url,
				&i.Type,
			); err != nil {
				return nil, err
			}
			items = append(items, model.UnitMediaModel(i))
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

	return res, err
}

func (r *repo) bulkDelete(ctx context.Context, uid uuid.UUID, ids []int64, table_name, info_id_field string) error {
	if len(ids) == 0 {
		return nil
	}

	ids_i := make([]interface{}, len(ids))
	for i, v := range ids {
		ids_i[i] = v
	}
	ib := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	ib.DeleteFrom(table_name)
	ib.Where(
		ib.Equal("unit_id", uid),
		ib.In(info_id_field, ids_i...),
	)
	sql, args := ib.Build()
	_, err := r.dao.ExecContext(ctx, sql, args...)
	return err
}

func (r *repo) DeleteUnitAmenities(ctx context.Context, uid uuid.UUID, ids []int64) error {
	return r.bulkDelete(ctx, uid, ids, "unit_amenity", "amenity_id")
}

func (r *repo) DeleteUnitMedium(ctx context.Context, uid uuid.UUID, ids []int64) error {
	return r.bulkDelete(ctx, uid, ids, "unit_media", "id")
}

func (r *repo) GetAllAmenities(ctx context.Context) ([]model.UAmenity, error) {
	resDb, err := r.dao.GetAllUnitAmenities(ctx)
	if err != nil {
		return nil, err
	}
	var res []model.UAmenity
	for _, i := range resDb {
		res = append(res, model.UAmenity(i))
	}
	return res, nil
}

func (r *repo) CheckUnitOwnership(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error) {
	res, err := r.dao.CheckUnitOwnership(ctx, db.CheckUnitOwnershipParams{
		ID:      id,
		OwnerID: userId,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}
