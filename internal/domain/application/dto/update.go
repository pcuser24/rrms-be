package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type UpdateApplicationStatus struct {
	Status  database.APPLICATIONSTATUS `json:"status" validate:"required,oneof=PENDING CONDITIONALLY_APPROVED APPROVED REJECTED"`
	Comment *string                    `json:"comment" validate:"omitempty"`
}

func (u *UpdateApplicationStatus) Validate() error {
	validator := validator.New()
	err := validator.Struct(*u)
	if err != nil {
		return err
	}
	if u.Status == "REJECTED" {
		return validator.Var(u.Comment, "required")
	}
	return nil
}
