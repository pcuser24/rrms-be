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
	PTypes          []string   `query:"ptypes"`
	PCreatorID      *string    `query:"pcreatorId"`
	PName           *string    `query:"pname"`
	PBuilding       *string    `query:"pbuilding"`
	PProject        *string    `query:"pproject"`
	PFullAddress    *string    `query:"pfullAddress"`
	PCity           *string    `query:"pcity"`
	PDistrict       *string    `query:"pdistrict"`
	PWard           *string    `query:"pward"`
	PMinArea        *float32   `query:"pminArea"`
	PMaxArea        *float32   `query:"pmaxArea"`
	PNumberOfFloors *int32     `query:"pnumberOfFloors"`
	PYearBuilt      *int32     `query:"pyearBuilt"`
	POrientation    *string    `query:"porientation"`
	PMinFacade      *int32     `query:"pfacade"`
	PIsPublic       *bool      `query:"pisPublic"`
	PFeatures       []int32    `query:"pfeatures"`
	PTags           []string   `query:"ptags"`
	PMinCreatedAt   *time.Time `query:"pminCreatedAt"`
	PMaxCreatedAt   *time.Time `query:"pmaxCreatedAt"`
	PMinUpdatedAt   *time.Time `query:"pminUpdatedAt"`
	PMaxUpdatedAt   *time.Time `query:"pmaxUpdatedAt"`
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
	// if len(q.PFeatures) == 1 {
	// 	q.PFeatures = strings.Split(q.PFeatures[0], ",")
	// }
	if len(q.PTags) == 1 {
		q.PTags = strings.Split(q.PTags[0], ",")
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
	SortBy string                          `json:"sortby"`
	Order  string                          `json:"order"`
	Items  []SearchPropertyCombinationItem `json:"items"`
}
