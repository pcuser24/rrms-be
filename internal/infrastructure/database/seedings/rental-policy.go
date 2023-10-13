package seedings

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func SeedRentalPolicies(d database.DAO) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("rental_policies")
	ib.Cols("policy")
	for _, i := range []string{
		"rental_policy-payment",
		"rental_policy-maintenance",
		"rental_policy-insurance",
		"rental_policy-noise",
		"rental_policy-lease_renewal",
		"rental_policy-change_to_property",
		"rental_policy-parking",
		"rental_policy-pets",
		"rental_policy-subletting",
		"rental_policy-business",
		"rental_policy-consequences",
		"rental_policy-other",
	} {
		ib.Values(i)
	}
	sql, args := ib.Build()
	_, err := d.ExecContext(context.Background(), sql, args...)
	return err
}
