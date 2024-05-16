package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateListing struct {
	Title             *string               `json:"title" validate:"omitempty"`
	Description       *string               `json:"description" validate:"omitempty"`
	Price             *int64                `json:"price" validate:"omitempty"`
	SecurityDeposit   *int64                `json:"securityDeposit" validate:"omitempty"`
	LeaseTerm         *int32                `json:"leaseTerm" validate:"omitempty"`
	PetsAllowed       *bool                 `json:"petsAllowed" validate:"omitempty"`
	NumberOfResidents *int32                `json:"numberOfResidents" validate:"omitempty"`
	Policies          []CreateListingPolicy `json:"policies" validate:"omitempty,dive"`
	Units             []CreateListingUnit   `json:"units" validate:"omitempty,dive"`
	Tags              []string              `json:"tags" validate:"omitempty,dive"`
}

func (u *UpdateListing) ToUpdateListingDB(id uuid.UUID) *database.UpdateListingParams {
	return &database.UpdateListingParams{
		Title:             types.StrN(u.Title),
		Description:       types.StrN(u.Description),
		Price:             types.Int64N(u.Price),
		SecurityDeposit:   types.Int64N(u.SecurityDeposit),
		LeaseTerm:         types.Int32N(u.LeaseTerm),
		PetsAllowed:       types.BoolN(u.PetsAllowed),
		NumberOfResidents: types.Int32N(u.NumberOfResidents),
		ID:                id,
	}
}
