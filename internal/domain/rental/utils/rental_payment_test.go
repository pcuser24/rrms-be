package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetRentalPaymentCode(t *testing.T) {
	var rentalID int64 = 123456789
	var paymentType RentalPaymentType = RENTALPAYMENTTYPERENTAL
	startDate, err := time.Parse("2006-01-02", "2021-01-01")
	require.NoError(t, err)
	endDate, err := time.Parse("2006-01-02", "2021-02-01")
	require.NoError(t, err)

	code := GetRentalPaymentCode(rentalID, paymentType, startDate, endDate)
	require.Equal(t, "123456789_RENTAL_012021022021", code)

	startDate, err = time.Parse("2006-01-02", "2023-10-01")
	require.NoError(t, err)
	startDate, err = time.Parse("2006-01-02", "2023-11-24")
	require.NoError(t, err)
	code = GetRentalPaymentCode(rentalID, paymentType, startDate, endDate)
	require.Equal(t, "123456789_RENTAL_102023112023", code)
}

func TestGetRentalPaymentPrice(t *testing.T) {
	const (
		rentalPrice  = 1000.00
		defaultDelta = rentalPrice / 30
	)

	startDate, err := time.Parse("2006-01-02", "2021-01-01")
	require.NoError(t, err)
	endDate, err := time.Parse("2006-01-02", "2021-02-01")
	require.NoError(t, err)

	price := GetRentalPaymentPrice(startDate, endDate, 1, rentalPrice)
	require.Equal(t, rentalPrice, price)

	startDate, err = time.Parse("2006-01-02", "2021-01-15")
	require.NoError(t, err)
	price = GetRentalPaymentPrice(startDate, endDate, 1, rentalPrice)
	// Approximate value
	require.InDelta(t, 500.00, price, defaultDelta*2)
}
