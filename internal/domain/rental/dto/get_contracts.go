package dto

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

const RentalContractFieldsLocalKey = "rentalContractFields"

var (
	contractRetrievableFields = []string{"rental_id", "a_fullname", "a_dob", "a_phone", "a_address", "a_household_registration", "a_identity", "a_identity_issued_by", "a_identity_issued_at", "a_documents", "a_bank_account", "a_bank", "a_registration_number", "b_fullname", "b_organization_name", "b_organization_hq_address", "b_organization_code", "b_organization_code_issued_at", "b_organization_code_issued_by", "b_dob", "b_phone", "b_address", "b_household_registration", "b_identity", "b_identity_issued_by", "b_identity_issued_at", "b_bank_account", "b_bank", "b_tax_code", "payment_method", "payment_day", "n_copies", "created_at_place", "content", "status", "created_at", "updated_at", "created_by", "updated_by"}
)

type GetRentalContracts struct {
	Fields []string `query:"fields" validate:"rentalContractFields"`
	Limit  *int32   `query:"limit" validate:"omitempty,gte=0"`
	Offset *int32   `query:"offset" validate:"omitempty,gte=0"`
}

func (q *GetRentalContracts) parseFields() {
	if len(q.Fields) == 1 {
		q.Fields = strings.Split(q.Fields[0], ",")
	}
}

func (q *GetRentalContracts) QueryParser(ctx *fiber.Ctx) error {
	err := ctx.QueryParser(q)
	if err != nil {
		return err
	}
	q.parseFields()
	return nil
}

func (q *GetRentalContracts) ValidateQuery() error {
	v := validator.New()
	v.RegisterValidation(
		RentalContractFieldsLocalKey,
		func(fl validator.FieldLevel) bool {
			if fields, ok := fl.Field().Interface().([]string); ok {
				for _, f := range fields {
					if !slices.Contains(contractRetrievableFields, f) {
						return false
					}
				}
				return true
			}
			return false
		},
	)
	if errs := validation.ValidateStruct(v, *q); len(errs) > 0 {
		return errs[0]
	}
	return nil
}
