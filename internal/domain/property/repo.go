package property

import (
	"context"
	"database/sql"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	unitModel "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
)

type Repo interface {
	CreateProperty(ctx context.Context, data *dto.CreateProperty) (*model.PropertyModel, error)
	GetPropertyManagers(ctx context.Context, id uuid.UUID) ([]model.PropertyManagerModel, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error)
	// Get properties with custom fields by ids
	GetPropertiesByIds(ctx context.Context, ids []string, fields []string) ([]model.PropertyModel, error)
	GetUnitsOfProperty(ctx context.Context, id uuid.UUID) ([]unitModel.UnitModel, error)
	GetManagedProperties(ctx context.Context, userId uuid.UUID) ([]database.GetManagedPropertiesRow, error)
	SearchPropertyCombination(ctx context.Context, query *dto.SearchPropertyCombinationQuery) (*dto.SearchPropertyCombinationResponse, error)
	IsPublic(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	GetAllFeatures(ctx context.Context) ([]model.PFeature, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateProperty(ctx context.Context, data *dto.CreateProperty) (*model.PropertyModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d database.DAO) (interface{}, error) {

		var pm *model.PropertyModel

		prop, err := d.CreateProperty(ctx, *data.ToCreatePropertyDB())
		if err != nil {
			return nil, err
		}
		pm = model.ToPropertyModel(&prop)

		for _, m := range data.Managers {
			res, err := d.CreatePropertyManager(ctx, database.CreatePropertyManagerParams{
				PropertyID: prop.ID,
				ManagerID:  m.ManagerID,
				Role:       m.Role,
			})
			if err != nil {
				return nil, err
			}
			pm.Managers = append(pm.Managers, model.PropertyManagerModel(res))
		}

		for _, f := range data.Features {
			res, err := d.CreatePropertyFeature(ctx, *f.ToCreatePropertyFeatureDB(prop.ID))
			if err != nil {
				return nil, err
			}
			pm.Features = append(pm.Features, *model.ToPropertyFeatureModel(&res))
		}

		for _, m := range data.Media {
			res, err := d.CreatePropertyMedia(ctx, *m.ToCreatePropertyMediaDB(prop.ID))
			if err != nil {
				return nil, err
			}
			pm.Media = append(pm.Media, *model.ToPropertyMediaModel(&res))
		}

		for _, t := range data.Tags {
			res, err := d.CreatePropertyTag(ctx, database.CreatePropertyTagParams{
				PropertyID: prop.ID,
				Tag:        t.Tag,
			})
			if err != nil {
				return nil, err
			}
			pm.Tags = append(pm.Tags, model.PropertyTagModel(res))
		}

		return pm, nil
	})
	if err != nil {
		return nil, err
	}
	p := res.(*model.PropertyModel)

	return p, nil
}

func (r *repo) SearchPropertyCombination(ctx context.Context, query *dto.SearchPropertyCombinationQuery) (*dto.SearchPropertyCombinationResponse, error) {
	sqSql, args := sqlbuilders.SearchPropertyCombinationBuilder(query)
	rows, err := r.dao.Query(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res1, err := func() (*dto.SearchPropertyCombinationResponse, error) {
		defer rows.Close()
		var r dto.SearchPropertyCombinationResponse
		for rows.Next() {
			var i dto.SearchPropertyCombinationItem
			if err := rows.Scan(&i.LId, &r.Count); err != nil {
				return nil, err
			}
			r.Items = append(r.Items, i)
		}
		return &r, nil
	}()

	if err != nil {
		return nil, err
	}

	return &dto.SearchPropertyCombinationResponse{
		Count:  res1.Count,
		SortBy: *query.SortBy,
		Order:  *query.Order,
		Offset: *query.Offset,
		Limit:  *query.Limit,
		Items:  res1.Items,
	}, nil
}

func (r *repo) GetPropertiesByIds(ctx context.Context, ids []string, fields []string) ([]model.PropertyModel, error) {
	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"features", "tags", "media"}, f) {
			fkFields = append(fkFields, f)
		} else {
			nonFKFields = append(nonFKFields, f)
		}
	}
	// log.Println(nonFKFields, fkFields)

	// get non fk fields
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("properties")
	ib.Where(ib.In("id::text", sqlbuilder.List(ids)))
	query, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.PropertyModel
	var i database.Property
	var scanningFields []interface{} = []interface{}{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "name":
			scanningFields = append(scanningFields, &i.Name)
		case "building":
			scanningFields = append(scanningFields, &i.Building)
		case "project":
			scanningFields = append(scanningFields, &i.Project)
		case "area":
			scanningFields = append(scanningFields, &i.Area)
		case "number_of_floors":
			scanningFields = append(scanningFields, &i.NumberOfFloors)
		case "year_built":
			scanningFields = append(scanningFields, &i.YearBuilt)
		case "orientation":
			scanningFields = append(scanningFields, &i.Orientation)
		case "entrance_width":
			scanningFields = append(scanningFields, &i.EntranceWidth)
		case "facade":
			scanningFields = append(scanningFields, &i.Facade)
		case "full_address":
			scanningFields = append(scanningFields, &i.FullAddress)
		case "city":
			scanningFields = append(scanningFields, &i.City)
		case "district":
			scanningFields = append(scanningFields, &i.District)
		case "ward":
			scanningFields = append(scanningFields, &i.Ward)
		case "lat":
			scanningFields = append(scanningFields, &i.Lat)
		case "lng":
			scanningFields = append(scanningFields, &i.Lng)
		case "place_url":
			scanningFields = append(scanningFields, &i.PlaceUrl)
		case "type":
			scanningFields = append(scanningFields, &i.Type)
		case "description":
			scanningFields = append(scanningFields, &i.Description)
		case "is_public":
			scanningFields = append(scanningFields, &i.IsPublic)
		case "created_at":
			scanningFields = append(scanningFields, &i.CreatedAt)
		case "updated_at":
			scanningFields = append(scanningFields, &i.UpdatedAt)
		}
	}
	for rows.Next() {
		if err := rows.Scan(scanningFields...); err != nil {
			return nil, err
		}
		items = append(items, *model.ToPropertyModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get fk fields
	for i := 0; i < len(items); i++ {
		p := &items[i]
		if slices.Contains(fkFields, "features") {
			f, err := r.dao.GetPropertyFeatures(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, fdb := range f {
				p.Features = append(p.Features, *model.ToPropertyFeatureModel(&fdb))
			}
		}
		if slices.Contains(fkFields, "media") {
			m, err := r.dao.GetPropertyMedia(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, mdb := range m {
				p.Media = append(p.Media, *model.ToPropertyMediaModel(&mdb))
			}
		}
		if slices.Contains(fkFields, "tags") {
			t, err := r.dao.GetPropertyTags(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, tdb := range t {
				p.Tags = append(p.Tags, model.PropertyTagModel(tdb))
			}
		}

	}
	return items, nil
}

func (r *repo) GetPropertyById(ctx context.Context, id uuid.UUID) (*model.PropertyModel, error) {
	p, err := r.dao.GetPropertyById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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

func (r *repo) GetUnitsOfProperty(ctx context.Context, id uuid.UUID) ([]unitModel.UnitModel, error) {
	_res, err := r.dao.GetUnitsOfProperty(ctx, id)
	if err != nil {
		return nil, err
	}
	resm := make([]unitModel.UnitModel, len(_res))
	for _, res := range _res {
		resm = append(resm, *unitModel.ToUnitModel(&res))
	}
	return resm, nil
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

func (r *repo) GetManagedProperties(ctx context.Context, userId uuid.UUID) ([]database.GetManagedPropertiesRow, error) {
	return r.dao.GetManagedProperties(ctx, userId)
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

func (r *repo) IsPublic(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.dao.IsPropertyPublic(ctx, id)
}

func (r *repo) UpdateProperty(ctx context.Context, data *dto.UpdateProperty) error {
	return r.dao.UpdateProperty(ctx, *data.ToUpdatePropertyDB())
}

func (r *repo) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteProperty(ctx, id)
}
