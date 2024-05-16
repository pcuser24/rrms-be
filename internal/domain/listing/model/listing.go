package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type ListingPolicyModel struct {
	ListingID uuid.UUID `json:"listingId"`
	PolicyID  int64     `json:"policyId"`
	Note      *string   `json:"note"`
}

func ToListingPolicyModel(lp *database.ListingPolicy) ListingPolicyModel {
	lm := ListingPolicyModel{
		ListingID: lp.ListingID,
		PolicyID:  lp.PolicyID,
	}

	if lp.Note.Valid {
		val := lp.Note.String
		lm.Note = &val
	}

	return lm
}

type ListingUnitModel struct {
	ListingID uuid.UUID `json:"listingId"`
	UnitID    uuid.UUID `json:"unitId"`
	Price     int64     `json:"price"`
}

type ListingTagModel struct {
	ID        int64     `json:"id"`
	ListingID uuid.UUID `json:"listingId"`
	Tag       string    `json:"tag"`
}

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
	PriceNegotiable bool   `json:"priceNegotiable"`
	SecurityDeposit *int64 `json:"securityDeposit"`

	LeaseTerm         *int32 `json:"leaseTerm"`
	PetsAllowed       *bool  `json:"petsAllowed"`
	NumberOfResidents *int32 `json:"numberOfResidents"`

	Priority  int32                `json:"priority"`
	Active    bool                 `json:"active"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	ExpiredAt time.Time            `json:"expiredAt"`
	Policies  []ListingPolicyModel `json:"policies"`
	Units     []ListingUnitModel   `json:"units"`
	Tags      []ListingTagModel    `json:"tags"`
}

func ToListingModel(ldb *database.Listing) *ListingModel {
	lm := &ListingModel{
		ID:                ldb.ID,
		CreatorID:         ldb.CreatorID,
		PropertyID:        ldb.PropertyID,
		Title:             ldb.Title,
		Description:       ldb.Description,
		FullName:          ldb.FullName,
		Email:             ldb.Email,
		Phone:             ldb.Phone,
		ContactType:       ldb.ContactType,
		Price:             ldb.Price,
		PriceNegotiable:   ldb.PriceNegotiable,
		Priority:          ldb.Priority,
		Active:            ldb.Active,
		CreatedAt:         ldb.CreatedAt,
		UpdatedAt:         ldb.UpdatedAt,
		ExpiredAt:         ldb.ExpiredAt,
		SecurityDeposit:   types.PNInt64(ldb.SecurityDeposit),
		LeaseTerm:         types.PNInt32(ldb.LeaseTerm),
		PetsAllowed:       types.PNBool(ldb.PetsAllowed),
		NumberOfResidents: types.PNInt32(ldb.NumberOfResidents),
		Policies:          make([]ListingPolicyModel, 0),
		Units:             make([]ListingUnitModel, 0),
		Tags:              make([]ListingTagModel, 0),
	}

	return lm
}
