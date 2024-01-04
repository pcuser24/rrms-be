package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type CreatePropertyManager struct {
	PropertyID uuid.UUID `json:"propertyID" validate:"required,uuid"`
	ManagerID  uuid.UUID `json:"managerID" validate:"required,uuid"`
	Role       string    `json:"role" validate:"required"`
}

type CreatePropertyMedia struct {
	// PropertyID uuid.UUID          `json:"property_id" validate:"required,uuid"`
	Url         string             `json:"url" validate:"required,url"`
	Type        database.MEDIATYPE `json:"type" validate:"required,oneof=IMAGE VIDEO"`
	Description *string            `json:"description"`
}

func (m *CreatePropertyMedia) ToCreatePropertyMediaDB() *database.CreatePropertyMediaParams {
	return &database.CreatePropertyMediaParams{
		Url:         m.Url,
		Type:        m.Type,
		Description: types.StrN(m.Description),
	}
}

type CreatePropertyFeature struct {
	// PropertyID  uuid.UUID `json:"property_id" validate:"required,uuid"`
	FeatureID   int64   `json:"featureID" validate:"required"`
	Description *string `json:"description"`
}

func (f *CreatePropertyFeature) ToCreatePropertyFeatureDB() *database.CreatePropertyFeatureParams {
	return &database.CreatePropertyFeatureParams{
		FeatureID:   f.FeatureID,
		Description: types.StrN(f.Description),
	}
}

type CreatePropertyTag struct {
	// PropertyID uuid.UUID `json:"property_id" validate:"required,uuid"`
	Tag string `json:"tag" validate:"required"`
}

type CreateProperty struct {
	CreatorID      uuid.UUID               `json:"creatorId" validate:"required,uuid4"`
	Name           *string                 `json:"name" validate:"omitempty"`
	Building       *string                 `json:"building" validate:"omitempty"`
	Project        *string                 `json:"project" validate:"omitempty"`
	Area           float32                 `json:"area" validate:"required,gt=0"`
	NumberOfFloors *int32                  `json:"numberOfFloors" validate:"omitempty,gt=0"`
	YearBuilt      *int32                  `json:"yearBuilt" validate:"omitempty,gt=0"`
	Orientation    *string                 `json:"orientation" validate:"omitempty"`
	EntranceWidth  *float64                `json:"entranceWidth" validate:"omitempty"`
	Facade         *float64                `json:"facade" validate:"omitempty"`
	FullAddress    string                  `json:"fullAddress" validate:"required"`
	District       string                  `json:"district" validate:"required"`
	City           string                  `json:"city" validate:"required"`
	Ward           *string                 `json:"ward" validate:"omitempty"`
	PlaceUrl       string                  `json:"placeUrl" validate:"required,url"`
	Lat            *float64                `json:"lat" validate:"omitempty"`
	Lng            *float64                `json:"lng" validate:"omitempty"`
	Description    *string                 `json:"description" validate:"omitempty"`
	Type           database.PROPERTYTYPE   `json:"type" validate:"required,oneof=APARTMENT PRIVATE TOWNHOUSE SHOPHOUSE VILLA ROOM STORE OFFICE BLOCK COMPLEX"`
	Managers       []CreatePropertyManager `json:"managers" validate:"dive"`
	Media          []CreatePropertyMedia   `json:"media" validate:"dive"`
	Features       []CreatePropertyFeature `json:"features" validate:"dive"`
	Tags           []CreatePropertyTag     `json:"tags" validate:"dive"`
}

func (c *CreateProperty) ToCreatePropertyDB() *database.CreatePropertyParams {
	p := &database.CreatePropertyParams{
		CreatorID:   c.CreatorID,
		Area:        c.Area,
		FullAddress: c.FullAddress,
		District:    c.District,
		City:        c.City,
		PlaceUrl:    c.PlaceUrl,
		Type:        c.Type,
	}
	if c.Name != nil {
		p.Name = sql.NullString{
			Valid:  true,
			String: *c.Name,
		}
	}
	if c.Building != nil {
		p.Building = sql.NullString{
			Valid:  true,
			String: *c.Building,
		}
	}
	if c.Project != nil {
		p.Project = sql.NullString{
			Valid:  true,
			String: *c.Project,
		}
	}
	if c.NumberOfFloors != nil {
		p.NumberOfFloors = sql.NullInt32{
			Valid: true,
			Int32: *c.NumberOfFloors,
		}
	}
	if c.YearBuilt != nil {
		p.YearBuilt = sql.NullInt32{
			Valid: true,
			Int32: *c.YearBuilt,
		}
	}
	if c.Orientation != nil {
		p.Orientation = sql.NullString{
			Valid:  true,
			String: *c.Orientation,
		}
	}
	if c.EntranceWidth != nil {
		p.EntranceWidth = sql.NullFloat64{
			Valid:   true,
			Float64: *c.EntranceWidth,
		}
	}
	if c.Facade != nil {
		p.Facade = sql.NullFloat64{
			Valid:   true,
			Float64: *c.Facade,
		}
	}
	if c.Ward != nil {
		p.Ward = sql.NullString{
			Valid:  true,
			String: *c.Ward,
		}
	}
	if c.Lat != nil {
		p.Lat = sql.NullFloat64{
			Valid:   true,
			Float64: *c.Lat,
		}
	}
	if c.Lng != nil {
		p.Lng = sql.NullFloat64{
			Valid:   true,
			Float64: *c.Lng,
		}
	}
	if c.Description != nil {
		p.Description = sql.NullString{
			Valid:  true,
			String: *c.Description,
		}
	}
	return p
}
