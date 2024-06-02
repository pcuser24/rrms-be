package service

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionscoremode"
	"github.com/google/uuid"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) GetRecentListings(limit int32, fields []string) ([]listing_model.ListingModel, error) {
	ids, err := s.statisticRepo.GetRecentListings(context.Background(), limit)
	if err != nil {
		return nil, err
	}

	return s.listingRepo.GetListingsByIds(context.Background(), ids, fields)
}

func (s *service) GetListingSuggestion(id uuid.UUID, limit int) (dto.ListingsSuggestionResult, error) {
	var res dto.ListingsSuggestionResult

	listing, err := s.listingRepo.GetListingByID(context.Background(), id)
	if err != nil {
		return dto.ListingsSuggestionResult{}, err
	}
	property, err := s.propertyRepo.GetPropertyById(context.Background(), listing.PropertyID)
	if err != nil {
		return dto.ListingsSuggestionResult{}, err
	}
	unitIds := make([]uuid.UUID, 0, len(listing.Units))
	for _, unit := range listing.Units {
		unitIds = append(unitIds, unit.UnitID)
	}
	units, err := s.unitRepo.GetUnitsByIds(context.Background(), unitIds, []string{"area", "floor", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "amenities"})
	if err != nil {
		return dto.ListingsSuggestionResult{}, err
	}
	// get a list of ids of amenities of units
	amenityIds := make([]int64, 0)
	for _, u := range units {
		for _, a := range u.Amenities {
			amenityIds = append(amenityIds, a.AmenityID)
		}
	}

	// suggestion by es
	esClient := s.esClient.GetTypedClient()
	scoreFns := []estypes.FunctionScore{
		{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.district": {
						Value: property.District,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		},
		{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"price": estypes.NumberRangeQuery{
						Gte: types.Ptr(estypes.Float64(listing.Price * 0.8)),
						Lte: types.Ptr(estypes.Float64(listing.Price * 1.2)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](5),
		},
		{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"property.area": estypes.NumberRangeQuery{
						Gte: types.Ptr(estypes.Float64(property.Area * 0.8)),
						Lte: types.Ptr(estypes.Float64(property.Area * 1.2)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](5),
		},
		{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.type": {
						Value: property.Type,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](4),
		},
		{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_bedrooms": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(
									utils.PtrDerefence(units[0].NumberOfBedrooms, 1)),
								),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		},
		{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_bathrooms": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(
									utils.PtrDerefence(units[0].NumberOfBedrooms, 1)),
								),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		},
	}
	if property.Ward != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.ward": {
						Value: property.Ward,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		})
	}
	if property.Project != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.project": {
						Value: property.Project,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](4),
		})
	}
	if property.Orientation != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.orientation": {
						Value: property.Orientation,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if len(property.Features) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "property.features",
					Query: &estypes.Query{
						Terms: &estypes.TermsQuery{
							TermsQuery: map[string]estypes.TermsQueryField{
								"property.features.feature_id": func() []int64 {
									var result []int64
									for _, f := range property.Features {
										result = append(result, f.FeatureID)
									}
									return result
								}(),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if len(amenityIds) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units.amenities",
					Query: &estypes.Query{
						Terms: &estypes.TermsQuery{
							TermsQuery: map[string]estypes.TermsQueryField{
								"listing_units.amenities.amenity_id": amenityIds,
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	search := esClient.Search().
		Index(string(es.LISTINGINDEX)).
		Request(&search.Request{
			Size: types.Ptr(limit),
			Query: &estypes.Query{
				FunctionScore: &estypes.FunctionScoreQuery{
					Query: &estypes.Query{
						Bool: &estypes.BoolQuery{
							Should: []estypes.Query{
								{
									Term: map[string]estypes.TermQuery{
										"active": {
											Value: true,
										},
									},
								},
								{
									Range: map[string]estypes.RangeQuery{
										"expired_at": estypes.DateRangeQuery{
											Gte: types.Ptr("now"),
										},
									},
								},
								{
									Term: map[string]estypes.TermQuery{
										"property.city": {
											Value: property.City,
										},
									},
								},
							},
							MustNot: []estypes.Query{
								{
									Term: map[string]estypes.TermQuery{
										"id": {
											Value: id.String(),
										},
									},
								},
							},
						},
					},
					Functions: scoreFns,
					ScoreMode: &functionscoremode.Sum,
					BoostMode: &functionboostmode.Sum,
				},
			},
		},
		)

	searchRes, err := search.Do(context.Background())
	if err != nil {
		return dto.ListingsSuggestionResult{}, err
	}
	res.Hits = make([]dto.SuggestedListing, 0, len(searchRes.Hits.Hits))
	for _, h := range searchRes.Hits.Hits {
		var i dto.SuggestedListing
		err := json.Unmarshal(h.Source_, &i)
		if err != nil {
			return dto.ListingsSuggestionResult{}, err
		}
		res.Hits = append(res.Hits, i)
	}

	// suggestion by db
	return res, nil
}
