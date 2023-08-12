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
		"p-feature_security",
		"p-feature_fire-alarm",
		"p-feature_gym",
		"p-feature_fitness-center",
		"p-feature_swimming-pool",
		"p-feature_community-rooms",
		"p-feature_public-library",
		"p-feature_parking",
		"p-feature_outdoor-common-area",
		"p-feature_services",
		"p-feature_facilities",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
