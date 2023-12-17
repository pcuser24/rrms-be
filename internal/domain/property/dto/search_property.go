package dto

import "time"

type SearchPropertyQuery struct {
	PTypes          []string   `query:"ptypes"`
	PCreatorID      *string    `query:"pcreatorId"`
	PName           *string    `query:"pname"`
	PBuilding       *string    `query:"pbuilding"`
	PProject        *string    `query:"pproject"`
	PFullAddress    *string    `query:"pfullAddress"`
	PCity           *string    `query:"pcity"`
	PDistrict       *string    `query:"pdistrict"`
	PWard           *string    `query:"pward"`
	PMinArea        *float32   `query:"pminArea"`
	PMaxArea        *float32   `query:"pmaxArea"`
	PNumberOfFloors *int32     `query:"pnumberOfFloors"`
	PYearBuilt      *int32     `query:"pyearBuilt"`
	POrientation    *string    `query:"porientation"`
	PMinFacade      *int32     `query:"pfacade"`
	PIsPublic       *bool      `query:"pisPublic"`
	PFeatures       []int32    `query:"pfeatures"`
	PTags           []string   `query:"ptags"`
	PMinCreatedAt   *time.Time `query:"pminCreatedAt"`
	PMaxCreatedAt   *time.Time `query:"pmaxCreatedAt"`
	PMinUpdatedAt   *time.Time `query:"pminUpdatedAt"`
	PMaxUpdatedAt   *time.Time `query:"pmaxUpdatedAt"`
}
