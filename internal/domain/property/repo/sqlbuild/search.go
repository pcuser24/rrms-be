package sqlbuild

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	unit_sqlbuild "github.com/user2410/rrms-backend/internal/domain/unit/repo/sqlbuild"
	"github.com/user2410/rrms-backend/internal/utils"
)

func SearchPropertyBuilder(
	searchFields []string, query *dto.SearchPropertyQuery,
	connectID, connectCreator string,
) (string, []interface{}) {
	var searchQuery string = "SELECT " + strings.Join(searchFields, ", ") + " FROM properties"
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
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_features WHERE property_features.property_id = properties.id AND feature_id IN ($?))")
		args = append(args, sqlbuilder.List(query.PFeatures))
	}
	if len(query.PManagerIDS) > 0 {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_managers WHERE property_managers.property_id = properties.id AND property_managers.manager_id::TEXT IN ($?))")
		args = append(args, sqlbuilder.List(query.PManagerIDS))
	}
	if query.PManagerRole != nil {
		searchQueries = append(searchQueries, "EXISTS (SELECT 1 FROM property_managers WHERE property_managers.property_id = properties.id AND property_managers.role = $?)")
		args = append(args, *query.PManagerRole)
	}
	// no field is specified and check exisence only
	if len(searchQueries) == 0 && searchFields[0] == "1" {
		return "", []interface{}{}
	}

	if len(connectID) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.id = %v", connectID))
	}
	if len(connectCreator) > 0 {
		searchQueries = append(searchQueries, fmt.Sprintf("properties.creator_id = %v", connectCreator))
	}
	if len(searchQueries) > 0 {
		// some fields are specified
		searchQuery += " WHERE " + strings.Join(searchQueries, " AND ")
	}
	return searchQuery, args
}

func SearchPropertyCombinationBuilder(query *dto.SearchPropertyCombinationQuery) (string, []any) {
	sqlProp, argsProp := SearchPropertyBuilder(
		[]string{"properties.id", "count(*) OVER() AS full_count"},
		&query.SearchPropertyQuery,
		"", "",
	)
	sqlUnit, argsUnit := unit_sqlbuild.SearchUnitBuilder([]string{"1"}, &query.SearchUnitQuery, "", "properties.id")

	var queryStr string = sqlProp
	var argsLs []interface{} = argsProp
	// build order: unit -> property
	if len(argsProp) > 0 && len(sqlUnit) > 0 {
		queryStr += fmt.Sprintf(" AND EXISTS (%v)", sqlUnit)
		argsLs = append(argsLs, argsUnit...)
	} else if len(sqlUnit) > 0 {
		queryStr += fmt.Sprintf(" WHERE EXISTS (%v)", sqlUnit)
		argsLs = append(argsLs, argsUnit...)
	}

	sql, args := sqlbuilder.Build(queryStr, argsLs...).Build()
	sqSql := utils.SequelizePlaceholders(sql)

	sqSql += " ORDER BY "
	sortOrders := make([]string, 0, len(query.SortBy))
	for i := 0; i < len(query.SortBy); i++ {
		sortOrders = append(sortOrders, fmt.Sprintf("%v %v", query.SortBy[i], query.Order[i]))
	}
	sqSql += strings.Join(sortOrders, ", ")
	sqSql += fmt.Sprintf(" LIMIT %v", *query.Limit)
	sqSql += fmt.Sprintf(" OFFSET %v", *query.Offset)

	return sqSql, args
}
