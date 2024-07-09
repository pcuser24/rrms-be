package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_utils "github.com/user2410/rrms-backend/internal/domain/listing/utils"
	payment_dto "github.com/user2410/rrms-backend/internal/domain/payment/dto"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
)

type AggregatedIndex struct {
	ID                string                   `json:"id"`
	CreatorID         string                   `json:"creator_id"`
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	FullName          string                   `json:"full_name"`
	Email             string                   `json:"email"`
	Phone             string                   `json:"phone"`
	ContactType       string                   `json:"contact_type"`
	Price             float32                  `json:"price"`
	PriceNegotiable   bool                     `json:"price_negotiable"`
	SecurityDeposit   *float32                 `json:"security_deposit"`
	LeaseTerm         *int32                   `json:"lease_term"`
	PetsAllowed       *bool                    `json:"pets_allowed"`
	NumberOfResidents *int32                   `json:"number_of_residents"`
	Priority          int32                    `json:"priority"`
	Active            bool                     `json:"active"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	ExpiredAt         time.Time                `json:"expired_at"`
	Tags              []map[string]string      `json:"tags"`
	ListingUnits      []map[string]interface{} `json:"listing_units"`
	Property          map[string]interface{}   `json:"property"`
}

func buildAggregatedIndex(listing *listing_model.ListingModel, property *property_model.PropertyModel, propertyVStatus any, units []unit_model.UnitModel) AggregatedIndex {
	return AggregatedIndex{
		ID:                listing.ID.String(),
		CreatorID:         listing.CreatorID.String(),
		Title:             listing.Title,
		Description:       listing.Description,
		FullName:          listing.FullName,
		Email:             listing.Email,
		Phone:             listing.Phone,
		ContactType:       listing.ContactType,
		Price:             listing.Price,
		PriceNegotiable:   listing.PriceNegotiable,
		SecurityDeposit:   listing.SecurityDeposit,
		LeaseTerm:         listing.LeaseTerm,
		PetsAllowed:       listing.PetsAllowed,
		NumberOfResidents: listing.NumberOfResidents,
		Priority:          listing.Priority,
		Active:            listing.Active,
		CreatedAt:         listing.CreatedAt,
		UpdatedAt:         listing.UpdatedAt,
		ExpiredAt:         listing.ExpiredAt,
		Tags:              convertTags(listing.Tags),
		ListingUnits:      convertListingUnits(units, listing.Units),
		Property:          convertProperty(property, propertyVStatus),
	}
}

func convertTags(tags []listing_model.ListingTagModel) []map[string]string {
	var result []map[string]string
	for _, tag := range tags {
		result = append(result, map[string]string{"tag": tag.Tag})
	}
	return result
}

func convertListingUnits(units []unit_model.UnitModel, listingUnits []listing_model.ListingUnitModel) []map[string]interface{} {
	var result []map[string]interface{}
	unitMap := make(map[uuid.UUID]unit_model.UnitModel)
	for _, unit := range units {
		unitMap[unit.ID] = unit
	}
	for _, lu := range listingUnits {
		unit := unitMap[lu.UnitID]
		unitData := map[string]interface{}{
			"unit_id":                lu.UnitID.String(),
			"price":                  lu.Price,
			"name":                   unit.Name,
			"area":                   unit.Area,
			"floor":                  unit.Floor,
			"number_of_living_rooms": unit.NumberOfLivingRooms,
			"number_of_bedrooms":     unit.NumberOfBedrooms,
			"number_of_bathrooms":    unit.NumberOfBathrooms,
			"number_of_toilets":      unit.NumberOfToilets,
			"number_of_balconies":    unit.NumberOfBalconies,
			"number_of_kitchens":     unit.NumberOfKitchens,
			"type":                   unit.Type,
			"created_at":             unit.CreatedAt,
			"updated_at":             unit.UpdatedAt,
			"amenities":              convertAmenities(unit.Amenities),
		}
		result = append(result, unitData)
	}
	return result
}

func convertAmenities(amenities []unit_model.UnitAmenityModel) []map[string]interface{} {
	var result []map[string]interface{}
	for _, amenity := range amenities {
		result = append(result, map[string]interface{}{
			"amenity_id":  amenity.AmenityID,
			"description": amenity.Description,
		})
	}
	return result
}

func convertProperty(property *property_model.PropertyModel, pv any) map[string]interface{} {
	return map[string]interface{}{
		"id":                  property.ID.String(),
		"creator_id":          property.CreatorID.String(),
		"name":                property.Name,
		"building":            property.Building,
		"project":             property.Project,
		"area":                property.Area,
		"number_of_floors":    property.NumberOfFloors,
		"year_built":          property.YearBuilt,
		"orientation":         property.Orientation,
		"entrance_width":      property.EntranceWidth,
		"facade":              property.Facade,
		"full_address":        property.FullAddress,
		"city":                property.City,
		"district":            property.District,
		"ward":                property.Ward,
		"lat":                 property.Lat,
		"lng":                 property.Lng,
		"primary_image":       property.PrimaryImage,
		"description":         property.Description,
		"type":                property.Type,
		"is_public":           property.IsPublic,
		"verification_status": pv,
		"created_at":          property.CreatedAt,
		"updated_at":          property.UpdatedAt,
		"features":            convertFeatures(property.Features),
	}
}

func convertFeatures(features []property_model.PropertyFeatureModel) []map[string]interface{} {
	var result []map[string]interface{}
	for _, feature := range features {
		result = append(result, map[string]interface{}{
			"feature_id":  feature.FeatureID,
			"description": feature.Description,
		})
	}
	return result
}

func (s *service) CreateListing(data *dto.CreateListing) (*dto.CreateListingResponse, error) {
	var (
		res = new(dto.CreateListingResponse)
		err error
	)
	// create listing
	res.Listing, err = s.domainRepo.ListingRepo.CreateListing(context.Background(), data)
	if err != nil {
		return nil, err
	}

	// create payment info
	params := payment_dto.CreatePayment{UserId: data.CreatorID}
	amount, price, discount, err := listing_utils.CalculateListingPrice(int(data.Priority), data.PostDuration)
	if err != nil {
		return nil, err
	}
	params.Amount = amount
	params.OrderInfo = fmt.Sprintf("[%s%s%s] Phi dang tin nha cho thue", payment_service.PAYMENTTYPE_CREATELISTING, payment_service.PAYMENTTYPE_DELIMITER, res.Listing.ID.String())
	params.Items = []payment_dto.CreatePaymentItem{
		{
			Name:     "Phi dang tin",
			Price:    float32(price),
			Quantity: int32(data.PostDuration),
			Discount: int32(discount),
		},
	}

	res.Payment, err = s.domainRepo.PaymentRepo.CreatePayment(context.Background(), &params)
	if err != nil {
		return nil, err
	}

	// index new listing
	esClient := s.esClient.GetTypedClient()

	property, err := s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), res.Listing.PropertyID)
	if err != nil {
		return res, err
	}
	var pv any = nil
	pvs, err := s.domainRepo.PropertyRepo.GetPropertiesVerificationStatus(context.Background(), []uuid.UUID{property.ID})
	if err != nil {
		return res, err
	}
	if len(pvs) > 0 {
		pv = pvs[0].Status
	}
	units := make([]unit_model.UnitModel, 0, len(res.Listing.Units))
	for _, u := range res.Listing.Units {
		unit, err := s.domainRepo.UnitRepo.GetUnitById(context.Background(), u.UnitID)
		if err != nil {
			return res, err
		}
		units = append(units, *unit)
	}
	doc := buildAggregatedIndex(res.Listing, property, pv, units)
	_, err = esClient.Index(string(es.LISTINGINDEX)).Request(doc).Id(res.Listing.ID.String()).Do(context.Background())

	return res, err
}
