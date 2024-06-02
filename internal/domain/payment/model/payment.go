package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type PaymentItemModel struct {
	PaymentID int64   `json:"paymentId"`
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Quantity  int32   `json:"quantity"`
	Discount  int32   `json:"discount"`
}

type PaymentModel struct {
	ID        int64                  `json:"id"`
	UserID    uuid.UUID              `json:"userId"`
	OrderID   string                 `json:"orderId"`
	OrderInfo string                 `json:"orderInfo"`
	Amount    float32                `json:"amount"`
	Status    database.PAYMENTSTATUS `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`

	Items []PaymentItemModel `json:"items"`
}

func ToPaymentModel(p *database.Payment) *PaymentModel {
	return &PaymentModel{
		ID:        p.ID,
		UserID:    p.UserID,
		OrderID:   p.OrderID,
		OrderInfo: p.OrderInfo,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Items:     []PaymentItemModel{},
	}
}
