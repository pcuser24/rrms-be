package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

const RentalFieldsLocalKey = "rentalFields"

var (
	retrievableFields = []string{"creator_id", "property_id", "unit_id", "application_id", "tenant_id", "profile_image", "tenant_type", "tenant_name", "tenant_phone", "tenant_email", "organization_name", "organization_hq_address", "start_date", "movein_date", "rental_period", "payment_type", "rental_price", "rental_payment_basis", "rental_intention", "deposit", "deposit_paid", "electricity_setup_by", "electricity_payment_type", "electricity_customer_code", "electricity_provider", "electricity_price", "water_setup_by", "water_payment_type", "water_customer_code", "water_provider", "water_price", "note", "status", "created_at", "updated_at"}
	sortbyFields      = append(retrievableFields, "remaining_time")
)

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
}

type GetRentalsQuery struct {
	Fields []string `query:"fields" validate:"rentalFields"`
	// Limit  *int32   `query:"limit" validate:"omitempty,gte=0"`
	// Offset *int32   `query:"offset" validate:"omitempty,gte=0"`
	// SortBy []string `query:"sortby" validate:"omitempty"`
	// Order  []string `query:"order" validate:"omitempty,dive,oneof=asc desc"`
}

func (q *GetRentalsQuery) parseFields() {
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
}

func (q *GetRentalsQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}

func (q *GetRentalsQuery) ValidateQuery() error {
	validator := validator.New()
	validator.RegisterValidation(RentalFieldsLocalKey, ValidateQuery)
	if errs := validation.ValidateStruct(validator, *q); len(errs) > 0 {
		return errs[0]
	}
	// for _, s := range q.SortBy {
	// 	if !slices.Contains(sortbyFields, s) {
	// 		return errors.New("invalid sortby field")
	// 	}
	// }
	// if len(q.SortBy) != len(q.Order) {
	// 	return errors.New("sortby and order must have the same length")
	// }
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

type GetRentalsOfPropertyQuery struct {
	GetRentalsQuery
	Expired bool   `query:"expired" validate:"omitempty"`
	Limit   *int32 `query:"limit" validate:"omitempty,gte=0"`
	Offset  *int32 `json:"offset" validate:"omitempty,gte=0"`
}

func (q *GetRentalsOfPropertyQuery) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}

type GetManagedRentalPaymentsQuery struct {
	Limit  *int32                         `query:"limit" validate:"omitempty,gte=0"`
	Offset *int32                         `json:"offset" validate:"omitempty,gte=0"`
	Status []database.RENTALPAYMENTSTATUS `query:"status" validate:"required,dive,oneof=PLAN ISSUED PENDING REQUEST2PAY PAID CANCELLED"`
}

type GetManagedRentalPaymentsItem struct {
	Payment  rental_model.RentalPayment    `json:"payment"`
	Rental   rental_model.RentalModel      `json:"rental"`
	Property *property_model.PropertyModel `json:"property"`
	Unit     *unit_model.UnitModel         `json:"unit"`
}
