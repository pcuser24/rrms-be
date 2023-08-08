package property

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

func TestXxx(t *testing.T) {
	desc := "abcd efgh"
	amenities := []dto.CreatePropertyAmenity{
		{
			Amenity:     "a1",
			Description: &desc,
		},
		{
			Amenity:     "a2",
			Description: nil,
		},
	}
	uid, _ := uuid.Parse("d01bfb0b-dfbf-442f-8674-b0823b5eac60")

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("amenities")
	ib.Cols("property_id", "amenity", "description")
	for _, amenity := range amenities {
		ib.Values(uid, amenity.Amenity, types.StrN((amenity.Description)))
	}
	sql, args := ib.Build()
	fmt.Println(sql)
	fmt.Println(args)

}
