package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

var retrievableFields = []string{"name", "property_id", "area", "floor", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "type", "created_at", "updated_at", "amenities", "media"}

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
}

const UnitFieldsLocalKey = "unitFields"

type UnitFieldQuery struct {
	Fields []string `query:"fields" validate:"unitFields"`
}

func (q *UnitFieldQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
	return nil
}

type GetUnitsByIdsQuery struct {
	UnitFieldQuery
	IDs []uuid.UUID `query:"unitIds" validate:"required,dive,uuid4"`
}

func (q *GetUnitsByIdsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		fieldSet := set.NewSet[string]()
		fieldSet.AddAll(strings.Split(q.Fields[0], ",")...)
		q.Fields = fieldSet.ToSlice()
	}
	return nil
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
