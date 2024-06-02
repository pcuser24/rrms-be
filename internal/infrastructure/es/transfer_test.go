package es

import (
	"bytes"
	"context"
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
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

func BuildAggregatedIndex(listing *listing_model.ListingModel, property *property_model.PropertyModel, units []unit_model.UnitModel) ([]byte, error) {
	aggregated := AggregatedIndex{
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
		Property:          convertProperty(property),
	}
	return json.Marshal(aggregated)
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

func convertProperty(property *property_model.PropertyModel) map[string]interface{} {
	return map[string]interface{}{
		"id":               property.ID.String(),
		"creator_id":       property.CreatorID.String(),
		"name":             property.Name,
		"building":         property.Building,
		"project":          property.Project,
		"area":             property.Area,
		"number_of_floors": property.NumberOfFloors,
		"year_built":       property.YearBuilt,
		"orientation":      property.Orientation,
		"entrance_width":   property.EntranceWidth,
		"facade":           property.Facade,
		"full_address":     property.FullAddress,
		"city":             property.City,
		"district":         property.District,
		"ward":             property.Ward,
		"lat":              property.Lat,
		"lng":              property.Lng,
		"primary_image":    property.PrimaryImage,
		"description":      property.Description,
		"type":             property.Type,
		"is_public":        property.IsPublic,
		"created_at":       property.CreatedAt,
		"updated_at":       property.UpdatedAt,
		"features":         convertFeatures(property.Features),
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

func TestTransferData(t *testing.T) {
	const batchSize = 16
	var numWorkers = runtime.NumCPU()

	for offset := 0; ; offset += batchSize {
		// bulk indexer
		bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Index:         "listings",       // The default index name
			Client:        client,           // The Elasticsearch client
			NumWorkers:    numWorkers,       // The number of worker goroutines
			FlushBytes:    6e+7,             // The flush threshold in bytes
			FlushInterval: 30 * time.Second, // The periodic flush interval
		})
		require.NoError(t, err)

		listingsDB, err := dao.GetSomeListings(context.Background(), database.GetSomeListingsParams{
			Limit:  batchSize,
			Offset: int32(offset),
		})
		require.NoError(t, err)
		if len(listingsDB) == 0 {
			break
		}

		start := time.Now().UTC()

		for _, ldb := range listingsDB {
			listing := listing_model.ToListingModel(&ldb)
			ludbs, err := dao.GetListingUnits(context.Background(), listing.ID)
			require.NoError(t, err)
			for _, ludb := range ludbs {
				listing.Units = append(listing.Units, listing_model.ListingUnitModel(ludb))
			}

			propertyDB, err := dao.GetPropertyById(context.Background(), listing.PropertyID)
			require.NoError(t, err)
			property := property_model.ToPropertyModel(&propertyDB)
			features, err := dao.GetPropertyFeatures(context.Background(), listing.PropertyID)
			require.NoError(t, err)
			for _, f := range features {
				property.Features = append(property.Features, property_model.ToPropertyFeatureModel(&f))
			}

			units := make([]unit_model.UnitModel, 0, len(listing.Units))
			for _, u := range listing.Units {
				unitDB, err := dao.GetUnitById(context.Background(), u.UnitID)
				require.NoError(t, err)
				unit := unit_model.ToUnitModel(&unitDB)
				amenities, err := dao.GetUnitAmenities(context.Background(), u.UnitID)
				require.NoError(t, err)
				for _, a := range amenities {
					unit.Amenities = append(unit.Amenities, *unit_model.ToUnitAmenityModel(&a))
				}
				units = append(units, *unit)
			}

			docByte, err := BuildAggregatedIndex(listing, property, units)
			require.NoError(t, err)

			err = bi.Add(
				context.Background(),
				esutil.BulkIndexerItem{
					// Action field configures the operation to perform (index, create, delete, update)
					Action: "index",

					// DocumentID is the (optional) document ID
					DocumentID: listing.ID.String(),

					// Body is an `io.Reader` with the payload
					Body: bytes.NewReader(docByte),

					// OnSuccess is called for each successful operation
					OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					},

					// OnFailure is called for each failed operation
					OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
						if err != nil {
							t.Logf("ERROR: %s", err)
						} else {
							t.Logf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
						}
					},
				},
			)
			require.NoError(t, err)
		}

		// Close the indexer
		//
		err = bi.Close(context.Background())
		require.NoError(t, err)
		biStats := bi.Stats()

		dur := time.Since(start)

		if biStats.NumFailed > 0 {
			t.Logf(
				"Indexed [%d] documents with [%d] errors in %s (%d docs/sec)",
				biStats.NumFlushed,
				biStats.NumFailed,
				dur.Truncate(time.Millisecond),
				int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
			)
		} else {
			t.Logf(
				"Sucessfuly indexed [%d] documents in %d (%d docs/sec)",
				int64(biStats.NumFlushed),
				dur.Truncate(time.Millisecond),
				int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
			)
		}

		if len(listingsDB) < batchSize {
			break
		}
	}
}
