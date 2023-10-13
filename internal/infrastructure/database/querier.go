// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	ChangePropertyVisibility(ctx context.Context, arg ChangePropertyVisibilityParams) error
	CheckListingOwnership(ctx context.Context, arg CheckListingOwnershipParams) (int64, error)
	CheckUnitManageability(ctx context.Context, arg CheckUnitManageabilityParams) (int64, error)
	CheckUnitOfProperty(ctx context.Context, arg CheckUnitOfPropertyParams) (int64, error)
	CheckValidUnitForListing(ctx context.Context, arg CheckValidUnitForListingParams) (int64, error)
	CreateListing(ctx context.Context, arg CreateListingParams) (Listing, error)
	CreateListingPolicy(ctx context.Context, arg CreateListingPolicyParams) (ListingPolicy, error)
	CreateListingUnit(ctx context.Context, arg CreateListingUnitParams) (ListingUnit, error)
	CreateProperty(ctx context.Context, arg CreatePropertyParams) (Property, error)
	CreateUnit(ctx context.Context, arg CreateUnitParams) (Unit, error)
	DeleteListing(ctx context.Context, id uuid.UUID) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	DeletePropertyFeature(ctx context.Context, arg DeletePropertyFeatureParams) error
	DeletePropertyManager(ctx context.Context, arg DeletePropertyManagerParams) error
	DeletePropertyMedia(ctx context.Context, arg DeletePropertyMediaParams) error
	DeletePropertyTag(ctx context.Context, arg DeletePropertyTagParams) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
	DeleteUnitAmenity(ctx context.Context, arg DeleteUnitAmenityParams) error
	DeleteUnitMedia(ctx context.Context, arg DeleteUnitMediaParams) error
	GetAllPropertyFeatures(ctx context.Context) ([]PFeature, error)
	GetAllRentalPolicies(ctx context.Context) ([]RentalPolicy, error)
	GetAllUnitAmenities(ctx context.Context) ([]UAmenity, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (Listing, error)
	GetListingPolicies(ctx context.Context, listingID uuid.UUID) ([]ListingPolicy, error)
	GetListingUnits(ctx context.Context, listingID uuid.UUID) ([]ListingUnit, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (Property, error)
	GetPropertyFeatures(ctx context.Context, propertyID uuid.UUID) ([]PropertyFeature, error)
	GetPropertyManagers(ctx context.Context, propertyID uuid.UUID) ([]PropertyManager, error)
	GetPropertyMedia(ctx context.Context, propertyID uuid.UUID) ([]PropertyMedia, error)
	GetPropertyTags(ctx context.Context, propertyID uuid.UUID) ([]PropertyTag, error)
	GetUnitAmenities(ctx context.Context, unitID uuid.UUID) ([]UnitAmenity, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (Unit, error)
	GetUnitManagers(ctx context.Context, id uuid.UUID) ([]PropertyManager, error)
	GetUnitMedia(ctx context.Context, unitID uuid.UUID) ([]UnitMedia, error)
	GetUnitsOfProperty(ctx context.Context, propertyID uuid.UUID) ([]Unit, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	InsertPropertyFeature(ctx context.Context, arg InsertPropertyFeatureParams) (PropertyFeature, error)
	InsertPropertyManager(ctx context.Context, arg InsertPropertyManagerParams) error
	InsertPropertyMedia(ctx context.Context, arg InsertPropertyMediaParams) (PropertyMedia, error)
	InsertPropertyTag(ctx context.Context, arg InsertPropertyTagParams) (PropertyTag, error)
	InsertUnitAmenity(ctx context.Context, arg InsertUnitAmenityParams) (UnitAmenity, error)
	InsertUnitMedia(ctx context.Context, arg InsertUnitMediaParams) (UnitMedia, error)
	InsertUser(ctx context.Context, arg InsertUserParams) (User, error)
	IsPropertyPublic(ctx context.Context, id uuid.UUID) (bool, error)
	IsUnitPublic(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateListing(ctx context.Context, arg UpdateListingParams) error
	UpdateListingStatus(ctx context.Context, arg UpdateListingStatusParams) error
	UpdateProperty(ctx context.Context, arg UpdatePropertyParams) error
	UpdatePropertyManager(ctx context.Context, arg UpdatePropertyManagerParams) error
	UpdateUnit(ctx context.Context, arg UpdateUnitParams) error
}

var _ Querier = (*Queries)(nil)