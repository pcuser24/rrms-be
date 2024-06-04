package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PreUpdatePropertyMedia struct {
	ID          int64   `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Size        int64   `json:"size" validate:"required,gt=0"`
	Type        string  `json:"type" validate:"required"`
	Description *string `json:"description"`
	Url         string  `json:"url"`
}

type PreUpdateProperty struct {
	Media []PreCreatePropertyMedia `json:"media" validate:"dive"`
}

type UpdateProperty struct {
	ID              uuid.UUID
	Name            *string                      `json:"name" validate:"omitempty"`
	Building        *string                      `json:"building" validate:"omitempty"`
	Project         *string                      `json:"project" validate:"omitempty"`
	Area            *float32                     `json:"area" validate:"omitempty,gte=0"`
	NumberOfFloors  *int32                       `json:"numberOfFloors" validate:"omitempty,gte=0"`
	YearBuilt       *int32                       `json:"yearBuilt" validate:"omitempty,gte=0"`
	Orientation     *string                      `json:"orientation" validate:"omitempty,oneof=n s e w ne nw se sw"`
	EntranceWidth   *float32                     `json:"entranceWidth" validate:"omitempty,gte=0"`
	Facade          *float32                     `json:"facade" validate:"omitempty,gte=0"`
	FullAddress     *string                      `json:"fullAddress" validate:"omitempty"`
	District        *string                      `json:"district" validate:"omitempty"`
	City            *string                      `json:"city" validate:"omitempty"`
	Ward            *string                      `json:"ward" validate:"omitempty"`
	Lat             *float64                     `json:"lat" validate:"omitempty"`
	Lng             *float64                     `json:"lng" validate:"omitempty"`
	PrimaryImage    *int64                       `json:"primaryImage" validate:"omitempty"`
	PrimaryImageUrl *string                      `json:"primaryImageUrl" validate:"omitempty"`
	Description     *string                      `json:"description" validate:"omitempty"`
	IsPublic        *bool                        `json:"isPublic" validate:"omitempty"`
	Media           []model.PropertyMediaModel   `json:"media" validate:"omitempty,dive"`
	Features        []model.PropertyFeatureModel `json:"features" validate:"omitempty,dive"`
	Managers        []model.PropertyManagerModel `json:"managers" validate:"omitempty,dive"`
}

func (u *UpdateProperty) ToUpdatePropertyDB() database.UpdatePropertyParams {
	return database.UpdatePropertyParams{
		ID:             u.ID,
		Name:           types.StrN(u.Name),
		Building:       types.StrN(u.Building),
		Project:        types.StrN(u.Project),
		Area:           types.Float32N(u.Area),
		NumberOfFloors: types.Int32N(u.NumberOfFloors),
		YearBuilt:      types.Int32N(u.YearBuilt),
		Orientation:    types.StrN(u.Orientation),
		EntranceWidth:  types.Float32N(u.EntranceWidth),
		Facade:         types.Float32N(u.Facade),
		FullAddress:    types.StrN(u.FullAddress),
		District:       types.StrN(u.District),
		City:           types.StrN(u.City),
		Ward:           types.StrN(u.Ward),
		Lat:            types.Float64N(u.Lat),
		Lng:            types.Float64N(u.Lng),
		PrimaryImage:   types.Int64N(u.PrimaryImage),
		Description:    types.StrN(u.Description),
		IsPublic:       types.BoolN(u.IsPublic),
	}
}

type CreatePropertyManagerRequest struct {
	PropertyID uuid.UUID
	CreatorID  uuid.UUID
	UserID     uuid.UUID
	Email      string `json:"email" validate:"required,email"`
}
