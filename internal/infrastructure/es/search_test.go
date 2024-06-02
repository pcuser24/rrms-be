package es

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionscoremode"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestSearch(t *testing.T) {
	listingRepo := listing_repo.NewRepo(dao)
	propertyRepo := property_repo.NewRepo(dao)
	unitRepo := unit_repo.NewRepo(dao)

	listing, err := listingRepo.GetListingByID(context.Background(), uuid.MustParse("00808905-a3c4-4c5c-b989-15e6341bceca"))
	require.NoError(t, err)
	property, err := propertyRepo.GetPropertyById(context.Background(), listing.PropertyID)
	require.NoError(t, err)
	unitIds := make([]uuid.UUID, 0, len(listing.Units))
	for _, unit := range listing.Units {
		unitIds = append(unitIds, unit.UnitID)
	}
	units, err := unitRepo.GetUnitsByIds(context.Background(), unitIds, []string{"area", "floor", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "amenities"})
	require.NoError(t, err)

	search := typedClient.Search().
		Index("listings").
		Request(&search.Request{
			Size: types.Ptr(10),
			Source_: estypes.SourceFilter{
				Includes: []string{"id"},
			},
			Query: &estypes.Query{
				FunctionScore: &estypes.FunctionScoreQuery{
					Query: &estypes.Query{
						Bool: &estypes.BoolQuery{
							Must: []estypes.Query{
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
							},
						},
					},
					Functions: []estypes.FunctionScore{
						{
							Filter: &estypes.Query{
								Term: map[string]estypes.TermQuery{
									"property.city": {
										Value: property.City,
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](8),
						},
						{
							Filter: &estypes.Query{
								Term: map[string]estypes.TermQuery{
									"property.district": {
										Value: property.District,
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](4),
						},
						{
							Filter: &estypes.Query{
								Term: map[string]estypes.TermQuery{
									"property.ward": {
										Value: property.Ward,
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](2),
						},
						{
							Filter: &estypes.Query{
								Range: map[string]estypes.RangeQuery{
									"price": estypes.NumberRangeQuery{
										Gte: types.Ptr[estypes.Float64](estypes.Float64(listing.Price * 0.8)),
										Lte: types.Ptr[estypes.Float64](estypes.Float64(listing.Price * 1.2)),
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](16),
						},
						{
							Filter: &estypes.Query{
								Range: map[string]estypes.RangeQuery{
									"property.area": estypes.NumberRangeQuery{
										Gte: types.Ptr[estypes.Float64](estypes.Float64(property.Area * 0.8)),
										Lte: types.Ptr[estypes.Float64](estypes.Float64(property.Area * 1.2)),
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](8),
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
												Gte: types.Ptr[estypes.Float64](estypes.Float64(
													utils.PtrDerefence(units[0].NumberOfBedrooms, 1)),
												),
											},
										},
									},
								},
							},
							Weight: types.Ptr[estypes.Float64](2),
						},
					},
					ScoreMode: &functionscoremode.Sum,
					BoostMode: &functionboostmode.Replace,
				},
			},
		},
		)

	searchRes, err := search.Do(context.Background())
	require.NoError(t, err)
	hits := make([]dto.SuggestedListing, 0, len(searchRes.Hits.Hits))
	for _, h := range searchRes.Hits.Hits {
		var i dto.SuggestedListing
		err := json.Unmarshal(h.Source_, &i)
		require.NoError(t, err)
		hits = append(hits, i)
	}

	t.Log(hits)
}

func TestSearch2(t *testing.T) {
	listingRepo := listing_repo.NewRepo(dao)
	propertyRepo := property_repo.NewRepo(dao)
	unitRepo := unit_repo.NewRepo(dao)

	listing, err := listingRepo.GetListingByID(context.Background(), uuid.MustParse("00808905-a3c4-4c5c-b989-15e6341bceca"))
	require.NoError(t, err)
	property, err := propertyRepo.GetPropertyById(context.Background(), listing.PropertyID)
	require.NoError(t, err)
	unitIds := make([]uuid.UUID, 0, len(listing.Units))
	for _, unit := range listing.Units {
		unitIds = append(unitIds, unit.UnitID)
	}
	units, err := unitRepo.GetUnitsByIds(context.Background(), unitIds, []string{"area", "floor", "number_of_living_rooms", "number_of_bedrooms", "number_of_bathrooms", "number_of_toilets", "number_of_balconies", "number_of_kitchens", "amenities"})
	require.NoError(t, err)

	query := fmt.Sprintf(`
	{
		"size": %d,
		"_source": ["id"],
		"query": {
			"function_score": {
				"query": {
					"bool": {
						"must": [
							{ "term": { "active": true } },
							{ "range": { "expired_at": { "gte": "now" } } }
						]
					}
				},
				"functions": [
					{
						"filter": {
							"term": {"property.city": "%s"}
						},
						"weight": 16
					},
					{
						"filter": {
							"term": {"property.district": "%s"}
						},
						"weight": 8
					},
					{
						"filter": {
							"term": {"property.ward": "%s"}
						},
						"weight": 4
					},
					{
						"filter": {
							"range": {
								"price": {
									"gte": %f,
									"lte": %f
								}
							}
						},
						"weight": 16
					},
					{
						"filter": {
							"range": {
								"property.area": {
									"gte": %f,
									"lte": %f
								}
							}
						},
						"weight": 8
					},
					{
						"filter": {
							"term": {
								"property.type": "%s"
							}
						},
						"weight": 4
					},
					{
						"filter": {
							"nested": {
								"path": "listing_units",
								"query": {
									"range": {
										"listing_units.number_of_bedrooms": {
											"gte": %d
										}
									}
								}
							}
						},
						"weight": 2
					},
					{
						"filter": {
							"bool": {
								"should": [
								]
							}
						},
						"weight": 1 
					}
				],
				"score_mode": "sum",
				"boost_mode": "replace"
			}
		}
	}	
	`,
		10, property.City, property.District, utils.PtrDerefence(property.Ward, ""), listing.Price*0.8, listing.Price*1.2, property.Area*0.8, property.Area*1.2, property.Type, utils.PtrDerefence(units[0].NumberOfBedrooms, utils.PtrDerefence(units[0].NumberOfBedrooms, 0)),
	)
	t.Log(query)

	res, err := client.Search(
		client.Search.WithIndex("listings"),
		client.Search.WithBody(strings.NewReader(query)),
	)
	require.NoError(t, err)
	defer res.Body.Close()

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	require.NoError(t, err)
	t.Log(r)
}
