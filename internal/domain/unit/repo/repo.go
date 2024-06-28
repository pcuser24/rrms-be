package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/domain/unit/repo/sqlbuild"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/redisd"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

const CACHE_EXPIRATION = 24 * 60 * 60 * time.Second // 24h0m0s

type Repo interface {
	CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (*model.UnitModel, error)
	GetUnitsByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]model.UnitModel, error)
	GetUnitsOfProperty(ctx context.Context, pid uuid.UUID) ([]model.UnitModel, error)
	SearchUnitCombination(ctx context.Context, query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error)
	CheckUnitManageability(ctx context.Context, uid uuid.UUID, userId uuid.UUID) (bool, error)
	CheckUnitOfProperty(ctx context.Context, pid, uid uuid.UUID) (bool, error)
	IsPublic(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateUnit(ctx context.Context, data *dto.UpdateUnit) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
}

type repo struct {
	dao         database.DAO
	redisClient redisd.RedisClient
}

func NewRepo(d database.DAO, redisClient redisd.RedisClient) Repo {
	return &repo{
		dao:         d,
		redisClient: redisClient,
	}
}

func (r *repo) CreateUnit(ctx context.Context, data *dto.CreateUnit) (*model.UnitModel, error) {
	var um *model.UnitModel
	res, err := r.dao.CreateUnit(ctx, *data.ToCreateUnitDB())
	if err != nil {
		return nil, err
	}
	um = model.ToUnitModel(&res)
	err = func() error {
		for _, a := range data.Amenities {
			amenity, err := r.dao.CreateUnitAmenity(ctx, database.CreateUnitAmenityParams{
				UnitID:      res.ID,
				AmenityID:   a.AmenityID,
				Description: types.StrN(a.Description),
			})
			if err != nil {
				return err
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
				return err
			}
			um.Media = append(um.Media, *model.ToUnitMediaModel(&media))
		}

		return nil
	}()

	if err != nil {
		_ = r.dao.DeleteUnit(ctx, res.ID)
		return nil, err
	}

	// save to cache
	r.saveUnitToCache(ctx, um)
	r.saveUnitsOfPropertyToCache(ctx, res.PropertyID, []uuid.UUID{res.ID})

	return um, nil
}

func (r *repo) GetUnitById(ctx context.Context, id uuid.UUID) (*model.UnitModel, error) {
	// read from cache
	uc, err := r.readUnitFromCache(ctx, id)
	if err == nil && uc != nil {
		return uc, nil
	}

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

	// save to cache
	r.saveUnitToCache(ctx, um)
	r.saveUnitsOfPropertyToCache(ctx, u.PropertyID, []uuid.UUID{u.ID})

	return um, nil
}

func (r *repo) GetUnitsOfProperty(ctx context.Context, pid uuid.UUID) ([]model.UnitModel, error) {
	// read from cache
	uids, err := r.readUnitsOfPropertyFromCache(ctx, pid)
	if err == nil && len(uids) > 0 {
		ucs, _, err := r.readUnitsFromCache(ctx, uids)
		if err == nil && len(ucs) == len(uids) {
			return ucs, nil
		}
	}

	_res, err := r.dao.GetUnitsOfProperty(ctx, pid)
	if err != nil {
		return nil, err
	}
	res := make([]model.UnitModel, 0, len(_res))
	for _, u := range _res {
		um := model.ToUnitModel(&u)
		a, err := r.dao.GetUnitAmenities(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		for _, adb := range a {
			um.Amenities = append(um.Amenities, *model.ToUnitAmenityModel(&adb))
		}

		m, err := r.dao.GetUnitMedia(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		for _, mdb := range m {
			um.Media = append(um.Media, *model.ToUnitMediaModel(&mdb))
		}
		res = append(res, *um)

		// save to cache
		if r.redisClient.Exists(ctx, fmt.Sprintf("unit:%s", u.ID.String())).Val() == 0 {
			r.saveUnitToCache(ctx, um)
		}
	}
	return res, nil
}

func (r *repo) SearchUnitCombination(ctx context.Context, query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error) {
	sqlUnit, argsUnit := sqlbuild.SearchUnitBuilder(
		[]string{"units.id", "count(*) OVER() AS full_count"},
		&query.SearchUnitQuery,
		"", "",
	)

	sql, args := sqlbuilder.Build(sqlUnit, argsUnit...).Build()
	sqSql := utils.SequelizePlaceholders(sql)
	// sqSql += fmt.Sprintf(" ORDER BY %v %v", utils.PtrDerefence[string](query.SortBy, "created_at"), utils.PtrDerefence[string](query.Order, "desc"))
	sqSql += fmt.Sprintf(" LIMIT %v", utils.PtrDerefence[int32](query.Limit, 1000))
	sqSql += fmt.Sprintf(" OFFSET %v", utils.PtrDerefence[int32](query.Offset, 0))
	rows, err := r.dao.Query(context.Background(), sqSql, args...)
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

func (r *repo) GetUnitsByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]model.UnitModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// read from cache
	ucs, cacheMiss, err := r.readUnitsFromCache(ctx, ids)
	if err == nil {
		if len(ucs) == len(ids) {
			return ucs, nil
		} else {
			ids = cacheMiss
		}
	}

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
	ib.Where(ib.In("id::text", sqlbuilder.List(func() []string {
		var res []string
		for _, id := range ids {
			res = append(res, id.String())
		}
		return res
	}())))
	query, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, query, args...)
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
	rows.Close()
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

	// combine cached and non-cached items
	items = append(items, ucs...)
	return items, nil
}

func (r *repo) saveUnitToCache(ctx context.Context, p *model.UnitModel) error {
	dataMap, err := json.Marshal(p)
	if err != nil {
		return err
	}
	setRes := r.redisClient.Set(ctx, fmt.Sprintf("unit:%s", p.ID.String()), dataMap, CACHE_EXPIRATION)
	return setRes.Err()
}

func (r *repo) saveUnitsOfPropertyToCache(ctx context.Context, pid uuid.UUID, units []uuid.UUID) error {
	uidStrs := make([]string, 0, len(units))
	for _, u := range units {
		uidStrs = append(uidStrs, u.String())
	}
	res := r.redisClient.Sadd(ctx, fmt.Sprintf("property:%s:units", pid.String()), uidStrs)
	return res.Err()
}

func (r *repo) readUnitsOfPropertyFromCache(ctx context.Context, pid uuid.UUID) ([]uuid.UUID, error) {
	cacheRes := r.redisClient.SMembers(ctx, fmt.Sprintf("property:%s:units", pid.String()))
	if err := cacheRes.Err(); err != nil {
		return nil, err
	}
	// set expiration
	r.redisClient.Expire(ctx, fmt.Sprintf("property:%s:units", pid.String()), CACHE_EXPIRATION)

	data := cacheRes.Val()
	res := make([]uuid.UUID, 0, len(data))
	for _, datum := range data {
		id, err := uuid.Parse(datum)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *repo) readUnitFromCache(ctx context.Context, id uuid.UUID) (*model.UnitModel, error) {
	cacheRes := r.redisClient.Get(ctx, fmt.Sprintf("unit:%s", id.String()))
	if err := cacheRes.Err(); err != nil {
		return nil, err
	}
	dataStr := cacheRes.Val()
	if dataStr == "" {
		return nil, nil
	}
	var p model.UnitModel
	if err := json.Unmarshal([]byte(dataStr), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) readUnitsFromCache(ctx context.Context, ids []uuid.UUID) (res []model.UnitModel, cacheMiss []uuid.UUID, err error) {
	cacheRes := r.redisClient.MGet(ctx, func() []string {
		var res []string
		for _, id := range ids {
			res = append(res, fmt.Sprintf("unit:%s", id.String()))
		}
		return res
	}()...)
	if err := cacheRes.Err(); err != nil {
		return nil, nil, err
	}
	data := cacheRes.Val()
	for i, datum := range data {
		datumStr, ok := datum.(string)
		if !ok {
			cacheMiss = append(cacheMiss, ids[i])
			continue
		}
		var p model.UnitModel
		if err := json.Unmarshal([]byte(datumStr), &p); err != nil {
			return nil, nil, err
		}
		res = append(res, p)
	}
	return res, cacheMiss, nil
}
