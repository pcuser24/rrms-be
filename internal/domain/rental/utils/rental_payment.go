package utils

import (
	"fmt"
	"time"
)

// enum for rental payment type: "RENTAL", "ELECTRICITY", "WATER", "SERVICES"
type RentalPaymentType string

const (
	RENTALPAYMENTTYPERENTAL      RentalPaymentType = "RENTAL"
	RENTALPAYMENTTYPEDEPOSIT     RentalPaymentType = "DEPOSIT"
	RENTALPAYMENTTYPEELECTRICITY RentalPaymentType = "ELECTRICITY"
	RENTALPAYMENTTYPEWATER       RentalPaymentType = "WATER"
	RENTALPAYMENTTYPESERVICES    RentalPaymentType = "SERVICES"
)

func GetRentalPaymentCode(rentalID int64, paymentType RentalPaymentType, startDate, endDate time.Time) string {
	return fmt.Sprintf("%d_%s_%02d%d%02d%d", rentalID, paymentType, startDate.Month(), startDate.Year(), endDate.Month(), endDate.Year())
}

func GetRentalPaymentPrice(startDate, endDate time.Time, paymentBasis int32, rentalPrice float64) float64 {
	// calculate rental duration in days
	rentalDuration := endDate.Sub(startDate).Hours() / 24

	// calculate rental price based on number of days
	if startDate.AddDate(0, int(paymentBasis), 0).After(endDate) {
		// prorate rental price
		return rentalPrice * rentalDuration / float64(paymentBasis*30)
	} else {
		// full rental price
		return rentalPrice
	}
}
