package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const ListingFieldsLocalKey = "listingFields"

var retrievableFields = []string{"creator_id", "property_id", "title", "description", "full_name", "email", "phone", "contact_type", "price", "price_negotiable", "security_deposit", "lease_term", "pets_allowed", "number_of_residents", "priority", "active", "created_at", "updated_at", "expired_at", "policies", "units"}

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
}

type GetListingsByIdsQuery struct {
	GetListingsQuery
	IDs []uuid.UUID `query:"listingIds" validate:"required,dive,uuid4"`
}

func (q *GetListingsByIdsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
	return nil
}

type GetListingsQuery struct {
	Fields []string `query:"fields" validate:"listingFields"`
}

func (q *GetListingsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
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

type GetListingsOfPropertyQuery struct {
	GetListingsQuery
	Expired bool   `query:"expired" validate:"omitempty"`
	Limit   *int32 `query:"limit" validate:"omitempty,gt=0"`
	Offset  *int32 `json:"offset" validate:"omitempty,gt=0"`
}
