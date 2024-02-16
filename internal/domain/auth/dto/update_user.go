package dto

import "github.com/google/uuid"

type UpdateUser struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=8,max=32"`
	FirstName *string `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName  *string `json:"last_name" validate:"omitempty,min=1,max=32"`
	Phone     *string `json:"phone" validate:"omitempty,min=10,max=15"`
	Avatar    *string `json:"avatar" validate:"omitempty,url"`
	Address   *string `json:"address" validate:"omitempty,min=1,max=255"`
	City      *string `json:"city" validate:"omitempty,min=1,max=255"`
	District  *string `json:"district" validate:"omitempty,min=1,max=255"`
	Ward      *string `json:"ward" validate:"omitempty,min=1,max=255"`
	UpdatedBy uuid.UUID
}
