package sqlbuilders

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
)

func SearchPropertyBuilder(
	searchFields []string, query *dto.SearchPropertyQuery,
	connectID, connectCreator string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM properties WHERE "
	var searchQueries []string
	var args []interface{}

	if query.PIsPublic != nil {
		searchQueries = append(searchQueries, "properties.is_public = $?")
		args = append(args, *query.PIsPublic)
	}
	if query.PName != nil {
		searchQueries = append(searchQueries, "properties.name ILIKE $?")
		args = append(args, "%"+(*query.PName)+"%")
	}
	if query.PCreatorID != nil {
		searchQueries = append(searchQueries, "properties.creator_id = $?")
		args = append(args, *query.PCreatorID)
	}
	if query.PBuilding != nil {
		searchQueries = append(searchQueries, "properties.building ILIKE $?")
		args = append(args, "%"+(*query.PBuilding)+"%")
	}
	if query.PProject != nil {
		searchQueries = append(searchQueries, "properties.project ILIKE $?")
		args = append(args, "%"+(*query.PProject)+"%")
	}
	if query.PFullAddress != nil {
		searchQueries = append(searchQueries, "properties.full_address ILIKE $?")
		args = append(args, "%"+(*query.PFullAddress)+"%")
	}
	if query.PCity != nil {
		searchQueries = append(searchQueries, "properties.city = $?")
		args = append(args, *query.PCity)
	}
	if query.PDistrict != nil {
		searchQueries = append(searchQueries, "properties.district = $?")
		args = append(args, *query.PDistrict)
	}
	if query.PWard != nil {
		searchQueries = append(searchQueries, "properties.ward = $?")
		args = append(args, *query.PWard)
	}
	if query.PMinArea != nil {
		searchQueries = append(searchQueries, "properties.area >= $?")
		args = append(args, *query.PMinArea)
	}
	if query.PMaxArea != nil {
		searchQueries = append(searchQueries, "properties.area <= $?")
		args = append(args, *query.PMaxArea)
	}
	if query.PNumberOfFloors != nil {
		searchQueries = append(searchQueries, "properties.number_of_floors = $?")
		args = append(args, *query.PNumberOfFloors)
	}
	if query.PYearBuilt != nil {
		searchQueries = append(searchQueries, "properties.year_built = $?")
		args = append(args, *query.PYearBuilt)
	}
	if query.POrientation != nil {
		searchQueries = append(searchQueries, "properties.orientation = $?")
		args = append(args, *query.POrientation)
	}
	if query.PMinFacade != nil {
		searchQueries = append(searchQueries, "properties.facade >= $?")
		args = append(args, *query.PMinFacade)
	}
	if len(query.PTypes) > 0 {
		searchQueries = append(searchQueries, "properties.type IN ($?)")
		args = append(args, sqlbuilder.List(query.PTypes))
	}
	if query.PMinCreatedAt != nil {
		searchQueries = append(searchQueries, "properties.created_at >= $?")
		args = append(args, *query.PMinCreatedAt)
	}
	if query.PMaxCreatedAt != nil {
		searchQueries = append(searchQueries, "properties.created_at <= $?")
		args = append(args, *query.PMaxCreatedAt)
	}
	if query.PMinUpdatedAt != nil {
		searchQueries = append(searchQueries, "properties.updated_at >= $?")
		args = append(args, *query.PMinUpdatedAt)
	}
	if query.PMaxUpdatedAt != nil {
		searchQueries = append(searchQueries, "properties.updated_at <= $?")
		args = append(args, *query.PMaxUpdatedAt)
	}
	if len(query.PFeatures) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_features WHERE property_id = properties.id AND feature_id IN ($?))")
		args = append(args, sqlbuilder.List(query.PFeatures))
	}
	if len(query.PTags) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_tags WHERE property_id = properties.id AND tag IN ($?))")
		args = append(args, sqlbuilder.List(query.PTags))
	}

	if len(searchQueries) == 0 {
		return "", []interface{}{}
	}
	if len(connectID) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.id = %v", connectID))
	}
	if len(connectCreator) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.creator_id = %v", connectCreator))
	}
	searchQuery += strings.Join(searchQueries, " AND \n")
	return searchQuery, args
}
