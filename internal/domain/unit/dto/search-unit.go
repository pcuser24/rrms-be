package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
)

type SearchUnitQuery struct {
	UPropertyID          *string `query:"upropertyId" validate:"omitempty,uuid"`
	UName                *string `query:"uname" validate:"omitempty"`
	UMinArea             *int64  `query:"uminPrice" validate:"omitempty"`
	UMaxArea             *int64  `query:"umaxPrice" validate:"omitempty"`
	UFloor               *int32  `query:"ufloor" validate:"omitempty"`
	UPrice               int64   `query:"uprice" validate:"omitempty"`
	UNumberOfLivingRooms *int32  `query:"unumberOfLivingRooms" validate:"omitempty"`
	UNumberOfBedrooms    *int32  `query:"unumberOfBedrooms" validate:"omitempty"`
	UNumberOfBathrooms   *int32  `query:"unumberOfBathrooms" validate:"omitempty"`
	UNumberOfToilets     *int32  `query:"unumberOfToilets" validate:"omitempty"`
	UNumberOfKitchens    *int32  `query:"unumberOfKitchens" validate:"omitempty"`
	UNumberOfBalconies   *int32  `query:"unumberOfBalconies" validate:"omitempty"`
	UAmenities           []int32 `query:"uamenities" validate:"omitempty"`
}

type SearchUnitCombinationQuery struct {
	requests.SearchSortPaginationQuery
	SearchUnitQuery
}

type SearchUnitCombinationItem struct {
	UId uuid.UUID `json:"uid"`
}

type SearchUnitCombinationResponse struct {
	Count uint32                      `json:"count"`
	Items []SearchUnitCombinationItem `json:"items"`
}
