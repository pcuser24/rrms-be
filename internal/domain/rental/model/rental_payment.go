package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type RentalPayment struct {
	ID int64 `json:"id"`
	// {payment.id}_{ELECTRICITY | WATER | RENTAL | SERVICES{id}}_{payment.created_at}
	Code       string    `json:"code"`
	RentalID   int64     `json:"rentalId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	ExpiryDate time.Time `json:"expiryDate"`
	// the date the payment gets paid
	PaymentDate time.Time                    `json:"paymentDate"`
	UpdatedBy   uuid.UUID                    `json:"updatedBy"`
	Status      database.RENTALPAYMENTSTATUS `json:"status"`
	Amount      float32                      `json:"amount"`
	Discount    *float32                     `json:"discount"`
	Note        *string                      `json:"note"`
}

func ToRentalPaymentModel(prdb *database.RentalPayment) RentalPayment {
	return RentalPayment{
		ID:          prdb.ID,
		Code:        prdb.Code,
		RentalID:    prdb.RentalID,
		CreatedAt:   prdb.CreatedAt,
		UpdatedAt:   prdb.UpdatedAt,
		StartDate:   prdb.StartDate.Time,
		EndDate:     prdb.EndDate.Time,
		ExpiryDate:  prdb.ExpiryDate.Time,
		PaymentDate: prdb.PaymentDate.Time,
		UpdatedBy:   prdb.UpdatedBy.Bytes,
		Status:      prdb.Status,
		Amount:      prdb.Amount,
		Discount:    types.PNFloat32(prdb.Discount),
		Note:        types.PNStr(prdb.Note),
	}
}
