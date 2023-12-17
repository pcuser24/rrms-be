package dto

import (
	"time"

	"github.com/google/uuid"
	propertyDTO "github.com/user2410/rrms-backend/internal/domain/property/dto"
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
	propertyDTO.SearchPropertyQuery
	unitDTO.SearchUnitQuery
}

type SearchListingCombinationItem struct {
	LId         uuid.UUID `json:"lid"`
	LPropertyID uuid.UUID `json:"lpropertyId"`
	LTitle      string    `json:"ltitle"`
	LPrice      int64     `json:"lprice"`
	LPriority   int32     `json:"lpriority"`
	LPostAt     time.Time `json:"lpostAt"`
}

type SearchListingCombinationResponse struct {
	Count uint32                         `json:"count"`
	Items []SearchListingCombinationItem `json:"items"`
}
