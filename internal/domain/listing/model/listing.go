package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type ListingModel struct {
	ID          uuid.UUID `json:"id"`
	CreatorID   uuid.UUID `json:"creatorId"`
	PropertyID  uuid.UUID `json:"propertyId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`

	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ContactType string `json:"contactType"`

	Price           int64  `json:"price"`
	SecurityDeposit *int64 `json:"securityDeposit"`

	LeaseTerm         *int32 `json:"leaseTerm"`
	PetsAllowed       *bool  `json:"petsAllowed"`
	NumberOfResidents *int32 `json:"numberOfResidents"`

	Priority  int32                `json:"priority"`
	Active    bool                 `json:"active"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	ExpiredAt time.Time            `json:"expiredAt"`
	PostAt    time.Time            `json:"postAt"`
	Policies  []ListingPolicyModel `json:"policies"`
	Units     []ListingUnitModel   `json:"units"`
}

func ToListingModel(ldb *database.Listing) *ListingModel {
	lm := &ListingModel{
		ID:          ldb.ID,
		CreatorID:   ldb.CreatorID,
		PropertyID:  ldb.PropertyID,
		Title:       ldb.Title,
		Description: ldb.Description,
		FullName:    ldb.FullName,
		Email:       ldb.Email,
		Phone:       ldb.Phone,
		ContactType: ldb.ContactType,
		Price:       ldb.Price,
		Priority:    ldb.Priority,
		Active:      ldb.Active,
		CreatedAt:   ldb.CreatedAt,
		UpdatedAt:   ldb.UpdatedAt,
		ExpiredAt:   ldb.ExpiredAt,
		PostAt:      ldb.PostAt,
	}

	if ldb.SecurityDeposit.Valid {
		val := ldb.SecurityDeposit.Int64
		lm.SecurityDeposit = &val
	}

	if ldb.LeaseTerm.Valid {
		val := ldb.LeaseTerm.Int32
		lm.LeaseTerm = &val
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
