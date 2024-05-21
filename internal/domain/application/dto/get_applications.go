package dto

import (
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const ApplicationFieldsLocalKey = "applicationFields"

var retrievableFields = []string{"creator_id", "listing_id", "property_id", "unit_id", "listing_price", "offered_price", "status", "created_at", "updated_at", "full_name", "email", "phone", "dob", "profile_image", "movein_date", "preferred_term", "rental_intention", "organization_name", "organization_hq_address", "organization_scale", "rh_address", "rh_city", "rh_district", "rh_ward", "rh_rental_duration", "rh_monthly_payment", "rh_reason_for_leaving", "employment_status", "employment_company_name", "employment_position", "employment_monthly_income", "employment_comment", "minors", "coaps", "pets", "vehicles"}

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
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

type GetApplicationsQuery struct {
	Fields []string `query:"fields" validate:"applicationFields"`
}

func (q *GetApplicationsQuery) parseFields() {
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
}

func (q *GetApplicationsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}

type GetApplicationsByIdsQuery struct {
	GetApplicationsQuery
	IDs []int64 `query:"appIds" validate:"required,dive"`
}

func (q *GetApplicationsByIdsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}

type GetApplicationsToMeQuery struct {
	GetApplicationsQuery
	CreatedBefore time.Time `query:"createdBefore"`
	Limit         int32     `query:"limit" validate:"omitempty,gte=0"`
	Offset        int32     `query:"offset" validate:"omitempty,gte=0"`
}

func (q *GetApplicationsToMeQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	if q.CreatedBefore.IsZero() {
		q.CreatedBefore = time.Now().AddDate(0, 0, -30)
	}
	if q.Limit == 0 {
		q.Limit = 10
	}
	return nil
}

type GetApplicationsOfPropertyQuery struct {
	GetApplicationsQuery
	ListingIds []uuid.UUID `query:"listingIds" validate:"dive,uuid4"`
	Limit      *int32      `query:"limit" validate:"omitempty,gte=0"`
	Offset     *int32      `json:"offset" validate:"omitempty,gte=0"`
}

func (q *GetApplicationsOfPropertyQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}
