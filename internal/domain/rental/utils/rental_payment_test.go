package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func TestGetRentalPaymentCode(t *testing.T) {
	var rentalID int64 = 123456789
	var paymentType RentalPaymentType = RENTALPAYMENTTYPERENTAL
	startDate, err := time.Parse("2006-01-02", "2021-01-01")
	require.NoError(t, err)
	endDate, err := time.Parse("2006-01-02", "2021-02-01")
	require.NoError(t, err)

	code := GetRentalPaymentCode(rentalID, paymentType, 0, startDate, endDate)
	require.Equal(t, "123456789_RENTAL_012021022021", code)

	startDate, err = time.Parse("2006-01-02", "2020-11-24")
	require.NoError(t, err)
	code = GetRentalPaymentCode(rentalID, paymentType, 0, startDate, endDate)
	require.Equal(t, "123456789_RENTAL_112020022021", code)

	serviceId := int64(987654321)
	code = GetRentalPaymentCode(rentalID, RENTALPAYMENTTYPESERVICE, serviceId, startDate, endDate)
	require.Equal(t, "123456789_SERVICE_987654321_112020022021", code)
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

func TestGetRentalServiceName(t *testing.T) {
	code := "123456789_RENTAL_012021022021"
	rServices := []rental_model.RentalService{}
	serviceName, err := GetServiceName(code, rServices)
	require.NoError(t, err)
	require.Equal(t, "Thuê nhà", serviceName)

	code = "123456789_SERVICE_987654321_112020022021"
	rServices = []rental_model.RentalService{
		{
			ID:   987654321,
			Name: "Random service name",
		},
	}
	serviceName, err = GetServiceName(code, rServices)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s Random service name", mapRentalPaymentTypeToServiceName[RENTALPAYMENTTYPESERVICE]), serviceName)
}
