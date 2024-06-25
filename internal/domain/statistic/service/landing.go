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
	statistic_dto "github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) GetRecentListings(limit int32, fields []string) ([]listing_model.ListingModel, error) {
	ids, err := s.domainRepo.StatisticRepo.GetRecentListings(context.Background(), limit)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.ListingRepo.GetListingsByIds(context.Background(), ids, fields)
}

func (s *service) GetSimilarListingsToListing(id uuid.UUID, limit int) (statistic_dto.ListingsSuggestionResult, error) {
	listing, err := s.domainRepo.ListingRepo.GetListingByID(context.Background(), id)
	if err != nil {
		return statistic_dto.ListingsSuggestionResult{}, err
	}
	property, err := s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), listing.PropertyID)
	if err != nil {
		return statistic_dto.ListingsSuggestionResult{}, err
	}
	unitIds := make([]uuid.UUID, 0, len(listing.Units))
	for _, unit := range listing.Units {
		unitIds = append(unitIds, unit.UnitID)
	}
	units, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), unitIds, []string{"area", "floor", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "amenities"})
	if err != nil {
		return statistic_dto.ListingsSuggestionResult{}, err
	}
	// get a list of ids of amenities of units
	amenityIds := make([]int64, 0)
	for _, u := range units {
		for _, a := range u.Amenities {
			amenityIds = append(amenityIds, a.AmenityID)
		}
	}

	// prepare queries and score functions
	// should queries
	shouldQueries := []estypes.Query{
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
	}

	// mustnot queries
	mustNotQueries := []estypes.Query{
		{
			Term: map[string]estypes.TermQuery{
				"id": {
					Value: id.String(),
				},
			},
		},
	}

	// prepare query
	query := statistic_dto.ListingSuggestionQuery{
		PTypes:    []string{string(property.Type)},
		PCity:     []string{property.City},
		PDistrict: []string{property.District},
		PMinArea:  types.Ptr(property.Area * 0.8),
		PMaxArea:  types.Ptr(property.Area * 1.2),
		LMinPrice: types.Ptr(listing.Price * 0.8),
		LMaxPrice: types.Ptr(listing.Price * 1.2),
	}
	if property.Ward != nil {
		query.PWard = []string{*property.Ward}
	}
	if property.Orientation != nil {
		query.POrientation = []string{*property.Orientation}
	}
	for _, f := range property.Features {
		query.PFeatures = append(query.PFeatures, f.FeatureID)
	}
	query.UAmenities = amenityIds
	if units[0].NumberOfBedrooms != nil {
		query.UNumberOfBedrooms = types.Ptr(*units[0].NumberOfBedrooms)
	}
	if units[0].NumberOfBathrooms != nil {
		query.UNumberOfBathrooms = types.Ptr(*units[0].NumberOfBathrooms)
	}
	if units[0].NumberOfBalconies != nil {
		query.UNumberOfBalconies = types.Ptr(*units[0].NumberOfBalconies)
	}
	if units[0].NumberOfToilets != nil {
		query.UNumberOfToilets = types.Ptr(*units[0].NumberOfToilets)
	}
	if units[0].NumberOfKitchens != nil {
		query.UNumberOfKitchens = types.Ptr(*units[0].NumberOfKitchens)
	}
	if units[0].NumberOfLivingRooms != nil {
		query.UNumberOfLivingRooms = types.Ptr(*units[0].NumberOfLivingRooms)
	}

	var res statistic_dto.ListingsSuggestionResult
	searchRes, err := s.SuggestSimilarListings(&query, []estypes.Query{}, mustNotQueries, shouldQueries, limit)
	if err != nil {
		return res, err
	}
	res.Hits = make([]statistic_dto.ListingsSuggestionItem, 0, len(searchRes.Hits.Hits))
	for _, h := range searchRes.Hits.Hits {
		var i statistic_dto.ListingsSuggestionItem
		err := json.Unmarshal(h.Source_, &i)
		if err != nil {
			return res, err
		}
		res.Hits = append(res.Hits, i)
	}

	// suggestion by db
	return res, nil
}

func (s *service) SuggestSimilarListings(query *statistic_dto.ListingSuggestionQuery, mustQueries, mustNotQueries, shouldQueries []estypes.Query, limit int) (*search.Response, error) {
	esClient := s.esClient.GetTypedClient()

	// score functions
	scoreFns := []estypes.FunctionScore{
		{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.verificationStatus": {
						Value: "APPROVED",
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		},
		{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.verificationStatus": {
						Value: "PENDING",
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		},
		{
			Filter: &estypes.Query{
				Term: map[string]estypes.TermQuery{
					"property.verificationStatus": {
						Value: "REJECTED",
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		},
		{
			Filter: &estypes.Query{
				Bool: &estypes.BoolQuery{
					MustNot: []estypes.Query{
						{
							Exists: &estypes.ExistsQuery{
								Field: "property.verificationStatus",
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		},
	}
	if len(query.PCity) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Terms: &estypes.TermsQuery{
					TermsQuery: map[string]estypes.TermsQueryField{
						"property.city": query.PCity,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		})
	}
	if len(query.PDistrict) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Terms: &estypes.TermsQuery{
					TermsQuery: map[string]estypes.TermsQueryField{
						"property.district": query.PDistrict,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if len(query.PWard) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Terms: &estypes.TermsQuery{
					TermsQuery: map[string]estypes.TermsQueryField{
						"property.ward": query.PWard,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1.5),
		})
	}
	if query.LMaxPrice != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"price": estypes.NumberRangeQuery{
						Lte: types.Ptr(estypes.Float64(*query.LMaxPrice)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if query.LMinPrice != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"price": estypes.NumberRangeQuery{
						Gte: types.Ptr(estypes.Float64(*query.LMinPrice)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if query.PMaxArea != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"property.area": estypes.NumberRangeQuery{
						Lte: types.Ptr(estypes.Float64(*query.PMaxArea)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if query.PMinArea != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Range: map[string]estypes.RangeQuery{
					"property.area": estypes.NumberRangeQuery{
						Gte: types.Ptr(estypes.Float64(*query.PMinArea)),
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](2),
		})
	}
	if len(query.PTypes) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Terms: &estypes.TermsQuery{
					TermsQuery: map[string]estypes.TermsQueryField{
						"property.type": query.PTypes,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](3),
		})
	}
	if query.UNumberOfBedrooms != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_bedrooms": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfBedrooms)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if query.UNumberOfBathrooms != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_bathrooms": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfBedrooms)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if query.UNumberOfBalconies != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_balconies": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfBalconies)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if query.UNumberOfToilets != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_toilets": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfToilets)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if query.UNumberOfLivingRooms != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_living_rooms": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfLivingRooms)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if query.UNumberOfKitchens != nil {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units",
					Query: &estypes.Query{
						Range: map[string]estypes.RangeQuery{
							"listing_units.number_of_kitchens": estypes.NumberRangeQuery{
								Gte: types.Ptr(estypes.Float64(*query.UNumberOfKitchens)),
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if len(query.POrientation) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Terms: &estypes.TermsQuery{
					TermsQuery: map[string]estypes.TermsQueryField{
						"property.orientation": query.POrientation,
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if len(query.PFeatures) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "property.features",
					Query: &estypes.Query{
						Terms: &estypes.TermsQuery{
							TermsQuery: map[string]estypes.TermsQueryField{
								"property.features.feature_id": query.PFeatures,
							},
						},
					},
				},
			},
			Weight: types.Ptr[estypes.Float64](1),
		})
	}
	if len(query.UAmenities) > 0 {
		scoreFns = append(scoreFns, estypes.FunctionScore{
			Filter: &estypes.Query{
				Nested: &estypes.NestedQuery{
					Path: "listing_units.amenities",
					Query: &estypes.Query{
						Terms: &estypes.TermsQuery{
							TermsQuery: map[string]estypes.TermsQueryField{
								"listing_units.amenities.amenity_id": query.UAmenities,
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
							Should:  shouldQueries,
							Must:    mustQueries,
							MustNot: mustNotQueries,
						},
					},
					Functions: scoreFns,
					ScoreMode: &functionscoremode.Sum,
					BoostMode: &functionboostmode.Sum,
				},
			}})

	return search.Do(context.Background())
}
