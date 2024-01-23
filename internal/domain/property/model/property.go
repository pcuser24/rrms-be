package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PropertyTagModel = database.PropertyTag

type PropertyModel struct {
	ID             uuid.UUID `json:"id"`
	CreatorID      uuid.UUID `json:"creatorId"`
	Name           string    `json:"name"`
	Building       *string   `json:"building"`
	Project        *string   `json:"project"`
	Area           float32   `json:"area"`
	NumberOfFloors *int32    `json:"numberOfFloors"`
	YearBuilt      *int32    `json:"yearBuilt"`
	// n,s,w,e,nw,ne,sw,se
	Orientation   *string                `json:"orientation"`
	EntranceWidth *float32               `json:"entranceWidth"`
	Facade        *float32               `json:"facade"`
	FullAddress   string                 `json:"fullAddress"`
	District      string                 `json:"district"`
	City          string                 `json:"city"`
	Ward          *string                `json:"ward"`
	Lat           *float64               `json:"lat"`
	Lng           *float64               `json:"lng"`
	PlaceUrl      string                 `json:"placeUrl"`
	Description   *string                `json:"description"`
	Type          database.PROPERTYTYPE  `json:"type"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	Managers      []PropertyManagerModel `json:"managers"`
	Features      []PropertyFeatureModel `json:"features"`
	Media         []PropertyMediaModel   `json:"media"`
	Tags          []PropertyTagModel     `json:"tags"`
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
		PlaceUrl:       p.PlaceUrl,
		Type:           p.Type,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Description:    types.PNStr(p.Description),
		Managers:       make([]PropertyManagerModel, 0),
		Features:       make([]PropertyFeatureModel, 0),
		Media:          make([]PropertyMediaModel, 0),
		Tags:           make([]PropertyTagModel, 0),
	}
}
