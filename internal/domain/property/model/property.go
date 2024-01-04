package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
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
	EntranceWidth *float64               `json:"entranceWidth"`
	Facade        *float64               `json:"facade"`
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
	m := &PropertyModel{
		ID:          p.ID,
		CreatorID:   p.CreatorID,
		Name:        p.Name,
		Area:        p.Area,
		FullAddress: p.FullAddress,
		District:    p.District,
		City:        p.City,
		PlaceUrl:    p.PlaceUrl,
		Type:        p.Type,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.Building.Valid {
		b := p.Building.String
		m.Building = &b
	}
	if p.Project.Valid {
		p := p.Project.String
		m.Project = &p
	}
	if p.NumberOfFloors.Valid {
		n := p.NumberOfFloors.Int32
		m.NumberOfFloors = &n
	}
	if p.YearBuilt.Valid {
		y := p.YearBuilt.Int32
		m.YearBuilt = &y
	}
	if p.Orientation.Valid {
		o := p.Orientation.String
		m.Orientation = &o
	}
	if p.EntranceWidth.Valid {
		e := p.EntranceWidth.Float64
		m.EntranceWidth = &e
	}
	if p.Facade.Valid {
		f := p.Facade.Float64
		m.Facade = &f
	}
	if p.Ward.Valid {
		w := p.Ward.String
		m.Ward = &w
	}
	if p.Lat.Valid {
		l := p.Lat.Float64
		m.Lat = &l
	}
	if p.Lng.Valid {
		l := p.Lng.Float64
		m.Lng = &l
	}
	if p.Description.Valid {
		d := p.Description.String
		m.Description = &d
	}
	return m
}
