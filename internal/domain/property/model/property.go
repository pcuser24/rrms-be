package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PROPERTYTYPE string

const (
	PROPERTYTYPEAPARTMENT     PROPERTYTYPE = "APARTMENT"
	PROPERTYTYPEPRIVATE       PROPERTYTYPE = "PRIVATE"
	PROPERTYTYPEROOM          PROPERTYTYPE = "ROOM"
	PROPERTYTYPESTORE         PROPERTYTYPE = "STORE"
	PROPERTYTYPEOFFICE        PROPERTYTYPE = "OFFICE"
	PROPERTYTYPEVILLA         PROPERTYTYPE = "VILLA"
	PROPERTYTYPEMINIAPARTMENT PROPERTYTYPE = "MINIAPARTMENT"
)

func (pt PROPERTYTYPE) MarshalBinary() ([]byte, error) {
	return json.Marshal(pt)
}

func (pt *PROPERTYTYPE) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pt)
}

type PropertyTagModel struct {
	ID         int64     `json:"id" redis:"id"`
	PropertyID uuid.UUID `json:"propertyId" redis:"propertyId"`
	Tag        string    `json:"tag" redis:"tag"`
}

func (pt PropertyTagModel) MarshalBinary() (data []byte, err error) {
	return json.Marshal(pt)
}

func (pt *PropertyTagModel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pt)
}

type PropertyModel struct {
	ID             uuid.UUID              `json:"id" redis:"id"`
	CreatorID      uuid.UUID              `json:"creatorId" redis:"creatorId"`
	Name           string                 `json:"name" redis:"name"`
	Building       *string                `json:"building" redis:"building"`
	Project        *string                `json:"project" redis:"project"`
	Area           float32                `json:"area" redis:"area"`
	NumberOfFloors *int32                 `json:"numberOfFloors" redis:"numberOfFloors"`
	YearBuilt      *int32                 `json:"yearBuilt" redis:"yearBuilt"`
	Orientation    *string                `json:"orientation" redis:"orientation"`
	EntranceWidth  *float32               `json:"entranceWidth" redis:"entranceWidth"`
	Facade         *float32               `json:"facade" redis:"facade"`
	FullAddress    string                 `json:"fullAddress" redis:"fullAddress"`
	District       string                 `json:"district" redis:"district"`
	City           string                 `json:"city" redis:"city"`
	Ward           *string                `json:"ward" redis:"ward"`
	Lat            *float64               `json:"lat" redis:"lat"`
	Lng            *float64               `json:"lng" redis:"lng"`
	PrimaryImage   int64                  `json:"primaryImage" redis:"primaryImage"`
	Description    *string                `json:"description" redis:"description"`
	Type           PROPERTYTYPE           `json:"type" redis:"type"`
	IsPublic       bool                   `json:"isPublic" redis:"isPublic"`
	CreatedAt      time.Time              `json:"createdAt" redis:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt" redis:"updatedAt"`
	Managers       []PropertyManagerModel `json:"managers" redis:"managers"`
	Features       []PropertyFeatureModel `json:"features" redis:"features"`
	Media          []PropertyMediaModel   `json:"media" redis:"media"`
	Tags           []PropertyTagModel     `json:"tags" redis:"tags"`
}

func ToPropertyModel(p *database.Property) *PropertyModel {
	return &PropertyModel{
		ID:             p.ID,
		CreatorID:      p.CreatorID,
		Name:           p.Name,
		Building:       types.PNStr(p.Building),
		Project:        types.PNStr(p.Project),
		NumberOfFloors: types.PNInt32(p.NumberOfFloors),
		YearBuilt:      types.PNInt32(p.YearBuilt),
		Orientation:    types.PNStr(p.Orientation),
		EntranceWidth:  types.PNFloat32(p.EntranceWidth),
		Facade:         types.PNFloat32(p.Facade),
		Area:           p.Area,
		FullAddress:    p.FullAddress,
		District:       p.District,
		City:           p.City,
		Ward:           types.PNStr(p.Ward),
		Lat:            types.PNFloat64(p.Lat),
		Lng:            types.PNFloat64(p.Lng),
		PrimaryImage:   p.PrimaryImage.Int64,
		Type:           PROPERTYTYPE(p.Type),
		IsPublic:       p.IsPublic,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Description:    types.PNStr(p.Description),
		Managers:       make([]PropertyManagerModel, 0),
		Features:       make([]PropertyFeatureModel, 0),
		Media:          make([]PropertyMediaModel, 0),
		Tags:           make([]PropertyTagModel, 0),
	}
}

func (p *PropertyModel) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
