package property

import (
	"log"
	"testing"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

func TestProperty(t *testing.T) {
	sql, args := SearchPropertyBuilder(
		[]string{"properties.id", "properties.name"},
		&dto.SearchPropertyQuery{
			PTypes:          []string{"APARTMENT", "HOUSE", "TEST"},
			PName:           types.Ptr[string]("test name"),
			PBuilding:       types.Ptr[string]("test building"),
			PProject:        types.Ptr[string]("test project"),
			PFullAddress:    types.Ptr[string]("test address"),
			PCity:           types.Ptr[string]("test city"),
			PDistrict:       types.Ptr[string]("test district"),
			PWard:           types.Ptr[string]("test ward"),
			PMinArea:        types.Ptr[float32](3000),
			PMaxArea:        types.Ptr[float32](2000),
			PNumberOfFloors: types.Ptr[int32](12),
			PYearBuilt:      types.Ptr[int32](2023),
			POrientation:    types.Ptr[string]("nw"),
			PMinFacade:      types.Ptr[int32](12),
			PIsPublic:       types.Ptr[bool](true),
			PFeatures:       []int32{1, 2, 3},
			PTags:           []string{"tag 1", "tag 2"},
			PMinCreatedAt:   types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
			PMaxCreatedAt:   types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		}, "", "")
	log.Println(sql)
	log.Println(args)
}
