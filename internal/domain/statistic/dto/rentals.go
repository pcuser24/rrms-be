package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

type RentalStatisticQuery struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
}

type RentalStatisticResponse struct {
	NewMaintenancesThisMonth []int64 `json:"newMaintenancesThisMonth"`
	NewMaintenancesLastMonth []int64 `json:"newMaintenancesLastMonth"`
}

type RentalPaymentStatisticQuery struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
	Limit     int32     `query:"limit"`
	Offset    int32     `query:"offset"`
}

type RentalPaymentIncomeItem struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Amount    float32   `json:"amount"`
}

type RentalPayment struct {
	model.RentalPayment
	ExpiryDuration int32     `json:"expiryDuration"`
	TenantId       uuid.UUID `json:"tenantId"`
	TenantName     string    `json:"tenantName"`
	PropertyID     uuid.UUID `json:"propertyId"`
	UnitID         uuid.UUID `json:"unitId"`
}

type RentalPaymentArrearsItem struct {
	StartTime time.Time       `json:"startTime"`
	EndTime   time.Time       `json:"endTime"`
	Payments  []RentalPayment `json:"payments"`
}
