package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type SuggestedListingUnit struct {
	Amenities []struct {
		AmenityID   int64   `json:"amenity_id"`
		Description *string `json:"description"`
	} `json:"amenities"`
	Area                int64             `json:"area"`
	CreatedAt           time.Time         `json:"created_at"`
	Floor               *int32            `json:"floor"`
	Name                string            `json:"name"`
	NumberOfBalconies   *int32            `json:"number_of_balconies"`
	NumberOfBathrooms   *int32            `json:"number_of_bathrooms"`
	NumberOfBedrooms    *int32            `json:"number_of_bedrooms"`
	NumberOfKitchens    *int32            `json:"number_of_kitchens"`
	NumberOfLivingRooms *int32            `json:"number_of_living_rooms"`
	NumberOfToilets     *int32            `json:"number_of_toilets"`
	Price               float32           `json:"price"`
	Type                database.UNITTYPE `json:"type"`
	UnitID              string            `json:"unit_id"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

type SuggestedProperty struct {
	ID             uuid.UUID `json:"id"`
	CreatorID      uuid.UUID `json:"creator_id"`
	Name           string    `json:"name"`
	Building       *string   `json:"building"`
	Project        *string   `json:"project"`
	Area           float32   `json:"area"`
	NumberOfFloors *int32    `json:"number_of_floors"`
	YearBuilt      *int32    `json:"year_built"`
	// n,s,w,e,nw,ne,sw,se
	Orientation   *string               `json:"orientation"`
	EntranceWidth *float32              `json:"entrance_width"`
	Facade        *float32              `json:"facade"`
	FullAddress   string                `json:"full_address"`
	District      string                `json:"district"`
	City          string                `json:"city"`
	Ward          *string               `json:"ward"`
	Lat           *float64              `json:"lat"`
	Lng           *float64              `json:"lng"`
	PrimaryImage  int64                 `json:"primary_image"`
	Description   *string               `json:"description"`
	Type          database.PROPERTYTYPE `json:"type"`
	IsPublic      bool                  `json:"is_public"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Features      []struct {
		FeatureID   int64   `json:"feature_id"`
		Description *string `json:"description"`
	} `json:"features"`
}

type SuggestedListing struct {
	ID                string                 `json:"id"`
	CreatorID         string                 `json:"creator_id"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	FullName          string                 `json:"full_name"`
	Email             string                 `json:"email"`
	Phone             string                 `json:"phone"`
	ContactType       string                 `json:"contact_type"`
	Price             float32                `json:"price"`
	PriceNegotiable   bool                   `json:"price_negotiable"`
	SecurityDeposit   *float32               `json:"security_deposit"`
	LeaseTerm         *int32                 `json:"lease_term"`
	PetsAllowed       *bool                  `json:"pets_allowed"`
	NumberOfResidents *int32                 `json:"number_of_residents"`
	Priority          int32                  `json:"priority"`
	Active            bool                   `json:"active"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	ExpiredAt         time.Time              `json:"expired_at"`
	Tags              []string               `json:"tags"`
	ListingUnits      []SuggestedListingUnit `json:"listing_units"`
	Property          SuggestedProperty      `json:"property"`
}

type ListingsSuggestionResult struct {
	Hits []SuggestedListing `json:"hits"`
}
