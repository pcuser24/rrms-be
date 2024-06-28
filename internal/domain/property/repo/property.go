package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/domain/property/repo/sqlbuild"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

const CACHE_EXPIRATION = 24 * 60 * 60 * time.Second

func (r *repo) CreateProperty(ctx context.Context, data *property_dto.CreateProperty) (*property_model.PropertyModel, error) {
	var pm *property_model.PropertyModel
	prop, err := r.dao.CreateProperty(ctx, *data.ToCreatePropertyDB())
	if err != nil {
		return nil, err
	}
	pm = property_model.ToPropertyModel(&prop)

	err = func() error {
		for _, m := range data.Managers {
			res, err := r.dao.CreatePropertyManager(ctx, database.CreatePropertyManagerParams{
				PropertyID: prop.ID,
				ManagerID:  m.ManagerID,
				Role:       m.Role,
			})
			if err != nil {
				return err
			}
			pm.Managers = append(pm.Managers, property_model.PropertyManagerModel(res))
		}

		for _, f := range data.Features {
			res, err := r.dao.CreatePropertyFeature(ctx, *f.ToCreatePropertyFeatureDB(prop.ID))
			if err != nil {
				return err
			}
			pm.Features = append(pm.Features, property_model.ToPropertyFeatureModel(&res))
		}

		var primaryImageID int64
		for _, m := range data.Media {
			res, err := r.dao.CreatePropertyMedia(ctx, *m.ToCreatePropertyMediaDB(prop.ID))
			if err != nil {
				return err
			}
			if m.Type == database.MEDIATYPEIMAGE && res.Url == data.PrimaryImage {
				primaryImageID = res.ID
			}
			pm.Media = append(pm.Media, property_model.ToPropertyMediaModel(&res))
		}
		err = r.dao.UpdateProperty(ctx, database.UpdatePropertyParams{
			ID:           prop.ID,
			PrimaryImage: pgtype.Int8{Valid: true, Int64: primaryImageID},
		})
		if err != nil {
			return err
		}
		pm.PrimaryImage = primaryImageID

		for _, t := range data.Tags {
			res, err := r.dao.CreatePropertyTag(ctx, database.CreatePropertyTagParams{
				PropertyID: prop.ID,
				Tag:        t.Tag,
			})
			if err != nil {
				return err
			}
			pm.Tags = append(pm.Tags, property_model.PropertyTagModel(res))
		}

		return nil
	}()

	if err != nil {
		// rollback and ignore any error
		_ = r.dao.DeleteProperty(ctx, pm.ID)
		return nil, err
	}

	// save the new property to redis
	r.savePropertyToCache(ctx, pm)

	return pm, nil
}

func (r *repo) SearchPropertyCombination(ctx context.Context, query *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error) {
	sqSql, args := sqlbuild.SearchPropertyCombinationBuilder(query)
	rows, err := r.dao.Query(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res1, err := func() (*property_dto.SearchPropertyCombinationResponse, error) {
		defer rows.Close()
		var r property_dto.SearchPropertyCombinationResponse
		for rows.Next() {
			var i property_dto.SearchPropertyCombinationItem
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

	return &property_dto.SearchPropertyCombinationResponse{
		Count:  res1.Count,
		SortBy: query.SortBy,
		Order:  query.Order,
		Offset: *query.Offset,
		Limit:  *query.Limit,
		Items:  res1.Items,
	}, nil
}

func (r *repo) GetPropertiesByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]property_model.PropertyModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// read from cache
	cachedProperties, cacheMiss, cerr := r.readPropertiesFromCache(ctx, ids)
	if cerr == nil {
		if len(cachedProperties) == len(ids) {
			return cachedProperties, nil
		} else {
			ids = cacheMiss
		}
	}

	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"features", "tags", "media", "managers"}, f) {
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
	ib.Where(ib.In("properties.id::text", sqlbuilder.List(func() []string {
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
	var items []property_model.PropertyModel
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
		case "primary_image":
			scanningFields = append(scanningFields, &i.PrimaryImage)
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
		items = append(items, *property_model.ToPropertyModel(&i))
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
				p.Features = append(p.Features, property_model.ToPropertyFeatureModel(&fdb))
			}
		}
		if slices.Contains(fkFields, "media") {
			m, err := r.dao.GetPropertyMedia(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, mdb := range m {
				p.Media = append(p.Media, property_model.ToPropertyMediaModel(&mdb))
			}
		}
		if slices.Contains(fkFields, "managers") {
			t, err := r.dao.GetPropertyManagers(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, tdb := range t {
				p.Managers = append(p.Managers, property_model.PropertyManagerModel(tdb))
			}
		}
	}

	// combine cached properties with fetched properties from db
	items = append(items, cachedProperties...)
	return items, nil
}

func (r *repo) GetPropertyById(ctx context.Context, id uuid.UUID) (*property_model.PropertyModel, error) {
	// read from cache
	cachedProperty, err := r.readPropertyFromCache(ctx, id)
	if err == nil && cachedProperty != nil {
		return cachedProperty, nil
	}

	p, err := r.dao.GetPropertyById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	pm := property_model.ToPropertyModel(&p)

	mn, err := r.dao.GetPropertyManagers(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, mdb := range mn {
		pm.Managers = append(pm.Managers, property_model.PropertyManagerModel(mdb))
	}

	f, err := r.dao.GetPropertyFeatures(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, fdb := range f {
		pm.Features = append(pm.Features, property_model.ToPropertyFeatureModel(&fdb))
	}

	t, err := r.dao.GetPropertyTags(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, tdb := range t {
		pm.Tags = append(pm.Tags, property_model.PropertyTagModel(tdb))
	}

	m, err := r.dao.GetPropertyMedia(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, mdb := range m {
		pm.Media = append(pm.Media, property_model.ToPropertyMediaModel(&mdb))
	}

	// save to cache
	r.savePropertyToCache(ctx, pm)

	return pm, nil
}

func (r *repo) GetListingsOfProperty(ctx context.Context, id uuid.UUID, query *listing_dto.GetListingsOfPropertyQuery) ([]uuid.UUID, error) {
	params := database.GetListingsOfPropertyParams{
		PropertyID: id,
		Expired:    query.Expired,
	}
	if query.Offset != nil {
		params.Offset = *query.Offset
	} else {
		params.Offset = 0
	}
	if query.Limit != nil {
		params.Limit = *query.Limit
	} else {
		params.Limit = 100
	}

	return r.dao.GetListingsOfProperty(ctx, params)
}

func (r *repo) GetApplicationsOfProperty(ctx context.Context, id uuid.UUID, query *application_dto.GetApplicationsOfPropertyQuery) ([]int64, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id")
	sb.From("applications")
	andExprs := []string{
		sb.Equal("property_id::text", id.String()),
	}
	if len(query.ListingIds) > 0 {
		idsStr := make([]string, 0, len(query.ListingIds))
		for _, id := range query.ListingIds {
			idsStr = append(idsStr, id.String())
		}
		andExprs = append(andExprs, sb.In("listing_id::text", sqlbuilder.List(idsStr)))
	}
	sb.Where(andExprs...)
	if query.Limit != nil {
		sb.Limit(int(*query.Limit))
	} else {
		sb.Limit(100)
	}
	if query.Offset != nil {
		sb.Offset(int(*query.Offset))
	} else {
		sb.Offset(0)
	}
	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repo) GetPropertyManagers(ctx context.Context, id uuid.UUID) ([]property_model.PropertyManagerModel, error) {
	// read from cache
	cachedProperty, err := r.readPropertyFromCache(ctx, id)
	if err == nil && cachedProperty != nil {
		return cachedProperty.Managers, nil
	}

	res, err := r.dao.GetPropertyManagers(ctx, id)
	if err != nil {
		return nil, err
	}
	var items []property_model.PropertyManagerModel
	for _, i := range res {
		items = append(items, property_model.PropertyManagerModel(i))
	}
	return items, err
}

type GetManagedPropertiesRow struct {
	PropertyID uuid.UUID
	Role       string
}

func (r *repo) GetManagedProperties(ctx context.Context, userId uuid.UUID, query *property_dto.GetPropertiesQuery) ([]GetManagedPropertiesRow, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("property_id", "role")
	sb.From("property_managers")
	sb.Where(sb.Equal("manager_id", userId))
	if query.SortBy != nil {
		switch *query.SortBy {
		case "rentals":
			sb.OrderBy("(SELECT COUNT(*) FROM rentals WHERE rentals.property_id = property_managers.property_id)")
		default:
			sb.OrderBy(fmt.Sprintf("(SELECT %s FROM properties WHERE properties.id = property_managers.property_id)", *query.SortBy))
		}
	}
	if query.Order != nil && *query.Order == "asc" {
		sb.Asc()
	} else {
		sb.Desc()
	}
	if query.Limit != nil {
		sb.Limit(int(*query.Limit))
	} else {
		sb.Limit(100)
	}
	if query.Offset != nil {
		sb.Offset(int(*query.Offset))
	} else {
		sb.Offset(0)
	}

	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetManagedPropertiesRow
	for rows.Next() {
		var i GetManagedPropertiesRow
		if err := rows.Scan(&i.PropertyID, &i.Role); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repo) GetRentalsOfProperty(ctx context.Context, id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]int64, error) {
	params := database.GetRentalsOfPropertyParams{
		PropertyID: id,
		Expired:    query.Expired,
	}
	if query.Offset != nil {
		params.Offset = *query.Offset
	} else {
		params.Offset = 0
	}
	if query.Limit != nil {
		params.Limit = *query.Limit
	} else {
		params.Limit = 100
	}
	return r.dao.GetRentalsOfProperty(ctx, params)
}

func (r *repo) IsPropertyVisible(ctx context.Context, uid, pid uuid.UUID) (bool, error) {
	res, err := r.dao.IsPropertyVisible(ctx, database.IsPropertyVisibleParams{
		UserID:     uid,
		PropertyID: pid,
	})
	if err != nil {
		return false, err
	}
	return res.Bool, nil
}

func (r *repo) UpdateProperty(ctx context.Context, data *property_dto.UpdateProperty) error {
	// update property media
	if data.Media != nil {
		// delete all old media records
		sb := sqlbuilder.PostgreSQL.NewDeleteBuilder()
		sb.DeleteFrom("property_media")
		sb.Where(sb.Equal("property_id::text", data.ID.String()))
		sql, args := sb.Build()
		_, err := r.dao.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		// then insert new media records
		media := []database.PropertyMedium{}
		for _, m := range data.Media {
			newMedia, err := r.dao.CreatePropertyMedia(ctx, database.CreatePropertyMediaParams{
				PropertyID:  data.ID,
				Url:         m.Url,
				Type:        m.Type,
				Description: types.StrN(m.Description),
			})
			if err != nil {
				return err
			}
			media = append(media, newMedia)
		}
		err = r.dao.UpdateProperty(ctx, database.UpdatePropertyParams{
			ID:           data.ID,
			PrimaryImage: pgtype.Int8{Valid: true, Int64: media[*data.PrimaryImage].ID}, // data.PrimaryImage must exist
		})
		if err != nil {
			return err
		}
		data.PrimaryImage = nil
	}

	// update property features
	if data.Features != nil {
		// delete all old features records
		sb := sqlbuilder.PostgreSQL.NewDeleteBuilder()
		sb.DeleteFrom("property_features")
		sb.Where(sb.Equal("property_id::text", data.ID.String()))
		sql, args := sb.Build()
		_, err := r.dao.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		// then insert new features records
		for _, f := range data.Features {
			_, err := r.dao.CreatePropertyFeature(ctx, database.CreatePropertyFeatureParams{
				PropertyID:  data.ID,
				FeatureID:   f.FeatureID,
				Description: types.StrN(f.Description),
			})
			if err != nil {
				return nil
			}
		}
	}

	err := r.dao.UpdateProperty(ctx, data.ToUpdatePropertyDB())
	if err != nil {
		return err
	}

	// update cache
	property, err := r.GetPropertyById(ctx, data.ID)
	if err == nil {
		r.savePropertyToCache(ctx, property)
	}

	return nil
}

func (r *repo) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	return r.dao.DeleteProperty(ctx, id)
}

func (r *repo) FilterVisibleProperties(ctx context.Context, pids []uuid.UUID, uid uuid.UUID) ([]uuid.UUID, error) {
	buildSQL := func(pid uuid.UUID) (string, []interface{}) {
		return `
		SELECT (
			SELECT is_public FROM "properties" WHERE properties.id = $1 LIMIT 1
		) OR (
			SELECT EXISTS (SELECT 1 FROM "property_managers" WHERE property_managers.property_id = $1 AND property_managers.manager_id = $2 LIMIT 1)
		) OR (
			SELECT EXISTS (SELECT 1 FROM "new_property_manager_requests" WHERE new_property_manager_requests.property_id = $1 AND new_property_manager_requests.user_id = $2 LIMIT 1)
		) OR (
			SELECT EXISTS (SELECT 1 FROM "rentals" WHERE rentals.property_id = $1 AND rentals.tenant_id = $2 LIMIT 1)
		)
		`, []interface{}{pid, uid}
	}

	resIDs := make([]uuid.UUID, 0, len(pids))

	queries := make([]database.BatchedQueryRow, 0, len(pids))
	for _, uid := range pids {
		sql, args := buildSQL(uid)
		queries = append(queries, database.BatchedQueryRow{
			SQL:    sql,
			Params: args,
			Fn: func(row pgx.Row) error {
				var res bool
				if err := row.Scan(&res); err != nil {
					return err
				}
				resIDs = append(resIDs, uid)
				return nil
			},
		})
	}

	err := r.dao.QueryRowBatch(ctx, queries)

	return resIDs, err
}

func (r *repo) savePropertyToCache(ctx context.Context, p *property_model.PropertyModel) error {
	dataMap, err := json.Marshal(p)
	if err != nil {
		return err
	}
	setRes := r.redisClient.Set(ctx, fmt.Sprintf("property:%s", p.ID.String()), dataMap, CACHE_EXPIRATION)
	return setRes.Err()
}

func (r *repo) readPropertyFromCache(ctx context.Context, id uuid.UUID) (*property_model.PropertyModel, error) {
	cacheRes := r.redisClient.Get(ctx, fmt.Sprintf("property:%s", id.String()))
	if err := cacheRes.Err(); err != nil {
		return nil, err
	}
	dataStr := cacheRes.Val()
	if dataStr == "" {
		return nil, nil
	}
	var p property_model.PropertyModel
	if err := json.Unmarshal([]byte(dataStr), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) readPropertiesFromCache(ctx context.Context, ids []uuid.UUID) (res []property_model.PropertyModel, cacheMiss []uuid.UUID, err error) {
	cacheRes := r.redisClient.MGet(ctx, func() []string {
		var res []string
		for _, id := range ids {
			res = append(res, fmt.Sprintf("property:%s", id.String()))
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
		var p property_model.PropertyModel
		if err := json.Unmarshal([]byte(datumStr), &p); err != nil {
			return nil, nil, err
		}
		res = append(res, p)
	}
	return res, cacheMiss, nil
}
