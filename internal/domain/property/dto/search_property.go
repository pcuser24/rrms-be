package dto

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	unitDTO "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
)

type SearchPropertyQuery struct {
	PTypes          []string   `query:"ptypes" validate:"omitempty"`
	PCreatorID      *string    `query:"pcreatorId" validate:"omitempty"`
	PName           *string    `query:"pname" validate:"omitempty"`
	PBuilding       *string    `query:"pbuilding" validate:"omitempty"`
	PProject        *string    `query:"pproject" validate:"omitempty"`
	PFullAddress    *string    `query:"pfullAddress" validate:"omitempty"`
	PCity           *string    `query:"pcity" validate:"omitempty"`
	PDistrict       *string    `query:"pdistrict" validate:"omitempty"`
	PWard           *string    `query:"pward" validate:"omitempty"`
	PMinArea        *float32   `query:"pminArea" validate:"omitempty"`
	PMaxArea        *float32   `query:"pmaxArea" validate:"omitempty"`
	PNumberOfFloors *int32     `query:"pnumberOfFloors" validate:"omitempty"`
	PYearBuilt      *int32     `query:"pyearBuilt" validate:"omitempty"`
	POrientation    *string    `query:"porientation" validate:"omitempty"`
	PMinFacade      *int32     `query:"pfacade" validate:"omitempty"`
	PIsPublic       *bool      `query:"pisPublic" validate:"omitempty"`
	PFeatures       []int32    `query:"pfeatures" validate:"omitempty"`
	PMinCreatedAt   *time.Time `query:"pminCreatedAt" validate:"omitempty"`
	PMaxCreatedAt   *time.Time `query:"pmaxCreatedAt" validate:"omitempty"`
	PMinUpdatedAt   *time.Time `query:"pminUpdatedAt" validate:"omitempty"`
	PMaxUpdatedAt   *time.Time `query:"pmaxUpdatedAt" validate:"omitempty"`
}

type SearchPropertyCombinationQuery struct {
	requests.SearchSortPaginationQuery
	SearchPropertyQuery
	unitDTO.SearchUnitQuery
}

func (q *SearchPropertyCombinationQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.PTypes) == 1 {
		q.PTypes = strings.Split(q.PTypes[0], ",")
	}
	return nil
}

type SearchPropertyCombinationItem struct {
	LId uuid.UUID `json:"lid"`
}

type SearchPropertyCombinationResponse struct {
	Count  uint32                          `json:"count"`
	Limit  int32                           `json:"limit"`
	Offset int32                           `json:"offset"`
	SortBy []string                        `json:"sortby"`
	Order  []string                        `json:"order"`
	Items  []SearchPropertyCombinationItem `json:"items"`
}
