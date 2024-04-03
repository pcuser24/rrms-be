package repo

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/domain/listing/repo/sqlbuild"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListingsByIds(ctx context.Context, ids []string, fields []string) ([]model.ListingModel, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(ctx context.Context, data *dto.UpdateListing) error
	UpdateListingStatus(ctx context.Context, id uuid.UUID, active bool) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckListingVisibility(ctx context.Context, lid, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckListingExpired(ctx context.Context, lid uuid.UUID) (bool, error)
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
	var lm *model.ListingModel
	res, err := r.dao.CreateListing(ctx, *data.ToCreateListingDB())
	if err != nil {
		return nil, err
	}
	lm = model.ToListingModel(&res)

	err = func() error {
		for i := 0; i < len(data.Units); i++ {
			u := &data.Units[i]
			lu, err := r.dao.CreateListingUnit(ctx, database.CreateListingUnitParams{
				ListingID: lm.ID,
				UnitID:    u.UnitID,
				Price:     u.Price,
			})
			if err != nil {
				return err
			}
			lm.Units = append(lm.Units, model.ListingUnitModel(lu))
		}

		for i := 0; i < len(data.Policies); i++ {
			p := &data.Policies[i]
			lp, err := r.dao.CreateListingPolicy(ctx, *p.ToCreateListingPolicyDB(lm.ID))
			if err != nil {
				return err
			}
			lm.Policies = append(lm.Policies, model.ToListingPolicyModel(&lp))
		}

		for i := 0; i < len(data.Tags); i++ {
			lt, err := r.dao.CreateListingTag(ctx, database.CreateListingTagParams{
				ListingID: lm.ID,
				Tag:       data.Tags[i],
			})
			if err != nil {
				return err
			}
			lm.Tags = append(lm.Tags, model.ListingTagModel(lt))
		}

		return nil
	}()

	if err != nil {
		_ = r.dao.DeleteListing(ctx, lm.ID)
		return nil, err
	}

	return lm, nil
}

func (r *repo) GetListingsByIds(ctx context.Context, ids []string, fields []string) ([]model.ListingModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
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
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, query, args...)
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
	rows.Close()
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
				l.Policies = append(l.Policies, model.ToListingPolicyModel(&i))
			}
		}
		if slices.Contains(fkFields, "tags") {
			t, err := r.dao.GetListingTags(ctx, l.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range t {
				l.Tags = append(l.Tags, model.ListingTagModel(i))
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
		res.Policies = append(res.Policies, model.ToListingPolicyModel(&i))
	}

	u, err := r.dao.GetListingUnits(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, i := range u {
		res.Units = append(res.Units, model.ListingUnitModel(i))
	}

	t, err := r.dao.GetListingTags(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, i := range t {
		res.Tags = append(res.Tags, model.ListingTagModel(i))
	}

	return res, nil
}

func (r *repo) UpdateListing(ctx context.Context, data *dto.UpdateListing) error {
	return r.dao.UpdateListing(ctx, *data.ToUpdateListingDB())
}

func (r *repo) UpdateListingStatus(ctx context.Context, id uuid.UUID, active bool) error {
	return r.dao.UpdateListingStatus(ctx, database.UpdateListingStatusParams{
		ID:     id,
		Active: active,
	})
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

func (r *repo) CheckListingExpired(ctx context.Context, lid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckListingExpired(ctx, lid)
	if err != nil {
		return false, err
	}
	return res.Bool, nil
}

func (r *repo) SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	sqSql, args := sqlbuild.SearchListingCombinationBuilder(query)
	rows, err := r.dao.Query(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res1, err := func() (*dto.SearchListingCombinationResponse, error) {
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

	return &dto.SearchListingCombinationResponse{
		SortBy: query.SortBy,
		Order:  query.Order,
		Limit:  *query.Limit,
		Offset: *query.Offset,
		Items:  res1.Items,
		Count:  res1.Count,
	}, nil
}

func (r *repo) CheckListingVisibility(ctx context.Context, lid, uid uuid.UUID) (bool, error) {
	return r.dao.CheckListingVisibility(ctx, database.CheckListingVisibilityParams{
		ID:        lid,
		ManagerID: uid,
	})
}
