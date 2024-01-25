package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var defaultSearchSortPaginationQuery = requests.SearchSortPaginationQuery{
	Limit:  types.Ptr[int32](1000),
	Offset: types.Ptr[int32](0),
	SortBy: types.Ptr[string]("created_at"),
	Order:  types.Ptr[string]("desc"),
}

func TestSearchPropertyCombination(t *testing.T) {
	// prepare testcases
	testCases := []struct {
		dataset     []dto.CreateProperty
		name        string
		query       dto.SearchPropertyCombinationQuery
		checkResult func(
			t *testing.T,
			ps []*model.PropertyModel,
			res *dto.SearchPropertyCombinationResponse,
			err error,
		)
	}{
		{
			name: "test exact fields",
			dataset: func() []dto.CreateProperty {
				dataset := []dto.CreateProperty{prepareRandomProperty(t, testAuthRepo)}
				dataset[0].City = "HCM"
				dataset = append(dataset, dataset[0])
				return dataset
			}(),
			query: dto.SearchPropertyCombinationQuery{
				SearchSortPaginationQuery: defaultSearchSortPaginationQuery,
				SearchPropertyQuery: dto.SearchPropertyQuery{
					PCity: types.Ptr[string]("HCM"),
				},
			},
			checkResult: func(
				t *testing.T,
				ps []*model.PropertyModel,
				res *dto.SearchPropertyCombinationResponse,
				err error,
			) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.GreaterOrEqual(t, res.Count, uint32(2))
				require.GreaterOrEqual(t, len(res.Items), 2)
				require.ElementsMatch(t, []uuid.UUID{ps[0].ID, ps[1].ID}, []uuid.UUID{res.Items[0].LId, res.Items[1].LId})
			},
		},
		{
			name: "test FK fields",
			dataset: func() []dto.CreateProperty {
				dataset := []dto.CreateProperty{prepareRandomProperty(t, testAuthRepo)}
				dataset[0].Features[0].FeatureID = 2
				dataset[0].Features[1].FeatureID = 5
				dataset[0].Features[2].FeatureID = 6
				dataset = append(dataset, dataset[0])
				return dataset
			}(),
			query: dto.SearchPropertyCombinationQuery{
				SearchSortPaginationQuery: defaultSearchSortPaginationQuery,
				SearchPropertyQuery: dto.SearchPropertyQuery{
					PFeatures: []int32{2, 5},
				},
			},
			checkResult: func(
				t *testing.T,
				ps []*model.PropertyModel,
				res *dto.SearchPropertyCombinationResponse,
				err error,
			) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.GreaterOrEqual(t, res.Count, uint32(2))
				require.GreaterOrEqual(t, len(res.Items), 2)
				require.ElementsMatch(t, []uuid.UUID{ps[0].ID, ps[1].ID}, []uuid.UUID{res.Items[0].LId, res.Items[1].LId})
			},
		},
		{
			name: "test ILIKE fields",
			dataset: func() []dto.CreateProperty {
				dataset := []dto.CreateProperty{prepareRandomProperty(t, testAuthRepo)}
				dataset[0].Name = "testproperty1"
				dataset = append(dataset, func() dto.CreateProperty {
					p := dataset[0]
					p.Name = "testproperty2"
					return p
				}())
				return dataset
			}(),
			query: dto.SearchPropertyCombinationQuery{
				SearchSortPaginationQuery: defaultSearchSortPaginationQuery,
				SearchPropertyQuery: dto.SearchPropertyQuery{
					PName: types.Ptr[string]("eStProP"),
				},
			},
			checkResult: func(
				t *testing.T,
				ps []*model.PropertyModel,
				res *dto.SearchPropertyCombinationResponse,
				err error,
			) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.GreaterOrEqual(t, res.Count, uint32(2))
				require.GreaterOrEqual(t, len(res.Items), 2)
				require.ElementsMatch(t, []uuid.UUID{ps[0].ID, ps[1].ID}, []uuid.UUID{res.Items[0].LId, res.Items[1].LId})
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ps := make([]*model.PropertyModel, 0, len(tc.dataset))
			for _, d := range tc.dataset {
				ps = append(ps, newRandomPropertyFromArg(t, testPropertyRepo, &d))
			}

			res, err := testPropertyRepo.SearchPropertyCombination(
				context.Background(),
				&tc.query,
			)

			tc.checkResult(t, ps, res, err)
		})
	}
}
