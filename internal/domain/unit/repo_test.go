package unit

import (
	"testing"

	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

func TestSearchUnit(t *testing.T) {
	sql, args := sqlbuilders.SearchUnitBuilder([]string{"1"}, &dto.SearchUnitQuery{
		UPropertyID:          types.Ptr[string]("1"),
		UName:                types.Ptr[string]("test name"),
		UMinArea:             types.Ptr[int64](3000),
		UMaxArea:             types.Ptr[int64](2000),
		UFloor:               types.Ptr[int32](12),
		UNumberOfLivingRooms: types.Ptr[int32](12),
		UNumberOfBedrooms:    types.Ptr[int32](12),
		UNumberOfBathrooms:   types.Ptr[int32](12),
		UNumberOfToilets:     types.Ptr[int32](12),
		UNumberOfKitchens:    types.Ptr[int32](12),
		UNumberOfBalconies:   types.Ptr[int32](12),
	}, "", "")
	t.Log(sql)
	t.Log(args)
}
