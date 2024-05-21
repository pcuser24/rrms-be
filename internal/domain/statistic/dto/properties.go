package dto

import "github.com/google/uuid"

type PropertiesStatisticQuery struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type PropertiesStatisticResponse struct {
	Properties                  []uuid.UUID                   `json:"properties"`
	OwnedProperties             []uuid.UUID                   `json:"ownedProperties"`
	OccupiedProperties          []uuid.UUID                   `json:"occupiedProperties"`
	Units                       []uuid.UUID                   `json:"units"`
	OccupiedUnits               []uuid.UUID                   `json:"occupiedUnits"`
	PropertiesWithActiveListing []uuid.UUID                   `json:"propertiesWithActiveListing"`
	MostRentedProperties        []ExtremelyRentedPropertyItem `json:"mostRentedProperties"`
	LeastRentedProperties       []ExtremelyRentedPropertyItem `json:"leastRentedProperties"`
	MostRentedUnits             []ExtremelyRentedUnitItem     `json:"mostRentedUnits"`
	LeastRentedUnits            []ExtremelyRentedUnitItem     `json:"leastRentedUnits"`
}

type ExtremelyRentedPropertyItem struct {
	PropertyID uuid.UUID `json:"propertyId"`
	Count      int64     `json:"count"`
}

type ExtremelyRentedUnitItem struct {
	PropertyID uuid.UUID `json:"propertyId"`
	UnitID     uuid.UUID `json:"unitId"`
	Count      int64     `json:"count"`
}
