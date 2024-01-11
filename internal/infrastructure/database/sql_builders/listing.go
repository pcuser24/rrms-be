package sqlbuilders

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
)

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
