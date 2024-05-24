package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreatePaymentItem struct {
	Name     string `json:"name" validate:"required"`
	Price    int64  `json:"price" validate:"required,gte=0"`
	Quantity int32  `json:"quantity" validate:"required,gte=0"`
	Discount int32  `json:"discount" validate:"required"`
}

type CreatePayment struct {
	UserId    uuid.UUID `json:"userId" validate:"required,uuid4"`
	OrderId   string    `json:"orderId" validate:"required"`
	OrderInfo string    `json:"orderInfo" validate:"required"`
	Amount    int64     `json:"amount" validate:"required"`

	Items []CreatePaymentItem `json:"items" validate:"required,dive"`
}

type GetPaymentsOfUserQuery struct {
	Limit  *int32 `query:"limit" validate:"omitempty,gte=0"`
	Offset *int32 `query:"offset" validate:"omitempty,gte=0"`
}

type UpdatePayment struct {
	ID        int64                   `json:"id" validate:"required"`
	OrderId   *string                 `json:"orderId" validate:"omitempty"`
	OrderInfo *string                 `json:"orderInfo" validate:"omitempty"`
	Amount    *int64                  `json:"amount" validate:"omitempty,gte=0"`
	Status    *database.PAYMENTSTATUS `json:"status" validate:"omitempty"`
}
