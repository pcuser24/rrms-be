package listing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/huandu/go-sqlbuilder"
	listingDTO "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	propertyDTO "github.com/user2410/rrms-backend/internal/domain/property/dto"
	unitDTO "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestSearchListing(t *testing.T) {
	sql, args := sqlbuilders.SearchListingBuilder([]string{"1"}, &listingDTO.SearchListingQuery{
		LTitle:                types.Ptr[string]("test_name"),
		LCreatorID:            types.Ptr[string]("5224e97d-94dd-4988-98fb-ab5d4a3e0967"),
		LPropertyID:           types.Ptr[string]("5224e97d-94dd-4988-98fb-ab5d4a3e0967"),
		LMinPrice:             types.Ptr[int64](3000),
		LMaxPrice:             types.Ptr[int64](2000),
		LPriceNegotiable:      types.Ptr[bool](true),
		LSecurityDeposit:      types.Ptr[int64](3000),
		LLeaseTerm:            types.Ptr[int32](12),
		LPetsAllowed:          types.Ptr[bool](true),
		LMinNumberOfResidents: types.Ptr[int32](12),
		LPriority:             types.Ptr[int32](12),
		LActive:               types.Ptr[bool](true),
		LPolicies:             []int32{1, 2, 3},
		LMinCreatedAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxCreatedAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinUpdatedAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxUpdatedAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinPostAt:            types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxPostAt:            types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinExpiredAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxExpiredAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
	}, "", "", "float32")
	t.Log(sql)
	t.Log(args)
}

func TestSearchListingCombination(t *testing.T) {
	sqlListing, argsListing := sqlbuilders.SearchListingBuilder([]string{"1"}, &listingDTO.SearchListingQuery{
		LTitle:                types.Ptr[string]("test_name"),
		LCreatorID:            types.Ptr[string]("5224e97d-94dd-4988-98fb-ab5d4a3e0967"),
		LPropertyID:           types.Ptr[string]("5224e97d-94dd-4988-98fb-ab5d4a3e0967"),
		LMinPrice:             types.Ptr[int64](3000),
		LMaxPrice:             types.Ptr[int64](2000),
		LPriceNegotiable:      types.Ptr[bool](true),
		LSecurityDeposit:      types.Ptr[int64](3000),
		LLeaseTerm:            types.Ptr[int32](12),
		LPetsAllowed:          types.Ptr[bool](true),
		LMinNumberOfResidents: types.Ptr[int32](12),
		LPriority:             types.Ptr[int32](12),
		LActive:               types.Ptr[bool](true),
		LPolicies:             []int32{1, 2, 3},
		LMinCreatedAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxCreatedAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinUpdatedAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxUpdatedAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinPostAt:            types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxPostAt:            types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMinExpiredAt:         types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		LMaxExpiredAt:         types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
	}, "", "", "")

	sqlUnit, argsUnit := sqlbuilders.SearchUnitBuilder([]string{"1"}, &unitDTO.SearchUnitQuery{
		UPropertyID:          types.Ptr[string]("5224e97d-94dd-4988-98fb-ab5d4a3e0967"),
		UName:                types.Ptr[string]("test_name"),
		UMinArea:             types.Ptr[int64](3000),
		UMaxArea:             types.Ptr[int64](2000),
		UFloor:               types.Ptr[int32](12),
		UNumberOfLivingRooms: types.Ptr[int32](12),
		UNumberOfBedrooms:    types.Ptr[int32](12),
		UNumberOfBathrooms:   types.Ptr[int32](12),
		UNumberOfToilets:     types.Ptr[int32](12),
		UNumberOfKitchens:    types.Ptr[int32](12),
		UNumberOfBalconies:   types.Ptr[int32](12),
	}, "", "")

	sqlProp, argsProp := sqlbuilders.SearchPropertyBuilder([]string{"1"}, &propertyDTO.SearchPropertyQuery{
		PTypes:          []string{"APARTMENT", "HOUSE"},
		PName:           types.Ptr[string]("test_name"),
		PBuilding:       types.Ptr[string]("test_building"),
		PProject:        types.Ptr[string]("test_project"),
		PFullAddress:    types.Ptr[string]("test_address"),
		PCity:           types.Ptr[string]("test_city"),
		PDistrict:       types.Ptr[string]("test_district"),
		PWard:           types.Ptr[string]("test_ward"),
		PMinArea:        types.Ptr[float32](3000),
		PMaxArea:        types.Ptr[float32](2000),
		PNumberOfFloors: types.Ptr[int32](12),
		PYearBuilt:      types.Ptr[int32](2023),
		POrientation:    types.Ptr[string]("nw"),
		PMinFacade:      types.Ptr[int32](12),
		PIsPublic:       types.Ptr[bool](true),
		PFeatures:       []int32{1, 2, 3},
		PTags:           []string{"tag 1", "tag 2"},
		PMinCreatedAt:   types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		PMaxCreatedAt:   types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
	}, "", "")

	// build order: unit -> property -> listing
	var queryStr string = sqlListing
	var argsLs []interface{} = argsListing
	t.Log("argsListing", argsListing)
	if len(sqlProp) > 0 {
		var tmp string = sqlProp
		t.Log("argsProp", argsProp)
		argsLs = append(argsLs, argsProp...)
		if len(sqlUnit) > 0 {
			tmp += fmt.Sprintf(" AND EXISTS (%v) ", sqlUnit)
			argsLs = append(argsLs, argsUnit...)
			t.Log("argsUnit", argsUnit)
		}
		queryStr = fmt.Sprintf("%v AND EXISTS (%v) ", sqlListing, tmp)
	}

	sql, args := sqlbuilder.Build(queryStr, argsLs...).Build()
	t.Log(sql)
	t.Log(args)

	dao, err := database.NewDAO("postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer dao.Close()

	sqSql := utils.SequelizePlaceholders(sql)
	t.Log("sqSql", sqSql)
	rows, err := dao.Query(context.Background(), sqSql, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	// print rows
	for rows.Next() {
		var id, creator_id, property_id, title, description, address, city, district, ward, orientation, status string
		var price, security_deposit, lease_term, number_of_residents, priority int64
		var price_negotiable, pets_allowed, active bool
		var created_at, updated_at, post_at, expired_at time.Time
		if err := rows.Scan(
			&id,
			&creator_id,
			&property_id,
			&title,
			&description,
			&address,
			&city,
			&district,
			&ward,
			&orientation,
			&price,
			&price_negotiable,
			&security_deposit,
			&lease_term,
			&pets_allowed,
			&number_of_residents,
			&priority,
			&active,
			&created_at,
			&updated_at,
			&post_at,
			&expired_at,
			&status,
		); err != nil {
			t.Fatal(err)
		}
		t.Log(id, creator_id, property_id, title, description, address, city, district, ward, orientation, price, price_negotiable, security_deposit, lease_term, pets_allowed, number_of_residents, priority, active, created_at, updated_at, post_at, expired_at, status)
	}
}
