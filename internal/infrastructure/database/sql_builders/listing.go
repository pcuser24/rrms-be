package sqlbuilders

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/utils"
)

func SearchListingBuilder(
	searchFields []string, query *dto.SearchListingQuery,
	connectID, connectCreator, connectProperty string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM listings"
	var searchQueries []string // WHERE field (=/ILIKE) ?
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

	// no field is specified and check for exisence only
	if len(searchQueries) == 0 && searchFields[0] == "1" {
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
	if len(searchQueries) > 0 {
		// some fields are specified
		searchQuery += " WHERE " + strings.Join(searchQueries, " AND ")
	}

	return searchQuery, args
}

// Search listing combined with property and unit
func SearchListingCombinationBuilder(query *dto.SearchListingCombinationQuery) (string, []any) {
	sqlListing, argsListing := SearchListingBuilder(
		[]string{"listings.id", "count(*) OVER() AS full_count"},
		&query.SearchListingQuery,
		"", "", "",
	)
	sqlProp, argsProp := SearchPropertyBuilder([]string{"1"}, &query.SearchPropertyQuery, "listings.property_id", "")
	sqlUnit, argsUnit := SearchUnitBuilder([]string{"1"}, &query.SearchUnitQuery, "", "listings.property_id")

	var queryStr string = sqlListing
	var argsLs []any = argsListing
	// NOTE: goofy code, will be refactored later
	if len(argsListing) > 0 {
		if len(sqlProp) > 0 {
			queryStr += fmt.Sprintf(" AND EXISTS (%v)", sqlProp)
			argsLs = append(argsLs, argsProp...)
		}
		if len(sqlUnit) > 0 {
			queryStr += fmt.Sprintf(" AND EXISTS (%v)", sqlUnit)
			argsLs = append(argsLs, argsUnit...)
		}
	} else if len(sqlProp) > 0 || len(sqlUnit) > 0 {
		queryStr += " WHERE "
		if len(sqlProp) > 0 {
			queryStr += fmt.Sprintf("EXISTS (%v)", sqlProp)
			argsLs = append(argsLs, argsProp...)
		}
		if len(sqlUnit) > 0 {
			if len(sqlProp) > 0 {
				queryStr += " AND "
			}
			queryStr += fmt.Sprintf("EXISTS (%v)", sqlUnit)
			argsLs = append(argsLs, argsUnit...)
		}
	}

	sql, args := sqlbuilder.Build(queryStr, argsLs...).Build()
	sqSql := utils.SequelizePlaceholders(sql)

	sqSql += fmt.Sprintf(" ORDER BY %v %v", *query.SortBy, *query.Order)
	sqSql += fmt.Sprintf(" LIMIT %v", *query.Limit)
	sqSql += fmt.Sprintf(" OFFSET %v", *query.Offset)

	return sqSql, args
}
