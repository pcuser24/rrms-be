package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type CreateListingPolicy struct {
	PolicyID int64   `json:"policy_id"`
	Note     *string `json:"note"`
}

type CreateListingUnit struct {
	UnitID uuid.UUID `json:"unit_id" validate:"required"`
}

type CreateListing struct {
	CreatorID         uuid.UUID             `json:"creator_id"`
	PropertyID        uuid.UUID             `json:"property_id" validate:"required"`
	Title             string                `json:"title" validate:"required"`
	Description       string                `json:"description" validate:"required"`
	Price             int64                 `json:"price" validate:"required,gt=0"`
	SecurityDeposit   *int64                `json:"security_deposit" validate:"gt=0"`
	LeaseTerm         int32                 `json:"lease_term" validate:"required,gt=0"`
	PetsAllowed       *bool                 `json:"pets_allowed"`
	NumberOfResidents *int32                `json:"number_of_residents" validate:"omitempty,gt=0"`
	Priority          int32                 `json:"priority" validate:"required,gte=1,lte=5"`
	ExpiredAt         time.Time             `json:"expired_at" validate:"required"`
	Policies          []CreateListingPolicy `json:"policies"`
	Units             []CreateListingUnit   `json:"units"`
}

func (c *CreateListing) ToCreateListingDB() *database.CreateListingParams {
	return &database.CreateListingParams{
		CreatorID:         c.CreatorID,
		PropertyID:        c.PropertyID,
		Title:             c.Title,
		Description:       c.Description,
		Price:             c.Price,
		SecurityDeposit:   types.Int64N(c.SecurityDeposit),
		LeaseTerm:         c.LeaseTerm,
		PetsAllowed:       types.BoolN(c.PetsAllowed),
		NumberOfResidents: types.Int32N(c.NumberOfResidents),
		Priority:          c.Priority,
		ExpiredAt:         c.ExpiredAt,
	}
}
