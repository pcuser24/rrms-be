package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateListing struct {
	Title             *string   `json:"title"`
	Description       *string   `json:"description"`
	Price             *int64    `json:"price"`
	SecurityDeposit   *int64    `json:"securityDeposit"`
	LeaseTerm         *int32    `json:"leaseTerm"`
	PetsAllowed       *bool     `json:"petsAllowed"`
	NumberOfResidents *int32    `json:"numberOfResidents"`
	ID                uuid.UUID `json:"id"`
}

func (u *UpdateListing) ToUpdateListingDB() *database.UpdateListingParams {
	return &database.UpdateListingParams{
		Title:             types.StrN(u.Title),
		Description:       types.StrN(u.Description),
		Price:             types.Int64N(u.Price),
		SecurityDeposit:   types.Int64N(u.SecurityDeposit),
		LeaseTerm:         types.Int32N(u.LeaseTerm),
		PetsAllowed:       types.BoolN(u.PetsAllowed),
		NumberOfResidents: types.Int32N(u.NumberOfResidents),
		ID:                u.ID,
	}
}
