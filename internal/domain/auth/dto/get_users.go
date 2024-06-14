package dto

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

var (
	retrievableFields = []string{"email", "password", "group_id", "created_at", "updated_at", "created_by", "updated_by", "deleted_f", "first_name", "last_name", "phone", "avatar", "address", "city", "district", "ward", "role"}
)

const UserFieldsLocalKey = "userFields"

func GetRetrievableFields() []string {
	rfs := make([]string, len(retrievableFields))
	copy(rfs, retrievableFields)
	return rfs
}

type UserFieldQuery struct {
	Fields []string `query:"fields" validate:"userFields"`
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
