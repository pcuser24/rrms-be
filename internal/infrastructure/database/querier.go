// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CheckListingOwnership(ctx context.Context, arg CheckListingOwnershipParams) (int64, error)
	CheckPropertyOwnerShip(ctx context.Context, arg CheckPropertyOwnerShipParams) (int64, error)
	CheckUnitOfProperty(ctx context.Context, arg CheckUnitOfPropertyParams) (int64, error)
	CheckUnitOwnership(ctx context.Context, arg CheckUnitOwnershipParams) (int64, error)
	CheckValidUnitForListing(ctx context.Context, arg CheckValidUnitForListingParams) (int64, error)
	CreateListing(ctx context.Context, arg CreateListingParams) (Listing, error)
	CreateListingPolicy(ctx context.Context, arg CreateListingPolicyParams) (ListingPolicy, error)
	CreateListingUnit(ctx context.Context, arg CreateListingUnitParams) (ListingUnit, error)
	CreateProperty(ctx context.Context, arg CreatePropertyParams) (Property, error)
	CreatePropertyAmenity(ctx context.Context, arg CreatePropertyAmenityParams) (PropertyAmenity, error)
	CreatePropertyFeature(ctx context.Context, arg CreatePropertyFeatureParams) (PropertyFeature, error)
	CreatePropertyMedia(ctx context.Context, arg CreatePropertyMediaParams) (PropertyMedium, error)
	CreateUnit(ctx context.Context, arg CreateUnitParams) (Unit, error)
	CreateUnitAmenity(ctx context.Context, arg CreateUnitAmenityParams) (UnitAmenity, error)
	CreateUnitMedia(ctx context.Context, arg CreateUnitMediaParams) (UnitMedium, error)
	DeleteAllPropertyAmenity(ctx context.Context, propertyID uuid.UUID) error
	DeleteAllPropertyFeature(ctx context.Context, propertyID uuid.UUID) error
	DeleteAllPropertyMedia(ctx context.Context, propertyID uuid.UUID) error
	DeleteAllPropertyTag(ctx context.Context, propertyID uuid.UUID) error
	DeleteAllUnitAmenity(ctx context.Context, unitID uuid.UUID) error
	DeleteAllUnitMedia(ctx context.Context, unitID uuid.UUID) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
	GetAllPropertyAmenities(ctx context.Context) ([]PAmenity, error)
	GetAllPropertyFeatures(ctx context.Context) ([]PFeature, error)
	GetAllRentalPolicies(ctx context.Context) ([]RentalPolicy, error)
	GetAllUnitAmenities(ctx context.Context) ([]UAmenity, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (Listing, error)
	GetListingPolicies(ctx context.Context, listingID uuid.UUID) ([]ListingPolicy, error)
	GetListingUnits(ctx context.Context, listingID uuid.UUID) ([]ListingUnit, error)
	GetPropertiesByOwnerId(ctx context.Context, arg GetPropertiesByOwnerIdParams) ([]Property, error)
	GetPropertyAmenities(ctx context.Context, propertyID uuid.UUID) ([]PropertyAmenity, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (Property, error)
	GetPropertyFeatures(ctx context.Context, propertyID uuid.UUID) ([]PropertyFeature, error)
	GetPropertyMedium(ctx context.Context, propertyID uuid.UUID) ([]PropertyMedium, error)
	GetPropertyTags(ctx context.Context, propertyID uuid.UUID) ([]PropertyTag, error)
	GetUnitAmenities(ctx context.Context, unitID uuid.UUID) ([]UnitAmenity, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (Unit, error)
	GetUnitMedia(ctx context.Context, unitID uuid.UUID) ([]UnitMedium, error)
	GetUnitsOfProperty(ctx context.Context, propertyID uuid.UUID) ([]Unit, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	InsertUser(ctx context.Context, arg InsertUserParams) (User, error)
	UpdateListing(ctx context.Context, arg UpdateListingParams) error
	UpdateProperty(ctx context.Context, arg UpdatePropertyParams) error
	UpdateUnit(ctx context.Context, arg UpdateUnitParams) error
}

var _ Querier = (*Queries)(nil)
