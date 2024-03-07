package dto

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	unitDTO "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
)

type SearchListingQuery struct {
	LTitle                *string    `json:"ltitle"`
	LCreatorID            *string    `json:"lcreatorId"`
	LPropertyID           *string    `json:"lpropertyId"`
	LMinPrice             *int64     `json:"lminPrice"`
	LMaxPrice             *int64     `json:"lmaxPrice"`
	LPriceNegotiable      *bool      `json:"lpriceNegotiable"`
	LSecurityDeposit      *int64     `json:"lsecurityDeposit"`
	LLeaseTerm            *int32     `json:"lleaseTerm"`
	LPetsAllowed          *bool      `json:"lpetsAllowed"`
	LMinNumberOfResidents *int32     `json:"lminNumberOfResidents"`
	LPriority             *int32     `json:"lpriority"`
	LActive               *bool      `json:"lactive"`
	LPolicies             []int32    `json:"lpolicies"`
	LMinCreatedAt         *time.Time `json:"lminCreatedAt"`
	LMaxCreatedAt         *time.Time `json:"lmaxCreatedAt"`
	LMinUpdatedAt         *time.Time `json:"lminUpdatedAt"`
	LMaxUpdatedAt         *time.Time `json:"lmaxUpdatedAt"`
	LMinPostAt            *time.Time `json:"lminPostAt"`
	LMaxPostAt            *time.Time `json:"lmaxPostAt"`
	LMinExpiredAt         *time.Time `json:"lminExpiredAt"`
	LMaxExpiredAt         *time.Time `json:"lmaxExpiredAt"`
}

type SearchListingCombinationQuery struct {
	requests.SearchSortPaginationQuery
	SearchListingQuery
	property_dto.SearchPropertyQuery
	unitDTO.SearchUnitQuery
}

func (q *SearchListingCombinationQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.PTypes) == 1 {
		q.PTypes = strings.Split(q.PTypes[0], ",")
	}
	// if len(q.PFeatures) == 1 {
	// 	q.PFeatures = strings.Split(q.PFeatures[0], ",")
	// }
	if len(q.PTags) == 1 {
		q.PTags = strings.Split(q.PTags[0], ",")
	}
	// if len(q.LPolicies) == 1 {
	// 	q.LPolicies = strings.Split(q.LPolicies[0], ",")
	// }
	// if len(q.UAmenities) == 1 {
	// 	q.UAmenities = strings.Split(q.UAmenities[0], ",")
	// }
	return nil
}

type SearchListingCombinationItem struct {
	LId uuid.UUID `json:"lid"`
}

type SearchListingCombinationResponse struct {
	Count  uint32                         `json:"count"`
	Limit  int32                          `json:"limit"`
	Offset int32                          `json:"offset"`
	SortBy string                         `json:"sortby"`
	Order  string                         `json:"order"`
	Items  []SearchListingCombinationItem `json:"items"`
}
