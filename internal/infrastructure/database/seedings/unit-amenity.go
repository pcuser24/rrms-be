package seedings

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func SeedUnitAmenities(d database.DAO) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("u_amenities")
	ib.Cols("amenity")
	for _, i := range []string{
		"u_amenity-furniture",
		"u_amenity-fridge",
		"u_amenity-air-cond",
		"u_amenity-washing-machine",
		"u_amenity-dishwasher",
		"u_amenity-water-heater",
		"u_amenity-tv",
		"u_amenity-internet",
		"u_amenity-wardrobe",
		"u_amenity-entresol",
		"u_amenity-bed",
		"u_amenity-other",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
