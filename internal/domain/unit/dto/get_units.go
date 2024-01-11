package dto

import "github.com/google/uuid"

type GetUnitByIdsQuery struct {
	Fields []string    `query:"fields" validate:"unitFields"`
	IDs    []uuid.UUID `query:"unitIds" validate:"required,dive,uuid4"`
}
