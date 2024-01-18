package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const PropertyFieldsLocalKey = "propertyFields"

type GetPropertiesByIdsQuery struct {
	Fields []string    `query:"fields" validate:"propertyFields"`
	IDs    []uuid.UUID `query:"propIds" validate:"required,dive,uuid4"`
}

func (q *GetPropertiesByIdsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
	return nil
}

type GetPropertiesQuery struct {
	Fields []string `query:"fields" validate:"propertyFields"`
}

func ValidateQuery(fl validator.FieldLevel) bool {
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
