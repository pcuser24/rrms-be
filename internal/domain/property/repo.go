package property

import (
	"context"
	"fmt"

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

		// insert property amenities
		ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
		ib.InsertInto("property_amenity")
		ib.Cols("property_id", "amenity", "description")
		for _, amenity := range data.Amenities {
			ib.Values(res.ID, amenity.Amenity, types.StrN((amenity.Description)))
		}
		ib.SQL("RETURNING *")
		sql, args := ib.Build()
		// fmt.Println("property_amenity: ", sql)
		// fmt.Println("property_amenity args: ", args)
		rows, err := d.QueryContext(ctx, sql, args...)
		if err != nil {
			return nil, err
		}
		pm.Amenities, err = func() ([]model.PropertyAmenityModel, error) {
			defer rows.Close()
			var items []model.PropertyAmenityModel
			for rows.Next() {
				var i db.PropertyAmenity
				if err := rows.Scan(
					&i.PropertyID,
					&i.Amenity,
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
			// fmt.Println("Read amenities error:", err)
			return nil, err
		}

		// insert property features
		ib = sqlbuilder.PostgreSQL.NewInsertBuilder()
		ib.InsertInto("property_feature")
		ib.Cols("property_id", "feature", "description")
		for _, ft := range data.Features {
			ib.Values(res.ID, ft.Feature, types.StrN((ft.Description)))
		}
		ib.SQL("RETURNING *")
		sql, args = ib.Build()
		// fmt.Println("property_feature:", sql)
		// fmt.Println("property_feature args:", args)
		rows, err = d.QueryContext(ctx, sql, args...)
		if err != nil {
			fmt.Println("Insert features error:", err)
			return nil, err
		}
		pm.Features, err = func() ([]model.PropertyFeatureModel, error) {
			defer rows.Close()
			var items []model.PropertyFeatureModel
			for rows.Next() {
				var i db.PropertyFeature
				if err := rows.Scan(
					&i.PropertyID,
					&i.Feature,
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
			// fmt.Println("Read features error:", err)
			return nil, err
		}

		// insert property medium
		ib = sqlbuilder.PostgreSQL.NewInsertBuilder()
		ib.InsertInto("property_media")
		ib.Cols("property_id", "url", "type")
		for _, media := range data.Medium {
			ib.Values(res.ID, media.Url, media.Type)
		}
		ib.SQL("RETURNING *")
		sql, args = ib.Build()
		// fmt.Println("property_media:", sql)
		// fmt.Println("property_media args:", args)
		rows, err = d.QueryContext(ctx, sql, args...)
		if err != nil {
			// fmt.Println("Insert media error:", err)
			return nil, err
		}
		pm.Medium, err = func() ([]model.PropertyMediaModel, error) {
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
			// fmt.Println("Read medium error:", err)
			return nil, err
		}

		// insert property tag
		ib = sqlbuilder.PostgreSQL.NewInsertBuilder()
		ib.InsertInto("property_tag")
		ib.Cols("property_id", "tag")
		for _, tag := range data.Tags {
			ib.Values(res.ID, tag.Tag)
		}
		ib.SQL("RETURNING *")
		sql, args = ib.Build()
		// fmt.Println("property_tag:", sql)
		// fmt.Println("property_tag args:", args)
		rows, err = d.QueryContext(ctx, sql, args...)
		if err != nil {
			fmt.Println("Insert tags error:", err)
			return nil, err
		}
		pm.Tags, err = func() ([]model.PropertyTagModel, error) {
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
			// fmt.Println("Read tags error:", err)
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
