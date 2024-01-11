package dto

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

type UnitFieldQuery struct {
	Fields []string `query:"fields" validate:"unitFields"`
}

func ValidateQuery(fl validator.FieldLevel) bool {
	if fields, ok := fl.Field().Interface().([]string); ok {
		for _, f := range fields {
			if !slices.Contains([]string{"name", "property_id", "area", "floor", "price", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "type", "created_at", "updated_at", "amenities", "media"}, f) {
				return false
			}
		}
		return true
	}
	return false
}
