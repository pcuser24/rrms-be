package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

type CreateRentalPayment struct {
	Code        string                       `json:"code" validate:"required"`
	RentalID    int64                        `json:"rentalId" validate:"required"`
	PaymentDate time.Time                    `json:"paymentDate" validate:"omitempty"`
	UserID      uuid.UUID                    `json:"userId" validate:"required"`
	Status      database.RENTALPAYMENTSTATUS `json:"status" validate:"omitempty"`
	Amount      float32                      `json:"amount" validate:"required"`
	Discount    *float32                     `json:"discount" validate:"omitempty,gte=0"`
	Note        *string                      `json:"note" validate:"omitempty"`
	StartDate   time.Time                    `json:"startDate" validate:"required"`
	EndDate     time.Time                    `json:"endDate" validate:"required"`
}

func (c *CreateRentalPayment) ToCreateRentalPaymentDB() database.CreateRentalPaymentParams {
	return database.CreateRentalPaymentParams{
		Code:     c.Code,
		RentalID: c.RentalID,
		PaymentDate: pgtype.Date{
			Time:  c.PaymentDate,
			Valid: !c.PaymentDate.IsZero(),
		},
		UserID: pgtype.UUID{
			Bytes: c.UserID,
			Valid: c.UserID != uuid.Nil,
		},
		Status: database.NullRENTALPAYMENTSTATUS{
			RENTALPAYMENTSTATUS: c.Status,
			Valid:               c.Status != "",
		},
		Amount:   c.Amount,
		Discount: types.Float32N(c.Discount),
		Note:     types.StrN(c.Note),
		StartDate: pgtype.Date{
			Time:  c.StartDate,
			Valid: !c.StartDate.IsZero(),
		},
		EndDate: pgtype.Date{
			Time:  c.EndDate,
			Valid: !c.EndDate.IsZero(),
		},
	}
}

type IUpdateRentalPayment interface {
	d()
}

type UpdateRentalPayment struct {
	ID          int64                        `json:"id" validate:"required"`
	Status      database.RENTALPAYMENTSTATUS `json:"status"`
	Note        *string                      `json:"note"`
	Amount      *float32                     `json:"amount"`
	Discount    *float32                     `json:"discount" validate:"omitempty,gte=0"`
	ExpiryDate  time.Time                    `json:"expiryDate"`
	PaymentDate time.Time                    `json:"paymentDate"`
	UserID      uuid.UUID                    `json:"userId"`
}

func (u *UpdateRentalPayment) d() {}

func (u *UpdateRentalPayment) ToUpdateRentalPaymentDB() database.UpdateRentalPaymentParams {
	return database.UpdateRentalPaymentParams{
		ID: u.ID,
		Status: database.NullRENTALPAYMENTSTATUS{
			RENTALPAYMENTSTATUS: u.Status,
			Valid:               u.Status != "",
		},
		Note:     types.StrN(u.Note),
		Amount:   types.Float32N(u.Amount),
		Discount: types.Float32N(u.Discount),
		ExpiryDate: pgtype.Date{
			Time:  u.ExpiryDate,
			Valid: !u.ExpiryDate.IsZero(),
		},
		PaymentDate: pgtype.Date{
			Time:  u.PaymentDate,
			Valid: !u.PaymentDate.IsZero(),
		},
		UserID: pgtype.UUID{
			Bytes: u.UserID,
			Valid: u.UserID != uuid.Nil,
		},
	}
}

type UpdatePlanRentalPayment struct {
	Amount     float32                      `json:"amount" validate:"omitempty"`
	Discount   *float32                     `json:"discount" validate:"omitempty"`
	ExpiryDate time.Time                    `json:"expiryDate" validate:"omitempty"`
	Status     database.RENTALPAYMENTSTATUS `json:"status" validate:"required,oneof=ISSUED PAID CANCELLED"`
}

func (u *UpdatePlanRentalPayment) d() {}

func (u UpdatePlanRentalPayment) Validate() error {
	if errs := validation.ValidateStruct(nil, u); len(errs) > 0 {
		return errs[0]
	}
	if u.Status == database.RENTALPAYMENTSTATUSISSUED &&
		(u.Amount == 0 || u.ExpiryDate.IsZero()) {
		return errors.New("amount and expiry date must be set when status is ISSUED")
	}

	return nil
}

type UpdateIssuedRentalPayment struct {
	Note   *string                      `json:"amount" validate:"omitempty"`
	Status database.RENTALPAYMENTSTATUS `json:"status" validate:"required,oneof=PENDING PLAN"`
}

func (u *UpdateIssuedRentalPayment) d() {}

type UpdatePendingRentalPayment struct {
	PaymentDate time.Time                    `json:"paymentDate" validate:"required"`
	Status      database.RENTALPAYMENTSTATUS `json:"status" validate:"required,oneof=REQUEST2PAY PAID"`
}

func (u *UpdatePendingRentalPayment) d() {}
