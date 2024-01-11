package dto

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

type GetListingsQuery struct {
	Fields []string `query:"fields" validate:"listingFields"`
}

var ValidateQuery validator.Func = func(fl validator.FieldLevel) bool {
	if fields, ok := fl.Field().Interface().([]string); ok {
		for _, f := range fields {
			if !slices.Contains([]string{"creator_id", "property_id", "title", "description", "full_name", "email", "phone", "contact_type", "price", "price_negotiable", "security_deposit", "lease_term", "pets_allowed", "number_of_residents", "priority", "post_at", "active", "created_at", "updated_at", "post_at", "expired_at", "policies", "units"}, f) {
				return false
			}
		}
		return true
	}
	return false
}
