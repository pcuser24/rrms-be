package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateProperty(ctx context.Context, data *dto.CreateProperty) (*model.PropertyModel, error)
	GetPropertyManagers(ctx context.Context, id uuid.UUID) ([]model.PropertyManagerModel, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error)
	IsPublic(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	AddPropertyFeatures(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error)
	AddPropertyMedia(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error)
	AddPropertyTag(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error)
	AddPropertyManagers(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyManager) ([]model.PropertyManagerModel, error)
	GetAllFeatures(ctx context.Context) ([]model.PFeature, error)
	DeletePropertyFeatures(ctx context.Context, puid uuid.UUID, fid []int64) error
	DeletePropertyMedia(ctx context.Context, puid uuid.UUID, mid []int64) error
	DeletePropertyTags(ctx context.Context, puid uuid.UUID, tid []int64) error
	DeletePropertyManager(ctx context.Context, puid uuid.UUID, mid uuid.UUID) error
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

		var pm *model.PropertyModel

		res, err := d.CreateProperty(ctx, *data.ToCreatePropertyDB())
		if err != nil {
			return nil, err
		}
		pm = model.ToPropertyModel(&res)

		pm.Managers, err = r.AddPropertyManagers(ctx, res.ID, data.Managers)
		if err != nil {
			return nil, err
		}

		pm.Features, err = r.AddPropertyFeatures(ctx, res.ID, data.Features)
		if err != nil {
			return nil, err
		}

		pm.Media, err = r.AddPropertyMedia(ctx, res.ID, data.Media)
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
	p := res.(*model.PropertyModel)

	return p, nil
}

func (r *repo) GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error) {
	p, err := r.dao.GetPropertyById(ctx, id)
	if err != nil {
		return nil, err
	}

	pm := model.ToPropertyModel(&p)

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

	m, err := r.dao.GetPropertyMedia(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, mdb := range m {
		pm.Media = append(pm.Media, *model.ToPropertyMediaModel(&mdb))
	}

	return pm, nil
}

func (r *repo) GetPropertyManagers(ctx context.Context, id uuid.UUID) ([]model.PropertyManagerModel, error) {
	res, err := r.dao.GetPropertyManagers(ctx, id)
	if err != nil {
		return nil, err
	}
	var items []model.PropertyManagerModel
	for _, i := range res {
		items = append(items, model.PropertyManagerModel(i))
	}
	return items, err
}

func (r *repo) IsPublic(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.dao.IsPropertyPublic(ctx, id)
}

func (r *repo) UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error {
	return r.dao.UpdateProperty(ctx, *data.ToUpdatePropertyDB())
}

func (r *repo) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteProperty(ctx, id)
}

func (r *repo) AddPropertyFeatures(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error) {
	var res []model.PropertyFeatureModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_features")
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

func (r *repo) AddPropertyManagers(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyManager) ([]model.PropertyManagerModel, error) {
	var res []model.PropertyManagerModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_managers")
	ib.Cols("property_id", "manager_id", "role")
	for _, m := range items {
		ib.Values(id, m.ManagerID, m.Role)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()
	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.PropertyManagerModel, error) {
		defer rows.Close()
		var items []model.PropertyManagerModel
		for rows.Next() {
			var i db.PropertyManager
			if err := rows.Scan(
				&i.PropertyID,
				&i.ManagerID,
				&i.Role,
			); err != nil {
				return nil, err
			}
			items = append(items, model.PropertyManagerModel(i))
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
	ib.InsertInto("property_tags")
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

func (r *repo) AddPropertyMedia(ctx context.Context, id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error) {
	var res []model.PropertyMediaModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("property_media")
	ib.Cols("property_id", "url", "type", "description")
	for _, media := range items {
		ib.Values(id, media.Url, media.Type, types.StrN((media.Description)))
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
			var i db.PropertyMedia
			if err := rows.Scan(
				&i.ID,
				&i.PropertyID,
				&i.Url,
				&i.Type,
				&i.Description,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToPropertyMediaModel(&i))
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

func (r *repo) DeletePropertyFeatures(ctx context.Context, puid uuid.UUID, fid []int64) error {
	return r.bulkDelete(ctx, puid, fid, "property_features", "feature_id")
}

func (r *repo) DeletePropertyTags(ctx context.Context, puid uuid.UUID, tid []int64) error {
	return r.bulkDelete(ctx, puid, tid, "property_tags", "id")
}

func (r *repo) DeletePropertyMedia(ctx context.Context, puid uuid.UUID, mid []int64) error {
	return r.bulkDelete(ctx, puid, mid, "property_media", "id")
}

func (r *repo) DeletePropertyManager(ctx context.Context, puid uuid.UUID, mid uuid.UUID) error {
	return r.dao.DeletePropertyManager(ctx, db.DeletePropertyManagerParams{
		PropertyID: puid,
		ManagerID:  mid,
	})
}

func SearchPropertyBuilder(
	searchFields []string, query *dto.SearchPropertyQuery,
	connectID, connectCreator string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM properties WHERE "
	var searchQueries []string
	var args []interface{}

	if query.PIsPublic != nil {
		searchQueries = append(searchQueries, "properties.is_public = $?")
		args = append(args, *query.PIsPublic)
	}
	if query.PName != nil {
		searchQueries = append(searchQueries, "properties.name ILIKE $?")
		args = append(args, "%"+(*query.PName)+"%")
	}
	if query.PCreatorID != nil {
		searchQueries = append(searchQueries, "properties.creator_id = $?")
		args = append(args, *query.PCreatorID)
	}
	if query.PBuilding != nil {
		searchQueries = append(searchQueries, "properties.building ILIKE $?")
		args = append(args, "%"+(*query.PBuilding)+"%")
	}
	if query.PProject != nil {
		searchQueries = append(searchQueries, "properties.project ILIKE $?")
		args = append(args, "%"+(*query.PProject)+"%")
	}
	if query.PFullAddress != nil {
		searchQueries = append(searchQueries, "properties.full_address ILIKE $?")
		args = append(args, "%"+(*query.PFullAddress)+"%")
	}
	if query.PCity != nil {
		searchQueries = append(searchQueries, "properties.city = $?")
		args = append(args, *query.PCity)
	}
	if query.PDistrict != nil {
		searchQueries = append(searchQueries, "properties.district = $?")
		args = append(args, *query.PDistrict)
	}
	if query.PWard != nil {
		searchQueries = append(searchQueries, "properties.ward = $?")
		args = append(args, *query.PWard)
	}
	if query.PMinArea != nil {
		searchQueries = append(searchQueries, "properties.area >= $?")
		args = append(args, *query.PMinArea)
	}
	if query.PMaxArea != nil {
		searchQueries = append(searchQueries, "properties.area <= $?")
		args = append(args, *query.PMaxArea)
	}
	if query.PNumberOfFloors != nil {
		searchQueries = append(searchQueries, "properties.number_of_floors = $?")
		args = append(args, *query.PNumberOfFloors)
	}
	if query.PYearBuilt != nil {
		searchQueries = append(searchQueries, "properties.year_built = $?")
		args = append(args, *query.PYearBuilt)
	}
	if query.POrientation != nil {
		searchQueries = append(searchQueries, "properties.orientation = $?")
		args = append(args, *query.POrientation)
	}
	if query.PMinFacade != nil {
		searchQueries = append(searchQueries, "properties.facade >= $?")
		args = append(args, *query.PMinFacade)
	}
	if len(query.PTypes) > 0 {
		searchQueries = append(searchQueries, "properties.type IN ($?)")
		args = append(args, sqlbuilder.List(query.PTypes))
	}
	if query.PMinCreatedAt != nil {
		searchQueries = append(searchQueries, "properties.created_at >= $?")
		args = append(args, *query.PMinCreatedAt)
	}
	if query.PMaxCreatedAt != nil {
		searchQueries = append(searchQueries, "properties.created_at <= $?")
		args = append(args, *query.PMaxCreatedAt)
	}
	if query.PMinUpdatedAt != nil {
		searchQueries = append(searchQueries, "properties.updated_at >= $?")
		args = append(args, *query.PMinUpdatedAt)
	}
	if query.PMaxUpdatedAt != nil {
		searchQueries = append(searchQueries, "properties.updated_at <= $?")
		args = append(args, *query.PMaxUpdatedAt)
	}
	if len(query.PFeatures) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_features WHERE property_id = properties.id AND feature_id IN ($?))")
		args = append(args, sqlbuilder.List(query.PFeatures))
	}
	if len(query.PTags) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_tags WHERE property_id = properties.id AND tag IN ($?))")
		args = append(args, sqlbuilder.List(query.PTags))
	}

	if len(searchQueries) == 0 {
		return "", []interface{}{}
	}
	if len(connectID) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.id = %v", connectID))
	}
	if len(connectCreator) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.creator_id = %v", connectCreator))
	}
	searchQuery += strings.Join(searchQueries, " AND \n")
	return searchQuery, args
}
