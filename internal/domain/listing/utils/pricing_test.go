package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateListingPrice(t *testing.T) {
	testcases := []struct {
		name         string
		priority     int
		postDuration int
		checkResult  func(int64, error)
	}{
		{
			name:         "Test case 1",
			priority:     1,
			postDuration: 7,
			// expected:     14000,
			checkResult: func(result int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(14000), result)
			},
		},
		{
			name:         "invalid postduration",
			priority:     2,
			postDuration: 17,
			// expected:     85000,
			checkResult: func(result int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(85000), result)
			},
		},
		{
			name:         "invalid priority",
			priority:     5,
			postDuration: 7,
			// expected:     60000,
			checkResult: func(result int64, err error) {
				require.ErrorIs(t, err, ErrInvalidPriority)
				require.Equal(t, int64(0), result)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			result, err := CalculateListingPrice(tc.priority, tc.postDuration)
			tc.checkResult(result, err)
		})
	}

}
