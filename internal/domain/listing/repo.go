package listing

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Repo interface {
	CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(ctx context.Context, data *dto.UpdateListing) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	AddListingPolicies(ctx context.Context, lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error)
	AddListingUnits(ctx context.Context, lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error)
	DeleteListingPolicies(ctx context.Context, lid uuid.UUID, ids []int64) error
	DeleteListingUnits(ctx context.Context, lid uuid.UUID, ids []uuid.UUID) error
	CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error)
}

type repo struct {
	dao db.DAO
}

func NewRepo(d db.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateListing(ctx context.Context, data *dto.CreateListing) (*model.ListingModel, error) {
	res, err := r.dao.QueryTx(ctx, func(d db.DAO) (interface{}, error) {
		var lm *model.ListingModel
		res, err := d.CreateListing(ctx, *data.ToCreateListingDB())
		if err != nil {
			return nil, err
		}
		lm = model.ToListingModel(&res)

		lm.Policies, err = r.AddListingPolicies(ctx, res.ID, data.Policies)
		if err != nil {
			return nil, err
		}

		lm.Units, err = r.AddListingUnits(ctx, res.ID, data.Units)
		if err != nil {
			return nil, err
		}

		return lm, nil
	})
	if err != nil {
		return nil, err
	}

	l := res.(*model.ListingModel)

	return l, nil
}

func (r *repo) AddListingPolicies(ctx context.Context, lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error) {
	var res []model.ListingPolicyModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("listing_policies")
	ib.Cols("listing_id", "policy_id", "note")
	for _, i := range items {
		ib.Values(lid, i.PolicyID, types.StrN(i.Note))
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()

	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.ListingPolicyModel, error) {
		defer rows.Close()
		var items []model.ListingPolicyModel
		for rows.Next() {
			var i db.ListingPolicy
			if err := rows.Scan(
				&i.ListingID,
				&i.PolicyID,
				&i.Note,
			); err != nil {
				return nil, err
			}
			items = append(items, *model.ToListingPolicyModel(&i))
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

func (r *repo) AddListingUnits(ctx context.Context, lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error) {
	var res []model.ListingUnitModel
	if len(items) == 0 {
		return res, nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("listing_unit")
	ib.Cols("listing_id", "unit_id")
	for _, i := range items {
		ib.Values(lid, i.UnitID)
	}
	ib.SQL("RETURNING *")
	sql, args := ib.Build()

	rows, err := r.dao.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	res, err = func() ([]model.ListingUnitModel, error) {
		defer rows.Close()
		var items []model.ListingUnitModel
		for rows.Next() {
			var i db.ListingUnit
			if err := rows.Scan(
				&i.ListingID,
				&i.UnitID,
			); err != nil {
				return nil, err
			}
			items = append(items, model.ListingUnitModel(i))
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

func (r *repo) SearchListing(ctx context.Context, data *dto.SearchListingQuery) error {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.
		Select("*").
		From("listings").
		Where(
			sb.Equal("active", true),
		)
	return nil
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

func (r *repo) bulkDelete(ctx context.Context, uid uuid.UUID, ids []interface{}, table_name, info_id_field string) error {
	if len(ids) == 0 {
		return nil
	}

	ib := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	ib.DeleteFrom(table_name)
	ib.Where(
		ib.Equal("listing_id", uid),
		ib.In(info_id_field, ids...),
	)
	sql, args := ib.Build()
	_, err := r.dao.ExecContext(ctx, sql, args...)
	return err
}

func (r *repo) DeleteListingPolicies(ctx context.Context, lid uuid.UUID, ids []int64) error {
	ids_i := make([]interface{}, len(ids))
	for i, v := range ids {
		ids_i[i] = v
	}
	return r.bulkDelete(ctx, lid, ids_i, "listing_policies", "policy_id")
}

func (r *repo) DeleteListingUnits(ctx context.Context, lid uuid.UUID, ids []uuid.UUID) error {
	ids_i := make([]interface{}, len(ids))
	for i, v := range ids {
		ids_i[i] = v
	}
	return r.bulkDelete(ctx, lid, ids_i, "listing_unit", "unit_id")
}

func (r *repo) CheckListingOwnership(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckListingOwnership(ctx, db.CheckListingOwnershipParams{
		ID:        lid,
		CreatorID: uid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repo) CheckValidUnitForListing(ctx context.Context, lid uuid.UUID, uid uuid.UUID) (bool, error) {
	res, err := r.dao.CheckValidUnitForListing(ctx, db.CheckValidUnitForListingParams{
		ID:   uid,
		ID_2: lid,
	})
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func SearchListingBuilder(
	searchFields []string, query *dto.SearchListingQuery,
	connectID, connectCreator, connectProperty string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM listings WHERE "
	var searchQueries []string
	var args []interface{}

	if query.LTitle != nil {
		searchQueries = append(searchQueries, "listings.title ILIKE $?")
		args = append(args, "%"+(*query.LTitle)+"%")
	}
	if query.LCreatorID != nil {
		searchQueries = append(searchQueries, "listings.creator_id = $?")
		args = append(args, *query.LCreatorID)
	}
	if query.LPropertyID != nil {
		searchQueries = append(searchQueries, "listings.property_id = $?")
		args = append(args, *query.LPropertyID)
	}
	if query.LMinPrice != nil {
		searchQueries = append(searchQueries, "listings.price >= $?")
		args = append(args, *query.LMinPrice)
	}
	if query.LMaxPrice != nil {
		searchQueries = append(searchQueries, "listings.price <= $?")
		args = append(args, *query.LMaxPrice)
	}
	if query.LPriceNegotiable != nil {
		searchQueries = append(searchQueries, "listings.price_negotiable = $?")
		args = append(args, *query.LPriceNegotiable)
	}
	if query.LSecurityDeposit != nil {
		searchQueries = append(searchQueries, "listings.security_deposit = $?")
		args = append(args, *query.LSecurityDeposit)
	}
	if query.LLeaseTerm != nil {
		searchQueries = append(searchQueries, "listings.lease_term = $?")
		args = append(args, *query.LLeaseTerm)
	}
	if query.LPetsAllowed != nil {
		searchQueries = append(searchQueries, "listings.pets_allowed = $?")
		args = append(args, *query.LPetsAllowed)
	}
	if query.LMinNumberOfResidents != nil {
		searchQueries = append(searchQueries, "listings.number_of_residents >= $?")
		args = append(args, *query.LMinNumberOfResidents)
	}
	if query.LPriority != nil {
		searchQueries = append(searchQueries, "listings.priority = $?")
		args = append(args, *query.LPriority)
	}
	if query.LActive != nil {
		searchQueries = append(searchQueries, "listings.active = $?")
		args = append(args, *query.LActive)
	}
	if query.LMinCreatedAt != nil {
		searchQueries = append(searchQueries, "listings.created_at >= $?")
		args = append(args, *query.LMinCreatedAt)
	}
	if query.LMaxCreatedAt != nil {
		searchQueries = append(searchQueries, "listings.created_at <= $?")
		args = append(args, *query.LMaxCreatedAt)
	}
	if query.LMinUpdatedAt != nil {
		searchQueries = append(searchQueries, "listings.updated_at >= $?")
		args = append(args, *query.LMinUpdatedAt)
	}
	if query.LMaxUpdatedAt != nil {
		searchQueries = append(searchQueries, "listings.updated_at <= $?")
		args = append(args, *query.LMaxUpdatedAt)
	}
	if query.LMinPostAt != nil {
		searchQueries = append(searchQueries, "listings.post_at >= $?")
		args = append(args, *query.LMinPostAt)
	}
	if query.LMaxPostAt != nil {
		searchQueries = append(searchQueries, "listings.post_at <= $?")
		args = append(args, *query.LMaxPostAt)
	}
	if query.LMinExpiredAt != nil {
		searchQueries = append(searchQueries, "listings.expired_at >= $?")
		args = append(args, *query.LMinExpiredAt)
	}
	if query.LMaxExpiredAt != nil {
		searchQueries = append(searchQueries, "listings.expired_at <= $?")
		args = append(args, *query.LMaxExpiredAt)
	}
	if len(query.LPolicies) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM listing_policies WHERE listing_id = listings.id AND policy_id IN ($?))")
		args = append(args, sqlbuilder.List(query.LPolicies))
	}

	if len(searchQueries) == 0 {
		return "", []interface{}{}
	}
	if len(connectID) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("listings.id = %v", connectID))
	}
	if len(connectCreator) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("listings.creator_id = %v", connectCreator))
	}
	if len(connectProperty) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("listings.property_id = %v", connectProperty))
	}

	searchQuery += strings.Join(searchQueries, " AND \n")
	return searchQuery, args
}

func (r *repo) SearchListingCombination(ctx context.Context, query *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	sqlListing, argsListing := SearchListingBuilder(
		[]string{"listings.id", "listings.property_id", "listings.title", "listings.price", "listings.priority", "listings.post_at", "count(*) OVER() AS full_count"},
		&query.SearchListingQuery,
		"", "", "",
	)
	sqlProp, argsProp := property.SearchPropertyBuilder([]string{"1"}, &query.SearchPropertyQuery, "listings.property_id", "")
	sqlUnit, argsUnit := unit.SearchUnitBuilder([]string{"1"}, &query.SearchUnitQuery, "", "properties.id")

	// build order: unit -> property -> listing
	var queryStr string = sqlListing
	var argsLs []interface{} = argsListing

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
	sqSql += fmt.Sprintf(" ORDER BY %v %v", *query.SortBy, *query.Order)
	sqSql += fmt.Sprintf(" LIMIT %v", *query.Limit)
	sqSql += fmt.Sprintf(" OFFSET %v", *query.Offset)
	rows, err := r.dao.QueryContext(context.Background(), sqSql, args...)
	if err != nil {
		return nil, err
	}

	res, err := func() (*dto.SearchListingCombinationResponse, error) {
		defer rows.Close()
		var r dto.SearchListingCombinationResponse
		for rows.Next() {
			var i dto.SearchListingCombinationItem
			if err := rows.Scan(
				&i.LId,
				&i.LPropertyID,
				&i.LTitle,
				&i.LPrice,
				&i.LPriority,
				&i.LPostAt,
				&r.Count,
			); err != nil {
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
