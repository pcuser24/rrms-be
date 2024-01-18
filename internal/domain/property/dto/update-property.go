package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateProperty struct {
	ID             uuid.UUID
	Name           *string  `json:"name"`
	Building       *string  `json:"building"`
	Project        *string  `json:"project"`
	Area           *float32 `json:"area" validate:"gte:0"`
	NumberOfFloors *int32   `json:"number_of_floors" validate:"gte:0"`
	YearBuilt      *int32   `json:"year_built" validate:"gte:0"`
	Orientation    *string  `json:"orientation" validate:"oneof:n,s,e,w,ne,nw,se,sw"`
	EntranceWidth  *float32 `json:"entrance_width" validate:"gte:0"`
	Facade         *float32 `json:"facade" validate:"gte:0"`
	FullAddress    *string  `json:"full_address"`
	District       *string  `json:"district"`
	City           *string  `json:"city"`
	Ward           *string  `json:"ward"`
	Lat            *float64 `json:"lat"`
	Lng            *float64 `json:"lng"`
	PlaceUrl       *string  `json:"place_url"`
	Description    *string  `json:"description"`
}

func (u *UpdateProperty) ToUpdatePropertyDB() *database.UpdatePropertyParams {
	return &database.UpdatePropertyParams{
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
		PlaceUrl:       types.StrN(u.PlaceUrl),
		Description:    types.StrN(u.Description),
	}
}
