package dto

import "time"

type PaymentsStatisticQuery struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
}

type PaymentsStatisticItem = RentalPaymentIncomeItem
