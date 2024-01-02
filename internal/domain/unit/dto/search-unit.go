package dto

type SearchUnitQuery struct {
	UPropertyID          *string `query:"upropertyId"`
	UName                *string `query:"uname"`
	UMinArea             *int64  `query:"uminPrice"`
	UMaxArea             *int64  `query:"umaxPrice"`
	UFloor               *int32  `query:"ufloor"`
	UPrice               int64   `query:"uprice"`
	UNumberOfLivingRooms *int32  `query:"unumberOfLivingRooms"`
	UNumberOfBedrooms    *int32  `query:"unumberOfBedrooms"`
	UNumberOfBathrooms   *int32  `query:"unumberOfBathrooms"`
	UNumberOfToilets     *int32  `query:"unumberOfToilets"`
	UNumberOfKitchens    *int32  `query:"unumberOfKitchens"`
	UNumberOfBalconies   *int32  `query:"unumberOfBalconies"`
	UAmenities           []int32 `query:"uamenities"`
}
