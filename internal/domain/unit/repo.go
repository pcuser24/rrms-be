package unit

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (*model.UnitModel, error)
	GetUnitsByIds(ctx context.Context, ids []string, fields []string) ([]model.UnitModel, error)
	GetUnitsOfProperty(ctx context.Context, id uuid.UUID) ([]model.UnitModel, error)
	SearchUnitCombination(ctx context.Context, query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error)
	CheckUnitManageability(ctx context.Context, uid uuid.UUID, userId uuid.UUID) (bool, error)
	CheckUnitOfProperty(ctx context.Context, pid, uid uuid.UUID) (bool, error)
	IsPublic(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateUnit(ctx context.Context, data *dto.UpdateUnit) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
	GetAllAmenities(ctx context.Context) ([]model.UAmenity, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d database.DAO) (interface{}, error) {
		var um *model.UnitModel
		res, err := d.CreateUnit(ctx, *data.ToCreateUnitDB())
		if err != nil {
			return nil, err
		}
		um = model.ToUnitModel(&res)

		for _, a := range data.Amenities {
			amenity, err := r.dao.CreateUnitAmenity(ctx, database.CreateUnitAmenityParams{
				UnitID:      res.ID,
				AmenityID:   a.AmenityID,
				Description: types.StrN(a.Description),
			})
			if err != nil {
				return nil, err
			}
			um.Amenities = append(um.Amenities, *model.ToUnitAmenityModel(&amenity))
		}

		for _, m := range data.Media {
			media, err := r.dao.CreateUnitMedia(ctx, database.CreateUnitMediaParams{
				UnitID:      res.ID,
				Url:         m.Url,
				Type:        m.Type,
				Description: types.StrN(m.Description),
			})
			if err != nil {
				return nil, err
			}
			um.Media = append(um.Media, *model.ToUnitMediaModel(&media))
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
		um.Media = append(um.Media, *model.ToUnitMediaModel(&mdb))
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
			um.Media = append(um.Media, *model.ToUnitMediaModel(&mdb))
		}
		res = append(res, um)
	}
	return res, nil
}

func (r *repo) SearchUnitCombination(ctx context.Context, query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error) {
	sqlUnit, argsUnit := sqlbuilders.SearchUnitBuilder(
		[]string{"units.id", "count(*) OVER() AS full_count"},
		&query.SearchUnitQuery,
		"", "",
	)

	sql, args := sqlbuilder.Build(sqlUnit, argsUnit...).Build()
	sqSql := utils.SequelizePlaceholders(sql)
	sqSql += fmt.Sprintf(" ORDER BY %v %v", utils.PtrDerefence[string](query.SortBy, "created_at"), utils.PtrDerefence[string](query.Order, "desc"))
	sqSql += fmt.Sprintf(" LIMIT %v", utils.PtrDerefence[int32](query.Limit, 1000))
	sqSql += fmt.Sprintf(" OFFSET %v", utils.PtrDerefence[int32](query.Offset, 0))
	rows, err := r.dao.QueryContext(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res, err := func() (*dto.SearchUnitCombinationResponse, error) {
		defer rows.Close()
		var r dto.SearchUnitCombinationResponse
		for rows.Next() {
			var i dto.SearchUnitCombinationItem
			if err := rows.Scan(&i.UId, &r.Count); err != nil {
				return nil, err
			}
			r.Items = append(r.Items, i)
		}
		return &r, nil
	}()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) UpdateUnit(ctx context.Context, data *dto.UpdateUnit) error {
	return r.dao.UpdateUnit(ctx, *data.ToUpdateUnitDB())
}

func (r *repo) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteUnit(ctx, id)
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

func (r *repo) CheckUnitManageability(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error) {
	res, err := r.dao.CheckUnitManageability(ctx, database.CheckUnitManageabilityParams{
		ID:        id,
		ManagerID: userId,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) CheckUnitOfProperty(ctx context.Context, pid, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckUnitOfProperty(ctx, database.CheckUnitOfPropertyParams{
		ID:         uid,
		PropertyID: pid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) IsPublic(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.dao.IsUnitPublic(ctx, id)
}

func (r *repo) GetUnitsByIds(ctx context.Context, ids []string, fields []string) ([]model.UnitModel, error) {
	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"amenities", "media"}, f) {
			fkFields = append(fkFields, f)
		} else {
			nonFKFields = append(nonFKFields, f)
		}
	}
	// log.Println(nonFKFields, fkFields)

	// get non fk fields
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("units")
	ib.Where(ib.In("id::text", sqlbuilder.List(ids)))
	query, args := ib.Build()
	log.Println(query, args)
	rows, err := r.dao.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.UnitModel
	var i database.Unit
	var scanningFields []interface{} = []interface{}{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "property_id":
			scanningFields = append(scanningFields, &i.PropertyID)
		case "name":
			scanningFields = append(scanningFields, &i.Name)
		case "area":
			scanningFields = append(scanningFields, &i.Area)
		case "floor":
			scanningFields = append(scanningFields, &i.Floor)
		case "price":
			scanningFields = append(scanningFields, &i.Price)
		case "number_of_living_rooms":
			scanningFields = append(scanningFields, &i.NumberOfLivingRooms)
		case "number_of_bedrooms":
			scanningFields = append(scanningFields, &i.NumberOfBedrooms)
		case "number_of_bathrooms":
			scanningFields = append(scanningFields, &i.NumberOfBathrooms)
		case "number_of_toilets":
			scanningFields = append(scanningFields, &i.NumberOfToilets)
		case "number_of_balconies":
			scanningFields = append(scanningFields, &i.NumberOfBalconies)
		case "number_of_kitchens":
			scanningFields = append(scanningFields, &i.NumberOfKitchens)
		case "type":
			scanningFields = append(scanningFields, &i.Type)
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
		items = append(items, *model.ToUnitModel(&i))
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get fk fields
	for i := 0; i < len(items); i++ {
		u := &items[i]
		if slices.Contains(fkFields, "amenities") {
			ua, err := r.dao.GetUnitAmenities(ctx, u.ID)
			if err != nil {
				return nil, err
			}
			for _, a := range ua {
				u.Amenities = append(u.Amenities, *model.ToUnitAmenityModel(&a))
			}
		}
		if slices.Contains(fkFields, "media") {
			um, err := r.dao.GetUnitMedia(ctx, u.ID)
			if err != nil {
				return nil, err
			}
			for _, m := range um {
				u.Media = append(u.Media, *model.ToUnitMediaModel(&m))
			}
		}
	}

	return items, nil
}
