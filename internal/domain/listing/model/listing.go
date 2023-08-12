package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type ListingModel struct {
	ID          uuid.UUID `json:"id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	PropertyID  uuid.UUID `json:"property_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	// Rental price per month in vietnamese dong
	Price           int64 `json:"price"`
	SecurityDeposit int64 `json:"security_deposit"`
	// Lease term in months
	LeaseTerm         int32  `json:"lease_term"`
	PetsAllowed       *bool  `json:"pets_allowed"`
	NumberOfResidents *int32 `json:"number_of_residents"`
	// Priority of the listing, range from 1 to 5, 1 is the lowest
	Priority  int32                `json:"priority"`
	Active    bool                 `json:"active"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	ExpiredAt time.Time            `json:"expired_at"`
	Policies  []ListingPolicyModel `json:"policies"`
	Units     []ListingUnitModel   `json:"units"`
}

func ToListingModel(ldb *database.Listing) *ListingModel {
	lm := &ListingModel{
		ID:              ldb.ID,
		CreatorID:       ldb.CreatorID,
		PropertyID:      ldb.PropertyID,
		Title:           ldb.Title,
		Description:     ldb.Description,
		Price:           ldb.Price,
		SecurityDeposit: ldb.SecurityDeposit,
		LeaseTerm:       ldb.LeaseTerm,
		Priority:        ldb.Priority,
		Active:          ldb.Active,
		CreatedAt:       ldb.CreatedAt,
		UpdatedAt:       ldb.UpdatedAt,
		ExpiredAt:       ldb.ExpiredAt,
	}

	if ldb.PetsAllowed.Valid {
		val := ldb.PetsAllowed.Bool
		lm.PetsAllowed = &val
	}

	if ldb.NumberOfResidents.Valid {
		val := ldb.NumberOfResidents.Int32
		lm.NumberOfResidents = &val
	}

	return lm
}
