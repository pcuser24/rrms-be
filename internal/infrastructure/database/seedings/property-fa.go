package seedings

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func SeedPropertyFeatures(d database.DAO) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("p_features")
	ib.Cols("feature")
	for _, i := range []string{
		"Security",
		"Fire alarm",
		"Gym",
		"Fitness center",
		"Swimming Pool",
		"Community rooms",
		"Public library",
		"Parking",
		"Outdoor common area",
		"Services",
		"Facilities",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}

func SeedPropertyAmenities(d database.DAO) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("p_amenities")
	ib.Cols("amenity")
	for _, i := range []string{
		"Elevator",
		"Security camera",
		"Yard",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
