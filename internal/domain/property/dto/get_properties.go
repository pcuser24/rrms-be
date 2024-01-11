package dto

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

type GetPropertiesQuery struct {
	Fields []string `query:"fields" validate:"propertyFields"`
}

func ValidQuery(fl validator.FieldLevel) bool {
	if fields, ok := fl.Field().Interface().([]string); ok {
		for _, f := range fields {
			if !slices.Contains([]string{"name", "building", "project", "area", "number_of_floors", "year_built", "orientation", "entrance_width", "facade", "full_address", "city", "district", "ward", "lat", "lng", "place_url", "description", "type", "is_public", "created_at", "updated_at", "features", "tags", "media"}, f) {
				return false
			}
		}
		return true
	}
	return false
}
