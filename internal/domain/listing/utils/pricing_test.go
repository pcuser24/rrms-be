package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
)

func TestCalculateListingPrice(t *testing.T) {
	testcases := []struct {
		name         string
		priority     int
		postDuration int
		checkResult  func(*testing.T, float32, error)
	}{
		{
			name:         "OK",
			priority:     1,
			postDuration: 7,
			// expected:     14000,
			checkResult: func(t *testing.T, result float32, err error) {
				require.NoError(t, err)
				require.Equal(t, float32(14000), result)
			},
		},
		{
			name:         "invalid postduration",
			priority:     2,
			postDuration: 17,
			// expected:     85000,
			checkResult: func(t *testing.T, result float32, err error) {
				require.NoError(t, err)
				require.Equal(t, float32(170000), result)
			},
		},
		{
			name:         "invalid priority",
			priority:     5,
			postDuration: 7,
			// expected:     60000,
			checkResult: func(t *testing.T, result float32, err error) {
				require.ErrorIs(t, err, ErrInvalidPriority)
				require.Equal(t, float32(0), result)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			result, _, _, err := CalculateListingPrice(tc.priority, tc.postDuration)
			tc.checkResult(t, result, err)
		})
	}

}

func TestCalculateUpgradeListingPrice(t *testing.T) {
	testcases := []struct {
		name        string
		listing     model.ListingModel
		priority    int
		checkResult func(*testing.T, float32, error)
	}{
		{
			name: "OK",
			listing: model.ListingModel{
				Priority:  1,
				CreatedAt: time.Now().AddDate(0, 0, -10),
			},
			priority: 2,
			checkResult: func(t *testing.T, result float32, err error) {
				require.NoError(t, err)
				require.Equal(t, float32(80000), result)
			},
		},
		{
			name: "invalid priority",
			listing: model.ListingModel{
				Priority:  3,
				CreatedAt: time.Now().AddDate(0, 0, -10),
			},
			priority: 2,
			checkResult: func(t *testing.T, result float32, err error) {
				require.ErrorIs(t, err, ErrInvalidPriority)
				require.Equal(t, float32(0), result)
			},
		},
		{
			name: "invalid priority range",
			listing: model.ListingModel{
				Priority:  1,
				CreatedAt: time.Now().AddDate(0, 0, -10),
			},
			priority: 0,
			checkResult: func(t *testing.T, result float32, err error) {
				require.ErrorIs(t, err, ErrInvalidPriority)
				require.Equal(t, float32(0), result)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			result, _, err := CalculateUpgradeListingPrice(&tc.listing, tc.priority)
			tc.checkResult(t, result, err)
		})
	}
}

func TestCalculateExtendListingPrice(t *testing.T) {
	testcases := []struct {
		name        string
		listing     model.ListingModel
		duration    int
		checkResult func(*testing.T, float32, error)
	}{
		{
			name: "OK",
			listing: model.ListingModel{
				Priority:  1,
				CreatedAt: time.Now().AddDate(0, 0, -10),
			},
			duration: 7,
			checkResult: func(t *testing.T, result float32, err error) {
				require.NoError(t, err)
				require.Equal(t, float32(14000), result)
			},
		},
		{
			name: "invalid duration",
			listing: model.ListingModel{
				Priority:  1,
				CreatedAt: time.Now().AddDate(0, 0, -10),
			},
			duration: 0,
			checkResult: func(t *testing.T, result float32, err error) {
				require.Error(t, err)
				require.Equal(t, float32(0), result)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			result, _, err := CalculateExtendListingPrice(&tc.listing, tc.duration)
			tc.checkResult(t, result, err)
		})
	}
}
