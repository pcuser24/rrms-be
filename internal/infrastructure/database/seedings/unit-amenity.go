package seedings

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func SeedUnitAmenities(d database.DAO) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("u_amenities")
	ib.Cols("feature")
	for _, i := range []string{
		"u-amenity_fridge",
		"u-amenity_air-cond",
		"u-amenity_washing-machine",
		"u-amenity_dishwasher",
		"u-amenity_water-heater",
		"u-amenity_tv",
		"u-amenity_internet",
		"u-amenity_wardrobe",
		"u-amenity_closet",
		"u-amenity_entresol",
		"u-amenity_bed",
		"u-amenity_sofa",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
