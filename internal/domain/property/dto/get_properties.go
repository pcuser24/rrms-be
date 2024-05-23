package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const PropertyFieldsLocalKey = "propertyFields"

var retrievableFields = []string{"name", "building", "project", "area", "number_of_floors", "year_built", "orientation", "entrance_width", "facade", "full_address", "city", "district", "ward", "lat", "lng", "primary_image", "description", "type", "is_public", "created_at", "updated_at", "features", "tags", "media"}

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
}

type GetPropertiesQuery struct {
	Fields []string `query:"fields" validate:"propertyFields"`
	Limit  *int32   `query:"limit" validate:"omitempty,gte=0"`
	Offset *int32   `query:"offset" validate:"omitempty,gte=0"`
	SortBy *string  `query:"sortby" validate:"omitempty,oneof=created_at area name rentals"`
	Order  *string  `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (q *GetPropertiesQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
	return nil
}

type GetPropertiesByIdsQuery struct {
	GetPropertiesQuery
	IDs []uuid.UUID `query:"propIds" validate:"required,dive,uuid4"`
}

func (q *GetPropertiesByIdsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	return q.GetPropertiesQuery.QueryParser(ctx)
}

func ValidateQuery(fl validator.FieldLevel) bool {
	if fields, ok := fl.Field().Interface().([]string); ok {
		for _, f := range fields {
			if !slices.Contains(retrievableFields, f) {
				return false
			}
		}
		return true
	}
	return false
}
