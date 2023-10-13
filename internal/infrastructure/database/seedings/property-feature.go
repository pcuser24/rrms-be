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
		"p_feature-security",
		"p_feature-fire_alarm",
		"p_feature-gym",
		"p_feature-fitness_center",
		"p_feature-swimming_pool",
		"p_feature-community_rooms",
		"p_feature-public_library",
		"p_feature-parking",
		"p_feature-outdoor_common_area",
		"p_feature-services",
		"p_feature-facilities",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
