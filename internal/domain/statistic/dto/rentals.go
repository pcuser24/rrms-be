package dto

import (
	"time"

	"github.com/google/uuid"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
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
	rental_model.RentalPayment
	ExpiryDuration int32                        `json:"expiryDuration"`
	TenantId       uuid.UUID                    `json:"tenantId"`
	TenantName     string                       `json:"tenantName"`
	PropertyID     uuid.UUID                    `json:"propertyId"`
	UnitID         uuid.UUID                    `json:"unitId"`
	Services       []rental_model.RentalService `json:"services"`
}

type RentalPaymentArrearsItem struct {
	StartTime time.Time       `json:"startTime"`
	EndTime   time.Time       `json:"endTime"`
	Payments  []RentalPayment `json:"payments"`
}

type TenantRentalStatisticResponse struct {
	CurrentRentals []rental_model.RentalModel `json:"currentRentals"`
	EndedRentalIds []int64                    `json:"nEndedRentalIds"`
}

type TenantMaintenanceStatisticResponse struct {
	Pending  int64 `json:"pending"`
	Resolved int64 `json:"resolved"`
	Closed   int64 `json:"closed"`
}

type TenantExpenditureStatisticQuery struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
	Limit     int32     `query:"limit"`
	Offset    int32     `query:"offset"`
}

type TenantExpenditureStatisticItem struct {
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Expenditure float32   `json:"expenditure"`
}

type TenantArrearsStatistic struct {
	Total    float32         `json:"total"`
	Payments []RentalPayment `json:"payments"`
}
