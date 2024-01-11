package listing

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
	"github.com/user2410/rrms-backend/internal/utils"
)

type Repo interface {
	CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListings(ctx context.Context, ids []string, fields []string) ([]model.ListingModel, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(ctx context.Context, data *dto.UpdateListing) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d database.DAO) (interface{}, error) {
		var lm *model.ListingModel
		res, err := d.CreateListing(ctx, *data.ToCreateListingDB())
		if err != nil {
			return nil, err
		}
		lm = model.ToListingModel(&res)

		for i := 0; i < len(data.Policies); i++ {
			p := &data.Policies[i]
			lp, err := r.dao.CreateListingPolicy(ctx, *p.ToCreateListingPolicyDB(lm.ID))
			if err != nil {
				return nil, err
			}
			lm.Policies = append(lm.Policies, *model.ToListingPolicyModel(&lp))
		}

		for i := 0; i < len(data.Units); i++ {
			u := &data.Units[i]
			lu, err := r.dao.CreateListingUnit(ctx, database.CreateListingUnitParams{
				ListingID: lm.ID,
				UnitID:    u.UnitID,
			})
			if err != nil {
				return nil, err
			}
			lm.Units = append(lm.Units, model.ListingUnitModel(lu))
		}

		return lm, nil
	})
	if err != nil {
		return nil, err
	}

	l := res.(*model.ListingModel)

	return l, nil
}

func (r *repo) GetListings(ctx context.Context, ids []string, fields []string) ([]model.ListingModel, error) {
	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"units", "policies"}, f) {
			fkFields = append(fkFields, f)
		} else {
			nonFKFields = append(nonFKFields, f)
		}
	}

	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("listings")
	ib.Where(ib.In("id::text", sqlbuilder.List(ids)))
	query, args := ib.Build()
	log.Println(query, args)
	rows, err := r.dao.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.ListingModel
	var i database.Listing
	var scanningFields []interface{} = []interface{}{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "creator_id":
			scanningFields = append(scanningFields, &i.CreatorID)
		case "property_id":
			scanningFields = append(scanningFields, &i.PropertyID)
		case "title":
			scanningFields = append(scanningFields, &i.Title)
		case "description":
			scanningFields = append(scanningFields, &i.Description)
		case "full_name":
			scanningFields = append(scanningFields, &i.FullName)
		case "email":
			scanningFields = append(scanningFields, &i.Email)
		case "phone":
			scanningFields = append(scanningFields, &i.Phone)
		case "contact_type":
			scanningFields = append(scanningFields, &i.ContactType)
		case "price":
			scanningFields = append(scanningFields, &i.Price)
		case "price_negotiable":
			scanningFields = append(scanningFields, &i.PriceNegotiable)
		case "security_deposit":
			scanningFields = append(scanningFields, &i.SecurityDeposit)
		case "lease_term":
			scanningFields = append(scanningFields, &i.LeaseTerm)
		case "pets_allowed":
			scanningFields = append(scanningFields, &i.PetsAllowed)
		case "number_of_residents":
			scanningFields = append(scanningFields, &i.NumberOfResidents)
		case "priority":
			scanningFields = append(scanningFields, &i.Priority)
		case "post_at":
			scanningFields = append(scanningFields, &i.PostAt)
		case "expired_at":
			scanningFields = append(scanningFields, &i.ExpiredAt)
		case "active":
			scanningFields = append(scanningFields, &i.Active)
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
		items = append(items, *model.ToListingModel(&i))
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get fk fields
	for i := 0; i < len(fkFields); i++ {
		l := &items[i]
		if slices.Contains(fkFields, "units") {
			u, err := r.dao.GetListingUnits(ctx, l.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				l.Units = append(l.Units, model.ListingUnitModel(i))
			}
		}
		if slices.Contains(fkFields, "policies") {
			p, err := r.dao.GetListingPolicies(ctx, l.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range p {
				l.Policies = append(l.Policies, *model.ToListingPolicyModel(&i))
			}
		}
	}

	return items, nil
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

func (r *repo) CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckListingOwnership(ctx, database.CheckListingOwnershipParams{
		ID:        lid,
		CreatorID: uid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckValidUnitForListing(ctx, database.CheckValidUnitForListingParams{
		ID:   uid,
		ID_2: lid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	sqlListing, argsListing := sqlbuilders.SearchListingBuilder(
		[]string{"listings.id", "count(*) OVER() AS full_count"},
		&query.SearchListingQuery,
		"", "", "",
	)
	sqlProp, argsProp := sqlbuilders.SearchPropertyBuilder([]string{"1"}, &query.SearchPropertyQuery, "listings.property_id", "")
	sqlUnit, argsUnit := sqlbuilders.SearchUnitBuilder([]string{"1"}, &query.SearchUnitQuery, "", "properties.id")

	var queryStr string = sqlListing
	var argsLs []interface{} = argsListing
	// build order: unit -> property -> listing
	if len(sqlProp) > 0 {
		var tmp string = sqlProp
		argsLs = append(argsLs, argsProp...)
		if len(sqlUnit) > 0 {
			tmp += fmt.Sprintf(" AND EXISTS (%v) ", sqlUnit)
			argsLs = append(argsLs, argsUnit...)
		}
		queryStr = fmt.Sprintf("%v AND EXISTS (%v) ", sqlListing, tmp)
	}

	sql, args := sqlbuilder.Build(queryStr, argsLs...).Build()
	sqSql := utils.SequelizePlaceholders(sql)
	sqSql += fmt.Sprintf(" ORDER BY %v %v", utils.PtrDerefence[string](query.SortBy, "created_at"), utils.PtrDerefence[string](query.Order, "desc"))
	sqSql += fmt.Sprintf(" LIMIT %v", utils.PtrDerefence[int32](query.Limit, 1000))
	sqSql += fmt.Sprintf(" OFFSET %v", utils.PtrDerefence[int32](query.Offset, 0))
	rows, err := r.dao.QueryContext(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res, err := func() (*dto.SearchListingCombinationResponse, error) {
		defer rows.Close()
		var r dto.SearchListingCombinationResponse
		for rows.Next() {
			var i dto.SearchListingCombinationItem
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

	return res, nil
}
