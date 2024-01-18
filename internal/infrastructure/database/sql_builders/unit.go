package sqlbuilders

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
)

func SearchUnitBuilder(
	searchFields []string, query *dto.SearchUnitQuery,
	connectID, connectProperty string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM units"
	var searchQueries []string
	var args []interface{}

	if query.UName != nil {
		searchQueries = append(searchQueries, "units.name ILIKE $?")
		args = append(args, "%"+(*query.UName)+"%")
	}
	if query.UPropertyID != nil {
		searchQueries = append(searchQueries, "units.property_id = $?")
		args = append(args, *query.UPropertyID)
	}
	if query.UMinArea != nil {
		searchQueries = append(searchQueries, "units.area >= $?")
		args = append(args, *query.UMinArea)
	}
	if query.UMaxArea != nil {
		searchQueries = append(searchQueries, "units.area <= $?")
		args = append(args, *query.UMaxArea)
	}
	if query.UFloor != nil {
		searchQueries = append(searchQueries, "units.floor = $?")
		args = append(args, *query.UFloor)
	}
	if query.UNumberOfLivingRooms != nil {
		searchQueries = append(searchQueries, "units.number_of_living_rooms = $?")
		args = append(args, *query.UNumberOfLivingRooms)
	}
	if query.UNumberOfBedrooms != nil {
		searchQueries = append(searchQueries, "units.number_of_bedrooms = $?")
		args = append(args, *query.UNumberOfBedrooms)
	}
	if query.UNumberOfBathrooms != nil {
		searchQueries = append(searchQueries, "units.number_of_bathrooms = $?")
		args = append(args, *query.UNumberOfBathrooms)
	}
	if query.UNumberOfToilets != nil {
		searchQueries = append(searchQueries, "units.number_of_toilets = $?")
		args = append(args, *query.UNumberOfToilets)
	}
	if query.UNumberOfKitchens != nil {
		searchQueries = append(searchQueries, "units.number_of_kitchens = $?")
		args = append(args, *query.UNumberOfKitchens)
	}
	if query.UNumberOfBalconies != nil {
		searchQueries = append(searchQueries, "units.number_of_balconies = $?")
		args = append(args, *query.UNumberOfBalconies)
	}
	if len(query.UAmenities) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM unit_amenities WHERE unit_id = units.id AND amenity_id IN ($?))")
		args = append(args, sqlbuilder.List(query.UAmenities))
	}

	// no field is specified and check exisence only
	if len(searchQueries) == 0 && searchFields[0] == "1" {
		return "", []interface{}{}
	}

	if len(connectID) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("units.id = %v", connectID))
	}
	if len(connectProperty) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("units.property_id = %v", connectProperty))
	}
	if len(searchQueries) > 0 {
		// some fields are specified
		searchQuery += " WHERE " + strings.Join(searchQueries, " AND ")
	}
	return searchQuery, args
}
