package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/require"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	client       *elasticsearch.Client
	typedClient  *elasticsearch.TypedClient
	dao          database.DAO
	retryBackoff = backoff.NewExponentialBackOff()
)

func TestMain(m *testing.M) {
	var err error

	cert, err := os.ReadFile("/media/pc/DATA/projects/rrms/rrms-backend/ca.crt")
	if err != nil {
		panic("Failed to read CA certificate")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "my12supers3cur3p4ssw0rd",
		CACert:   cert,

		// Retry on 429 TooManyRequests statuses
		//
		RetryOnStatus: []int{502, 503, 504, 429},

		// Configure the backoff function
		//
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	}
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	dao, err = database.NewPostgresDAO("postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable")
	if err != nil {
		panic(err)
	}

	esClient, err := NewElasticSearchClient(ElasticSearchClientParams{
		Addresses:  types.Ptr("https://localhost:9200"),
		Username:   types.Ptr("elastic"),
		Password:   types.Ptr("my12supers3cur3p4ssw0rd"),
		CACertPath: types.Ptr("/media/pc/DATA/projects/rrms/rrms-backend/ca.crt"),
	})
	if err != nil {
		panic(err)
	}
	typedClient = esClient.typedClient

	elasticsearch.NewDefaultClient()
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	// Send the request using the client
	res, err := client.Ping()
	require.NoError(t, err)
	defer res.Body.Close()

	// Print the response status and body
	require.False(t, res.IsError())
	t.Log(res.String())
}

func TestInfo(t *testing.T) {
	// Create the requesst
	req := esapi.InfoRequest{}

	// Send the request using the typed client
	res, err := req.Do(context.Background(), client.Transport)
	require.NoError(t, err)
	defer res.Body.Close()

	// Print the response status and body
	require.False(t, res.IsError())
	t.Log(res.String())
}

func TestTransferData_(t *testing.T) {
	// send first batch of 10 listings
	for o := 0; o <= 4001; o += 10 {

		listingsRows, err := dao.Query(
			context.Background(),
			fmt.Sprintf("SELECT id, creator_id, property_id, title, description, full_name, email, phone, contact_type, price, price_negotiable, security_deposit, lease_term, pets_allowed, number_of_residents, priority, active, created_at, updated_at, expired_at FROM listings LIMIT %d OFFSET %d", 10, o),
		)
		require.NoError(t, err)
		defer listingsRows.Close()
		for listingsRows.Next() {
			var listing database.Listing
			err = listingsRows.Scan(
				&listing.ID,
				&listing.CreatorID,
				&listing.PropertyID,
				&listing.Title,
				&listing.Description,
				&listing.FullName,
				&listing.Email,
				&listing.Phone,
				&listing.ContactType,
				&listing.Price,
				&listing.PriceNegotiable,
				&listing.SecurityDeposit,
				&listing.LeaseTerm,
				&listing.PetsAllowed,
				&listing.NumberOfResidents,
				&listing.Priority,
				&listing.Active,
				&listing.CreatedAt,
				&listing.UpdatedAt,
				&listing.ExpiredAt,
			)
			require.NoError(t, err)
			plisting := listing_model.ToListingModel(&listing)
			ludbs, err := dao.GetListingUnits(context.Background(), plisting.ID)
			require.NoError(t, err)
			for _, ludb := range ludbs {
				plisting.Units = append(plisting.Units, listing_model.ListingUnitModel(ludb))
			}

			propertyDB, err := dao.GetPropertyById(context.Background(), plisting.PropertyID)
			require.NoError(t, err)
			property := property_model.ToPropertyModel(&propertyDB)
			features, err := dao.GetPropertyFeatures(context.Background(), plisting.PropertyID)
			require.NoError(t, err)
			for _, f := range features {
				property.Features = append(property.Features, property_model.ToPropertyFeatureModel(&f))
			}

			units := make([]unit_model.UnitModel, 0, len(plisting.Units))
			for _, u := range plisting.Units {
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

			listingDoc := toListingDocument(plisting)
			propertyDoc := toPropertyDocument(property)
			unitDocs := toUnitDocument(units)
			// doc to bytes reader
			listingDocByte, err := json.Marshal(listingDoc)
			require.NoError(t, err)
			propertyDocByte, err := json.Marshal(propertyDoc)
			require.NoError(t, err)
			client.Index("listings", bytes.NewReader(listingDocByte))
			client.Index("properties", bytes.NewReader(propertyDocByte))
			for _, u := range unitDocs {
				unitDocByte, err := json.Marshal(u)
				require.NoError(t, err)
				client.Index("units", bytes.NewReader(unitDocByte))
			}

		}
		// sleep for 50ms
		time.Sleep(50 * time.Millisecond)
	}
}

func toListingDocument(listing *listing_model.ListingModel) map[string]interface{} {
	return map[string]interface{}{
		"id":                  listing.ID,
		"creator_id":          listing.CreatorID,
		"title":               listing.Title,
		"description":         listing.Description,
		"full_name":           listing.FullName,
		"email":               listing.Email,
		"phone":               listing.Phone,
		"contact_type":        listing.ContactType,
		"price":               listing.Price,
		"price_negotiable":    listing.PriceNegotiable,
		"security_deposit":    listing.SecurityDeposit,
		"lease_term":          listing.LeaseTerm,
		"pets_allowed":        listing.PetsAllowed,
		"number_of_residents": listing.NumberOfResidents,
		"priority":            listing.Priority,
		"active":              listing.Active,
		"created_at":          listing.CreatedAt,
		"updated_at":          listing.UpdatedAt,
		"expired_at":          listing.ExpiredAt,
	}
}

func toPropertyDocument(property *property_model.PropertyModel) map[string]interface{} {
	features := []int64{}
	for _, f := range property.Features {
		features = append(features, f.FeatureID)
	}

	return map[string]interface{}{
		"id":               property.ID,
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
		"features":         features,
		"created_at":       property.CreatedAt,
		"updated_at":       property.UpdatedAt,
	}
}

func toUnitDocument(units []unit_model.UnitModel) []map[string]interface{} {
	unitDocs := []map[string]interface{}{}
	for _, unit := range units {
		amenities := []int64{}
		for _, amenity := range unit.Amenities {
			amenities = append(amenities, amenity.AmenityID)
		}

		unitDocs = append(unitDocs, map[string]interface{}{
			"id":                     unit.ID,
			"name":                   unit.Name,
			"area":                   unit.Area,
			"floor":                  unit.Floor,
			"number_of_living_rooms": unit.NumberOfLivingRooms,
			"number_of_bedrooms":     unit.NumberOfBedrooms,
			"number_of_bathrooms":    unit.NumberOfBathrooms,
			"number_of_toilets":      unit.NumberOfToilets,
			"number_of_kitchens":     unit.NumberOfKitchens,
			"number_of_balconies":    unit.NumberOfBalconies,
			"type":                   unit.Type,
			"amenities":              amenities,
			"created_at":             unit.CreatedAt,
			"updated_at":             unit.UpdatedAt,
		})
	}
	return unitDocs
}

func TestCreateMappingIndices(t *testing.T) {
	CreateMappingIndices(client)
}
