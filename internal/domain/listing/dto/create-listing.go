package dto

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type CreateListingPolicy struct {
	PolicyID int64   `json:"policyId" validate:"required"`
	Note     *string `json:"note"`
}

type CreateListingUnit struct {
	UnitID uuid.UUID `json:"unitId" validate:"required,uuid4"`
}

type CreateListing struct {
	CreatorID         uuid.UUID             `json:"creatorId"`
	PropertyID        uuid.UUID             `json:"propertyId" validate:"required,uuid4"`
	Title             string                `json:"title" validate:"required"`
	Description       string                `json:"description" validate:"required"`
	FullName          string                `json:"fullName" validate:"required"`
	Email             string                `json:"email" validate:"required,email"`
	Phone             string                `json:"phone" validate:"required"`
	ContactType       string                `json:"contactType" validate:"required"`
	Price             int64                 `json:"price" validate:"required,gt=0"`
	PriceNegotiable   bool                  `json:"priceNegotiable"`
	SecurityDeposit   *int64                `json:"securityDeposit" validate:"omitempty,gt=0"`
	LeaseTerm         *int32                `json:"leaseTerm" validate:"required,gt=0"`
	PetsAllowed       *bool                 `json:"petsAllowed"`
	NumberOfResidents *int32                `json:"numberOfResidents" validate:"omitempty,gt=0"`
	Priority          int32                 `json:"priority" validate:"required,gte=1,lte=5"`
	PostAt            time.Time             `json:"postAt" validate:"required"`
	Active            bool                  `json:"active" validate:"required"`
	PostDuration      int                   `json:"postDuration" validate:"required"`
	Policies          []CreateListingPolicy `json:"policies" validate:"dive"`
	Units             []CreateListingUnit   `json:"units" validate:"dive"`
}

func (c *CreateListing) ToCreateListingDB() *database.CreateListingParams {
	ldb := &database.CreateListingParams{
		CreatorID:         c.CreatorID,
		PropertyID:        c.PropertyID,
		Title:             c.Title,
		Description:       c.Description,
		FullName:          c.FullName,
		Email:             c.Email,
		Phone:             c.Phone,
		ContactType:       c.ContactType,
		Price:             c.Price,
		PriceNegotiable:   sql.NullBool{Valid: true, Bool: c.PriceNegotiable},
		SecurityDeposit:   types.Int64N(c.SecurityDeposit),
		LeaseTerm:         types.Int32N(c.LeaseTerm),
		PetsAllowed:       types.BoolN(c.PetsAllowed),
		NumberOfResidents: types.Int32N(c.NumberOfResidents),
		Priority:          c.Priority,
		PostAt:            c.PostAt,
		Active:            c.Active,
	}
	// add duration to postAt to get expiry
	ldb.ExpiredAt = c.PostAt.AddDate(0, 0, c.PostDuration)
	// if the listing is not active but postAt is in the past, make it active
	if !c.Active && c.PostAt.Before(time.Now()) {
		ldb.Active = true
	}
	return ldb
}
