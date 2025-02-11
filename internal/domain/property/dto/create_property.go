package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreatePropertyManager struct {
	ManagerID uuid.UUID `json:"managerId" validate:"required,uuid4"`
	Role      string    `json:"role" validate:"required"`
}

type PreCreatePropertyMedia struct {
	ID   int64  `json:"id" validate:"required,gt=0"`
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required,gt=0"`
	Type string `json:"type" validate:"required"`
	Url  string `json:"url"`
}

type PreCreateProperty struct {
	Media []PreCreatePropertyMedia `json:"media" validate:"dive"`
}

type CreatePropertyMedia struct {
	Url         string             `json:"url" validate:"required,url"`
	Type        database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
	Description *string            `json:"description"`
}

func (m *CreatePropertyMedia) ToCreatePropertyMediaDB(propertyId uuid.UUID) *database.CreatePropertyMediaParams {
	return &database.CreatePropertyMediaParams{
		PropertyID:  propertyId,
		Url:         m.Url,
		Type:        m.Type,
		Description: types.StrN(m.Description),
	}
}

type CreatePropertyFeature struct {
	FeatureID   int64   `json:"featureId" validate:"required"`
	Description *string `json:"description"`
}

func (f *CreatePropertyFeature) ToCreatePropertyFeatureDB(propertyId uuid.UUID) *database.CreatePropertyFeatureParams {
	return &database.CreatePropertyFeatureParams{
		PropertyID:  propertyId,
		FeatureID:   f.FeatureID,
		Description: types.StrN(f.Description),
	}
}

type CreatePropertyTag struct {
	Tag string `json:"tag" validate:"required"`
}

type CreateProperty struct {
	CreatorID      uuid.UUID               `json:"creatorId"`
	Name           string                  `json:"name" validate:"required"`
	Building       *string                 `json:"building" validate:"omitempty"`
	Project        *string                 `json:"project" validate:"omitempty"`
	Area           float32                 `json:"area" validate:"required,gt=0"`
	NumberOfFloors *int32                  `json:"numberOfFloors" validate:"omitempty,gt=0"`
	YearBuilt      *int32                  `json:"yearBuilt" validate:"omitempty,gt=0"`
	Orientation    *string                 `json:"orientation" validate:"omitempty"`
	EntranceWidth  *float32                `json:"entranceWidth" validate:"omitempty"`
	Facade         *float32                `json:"facade" validate:"omitempty"`
	FullAddress    string                  `json:"fullAddress" validate:"required"`
	District       string                  `json:"district" validate:"required"`
	City           string                  `json:"city" validate:"required"`
	Ward           *string                 `json:"ward" validate:"omitempty"`
	Lat            *float64                `json:"lat" validate:"omitempty"`
	Lng            *float64                `json:"lng" validate:"omitempty"`
	Description    *string                 `json:"description" validate:"omitempty"`
	Type           database.PROPERTYTYPE   `json:"type" validate:"required,oneof=APARTMENT PRIVATE ROOM STORE OFFICE VILLA MINIAPARTMENT"`
	PrimaryImage   string                  `json:"primaryImage" validate:"required,url"`
	Managers       []CreatePropertyManager `json:"managers" validate:"dive"`
	Media          []CreatePropertyMedia   `json:"media" validate:"dive"`
	Features       []CreatePropertyFeature `json:"features" validate:"dive"`
	Tags           []CreatePropertyTag     `json:"tags" validate:"dive"`
}

func (c *CreateProperty) ToCreatePropertyDB() *database.CreatePropertyParams {
	return &database.CreatePropertyParams{
		CreatorID:      c.CreatorID,
		Name:           c.Name,
		Building:       types.StrN(c.Building),
		Project:        types.StrN(c.Project),
		NumberOfFloors: types.Int32N(c.NumberOfFloors),
		YearBuilt:      types.Int32N(c.YearBuilt),
		Orientation:    types.StrN(c.Orientation),
		EntranceWidth:  types.Float32N(c.EntranceWidth),
		Facade:         types.Float32N(c.Facade),
		Ward:           types.StrN(c.Ward),
		Lat:            types.Float64N(c.Lat),
		Lng:            types.Float64N(c.Lng),
		Area:           c.Area,
		FullAddress:    c.FullAddress,
		District:       c.District,
		City:           c.City,
		// PrimaryImage: omitted
		Type:        c.Type,
		Description: types.StrN(c.Description),
	}
}
