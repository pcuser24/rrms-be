package property

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateProperty(ctx context.Context, data *dto.CreateProperty) (*model.PropertyModel, error)
	CheckOwnership(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error)
	UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	AddPropertyAmenities(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error)
	AddPropertyFeatures(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error)
	AddPropertyMedium(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error)
	AddPropertyTag(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error)
	GetAllAmenities(ctx context.Context) ([]model.PAmenity, error)
	GetAllFeatures(ctx context.Context) ([]model.PFeature, error)
	DeletePropertyAmenities(ctx context.Context, puid uuid.UUID, aid []int64) error
	DeletePropertyFeatures(ctx context.Context, puid uuid.UUID, fid []int64) error
	DeletePropertyMedium(ctx context.Context, puid uuid.UUID, mid []int64) error
	DeletePropertyTags(ctx context.Context, puid uuid.UUID, tid []int64) error
	DeleteAllPropertyAmenities(ctx context.Context, puid uuid.UUID) error
	DeleteAllPropertyFeatures(ctx context.Context, puid uuid.UUID) error
	DeleteAllPropertyMedium(ctx context.Context, puid uuid.UUID) error
	DeleteAllPropertyTags(ctx context.Context, puid uuid.UUID) error
}

type repo struct {
	dao db.DAO
}

func NewRepo(d db.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateProperty(ctx context.Context, data *dto.CreateProperty) (*model.PropertyModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d db.DAO) (interface{}, error) {

		var pm model.PropertyModel

		// create property
		res, err := d.CreateProperty(ctx, *data.ToCreatePropertyDB())
		if err != nil {
			return nil, err
		}
		pm = *model.ToPropertyModel(&res)

		pm.Amenities, err = r.AddPropertyAmenities(ctx, res.ID, data.Amenities)
		if err != nil {
			return nil, err
		}

		pm.Features, err = r.AddPropertyFeatures(ctx, res.ID, data.Features)
		if err != nil {
			return nil, err
		}

		pm.Medium, err = r.AddPropertyMedium(ctx, res.ID, data.Medium)
		if err != nil {
			return nil, err
		}

		pm.Tags, err = r.AddPropertyTag(ctx, res.ID, data.Tags)
		if err != nil {
			return nil, err
		}

		return pm, nil
	})
	if err != nil {
		return nil, err
	}
	p := res.(model.PropertyModel)

	return &p, nil
}

func (r *repo) GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error) {
	p, err := r.dao.GetPropertyById(ctx, id)
	if err != nil {
		return nil, err
	}

	pm := model.ToPropertyModel(&p)

	a, err := r.dao.GetPropertyAmenities(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, adb := range a {
		pm.Amenities = append(pm.Amenities, *model.ToPropertyAmenityModel(&adb))
	}

	f, err := r.dao.GetPropertyFeatures(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, fdb := range f {
		pm.Features = append(pm.Features, *model.ToPropertyFeatureModel(&fdb))
	}

	t, err := r.dao.GetPropertyTags(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, tdb := range t {
		pm.Tags = append(pm.Tags, model.PropertyTagModel(tdb))
	}

	m, err := r.dao.GetPropertyMedium(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, mdb := range m {
		pm.Medium = append(pm.Medium, model.PropertyMediaModel(mdb))
	}

	return pm, nil
}

func (r *repo) CheckOwnership(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error) {
	res, err := r.dao.CheckPropertyOwnerShip(ctx, db.CheckPropertyOwnerShipParams{
		ID:      id,
		OwnerID: userId,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error {
	return r.dao.UpdateProperty(ctx, *data.ToUpdatePropertyDB())
}

func (r *repo) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteProperty(ctx, id)
}

func (r *repo) AddPropertyAmenities(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error) {
	var res []model.PropertyAmenityModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_amenity")
	ib.Cols("property_id", "amenity_id", "description")
	for _, amenity := range items {
		ib.Values(id, amenity.AmenityID, types.StrN((amenity.Description)))
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.PropertyAmenityModel, error) {
		defer rows.Close()
		var items []model.PropertyAmenityModel
		for rows.Next() {
			var i db.PropertyAmenity
			if err := rows.Scan(
				&i.PropertyID,
				&i.AmenityID,
				&i.Description,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToPropertyAmenityModel(&i))
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

func (r *repo) AddPropertyFeatures(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error) {
	var res []model.PropertyFeatureModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_feature")
	ib.Cols("property_id", "feature_id", "description")
	for _, ft := range items {
		ib.Values(id, ft.FeatureID, types.StrN((ft.Description)))
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.PropertyFeatureModel, error) {
		defer rows.Close()
		var items []model.PropertyFeatureModel
		for rows.Next() {
			var i db.PropertyFeature
			if err := rows.Scan(
				&i.PropertyID,
				&i.FeatureID,
				&i.Description,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToPropertyFeatureModel(&i))
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

func (r *repo) AddPropertyTag(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error) {
	var res []model.PropertyTagModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_tag")
	ib.Cols("property_id", "tag")
	for _, tag := range items {
		ib.Values(id, tag.Tag)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.PropertyTagModel, error) {
		defer rows.Close()
		var items []model.PropertyTagModel
		for rows.Next() {
			var i db.PropertyTag
			if err := rows.Scan(
				&i.ID,
				&i.PropertyID,
				&i.Tag,
			); err != nil {
				return nil, err
			}
			items = append(items, model.PropertyTagModel(i))
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

func (r *repo) AddPropertyMedium(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error) {
	var res []model.PropertyMediaModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_media")
	ib.Cols("property_id", "url", "type")
	for _, media := range items {
		ib.Values(id, media.Url, media.Type)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.PropertyMediaModel, error) {
		defer rows.Close()
		var items []model.PropertyMediaModel
		for rows.Next() {
			var i db.PropertyMedium
			if err := rows.Scan(
				&i.ID,
				&i.PropertyID,
				&i.Url,
				&i.Type,
			); err != nil {
				return nil, err
			}
			items = append(items, model.PropertyMediaModel(i))
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

func (r *repo) GetAllAmenities(ctx context.Context) ([]model.PAmenity, error) {
	resDb, err := r.dao.GetAllPropertyAmenities(ctx)
	if err != nil {
		return nil, err
	}
	var res []model.PAmenity
	for _, i := range resDb {
		res = append(res, model.PAmenity(i))
	}
	return res, nil
}

func (r *repo) GetAllFeatures(ctx context.Context) ([]model.PFeature, error) {
	resDb, err := r.dao.GetAllPropertyFeatures(ctx)
	if err != nil {
		return nil, err
	}
	var res []model.PFeature
	for _, i := range resDb {
		res = append(res, model.PFeature(i))
	}
	return res, nil
}

func (r *repo) bulkDelete(ctx context.Context, puid uuid.UUID, ids []int64, table_name, info_id_field string) error {
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
		ib.Equal("property_id", puid),
		ib.In(info_id_field, ids_i...),
	)
	sql, args := ib.Build()
	_, err := r.dao.ExecContext(ctx, sql, args...)
	return err
}

func (r *repo) DeletePropertyAmenities(ctx context.Context, puid uuid.UUID, aid []int64) error {
	return r.bulkDelete(ctx, puid, aid, "property_amenity", "amenity_id")
}

func (r *repo) DeletePropertyFeatures(ctx context.Context, puid uuid.UUID, fid []int64) error {
	return r.bulkDelete(ctx, puid, fid, "property_feature", "feature_id")
}

func (r *repo) DeletePropertyTags(ctx context.Context, puid uuid.UUID, tid []int64) error {
	return r.bulkDelete(ctx, puid, tid, "property_tag", "id")
}

func (r *repo) DeletePropertyMedium(ctx context.Context, puid uuid.UUID, mid []int64) error {
	return r.bulkDelete(ctx, puid, mid, "property_media", "id")
}

func (r *repo) DeleteAllPropertyAmenities(ctx context.Context, puid uuid.UUID) error {
	return r.dao.DeleteAllPropertyAmenity(ctx, puid)
}

func (r *repo) DeleteAllPropertyFeatures(ctx context.Context, puid uuid.UUID) error {
	return r.dao.DeleteAllPropertyFeature(ctx, puid)
}

func (r *repo) DeleteAllPropertyTags(ctx context.Context, puid uuid.UUID) error {
	return r.dao.DeleteAllPropertyTag(ctx, puid)
}

func (r *repo) DeleteAllPropertyMedium(ctx context.Context, puid uuid.UUID) error {
	return r.dao.DeleteAllPropertyMedia(ctx, puid)
}
