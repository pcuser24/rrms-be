package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Repo interface {
	CreatePayment(ctx context.Context, data *dto.CreatePayment) (*model.PaymentModel, error)
	GetPaymentsOfUser(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]model.PaymentModel, error)
	GetPaymentById(ctx context.Context, id int64) (*model.PaymentModel, error)
	UpdatePayment(ctx context.Context, data *dto.UpdatePayment) error
	CheckPaymentAccessible(ctx context.Context, userId uuid.UUID, id int64) (bool, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) CreatePayment(ctx context.Context, data *dto.CreatePayment) (*model.PaymentModel, error) {
	p, err := r.dao.CreatePayment(ctx, database.CreatePaymentParams{
		UserID:    data.UserId,
		OrderID:   data.OrderId,
		OrderInfo: data.OrderInfo,
		Amount:    data.Amount,
	})
	if err != nil {
		return nil, err
	}
	payment := model.ToPaymentModel(&p)

	for _, item := range data.Items {
		i, err := r.dao.CreatePaymentItem(ctx, database.CreatePaymentItemParams{
			PaymentID: p.ID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Discount:  item.Discount,
		})
		if err != nil {
			_ = r.dao.DeletePayment(ctx, p.ID)
			return nil, err
		}
		payment.Items = append(payment.Items, model.PaymentItemModel(i))
	}

	return payment, nil
}

func (r *repo) GetPaymentsOfUser(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]model.PaymentModel, error) {
	payments, err := r.dao.GetPaymentsOfUser(ctx, database.GetPaymentsOfUserParams{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	var result []model.PaymentModel
	for _, p := range payments {
		payment := model.ToPaymentModel(&p)

		items, err := r.dao.GetPaymentItemsByPaymentId(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			payment.Items = append(payment.Items, model.PaymentItemModel(item))
		}

		result = append(result, *payment)
	}
	return result, nil
}

func (r *repo) GetPaymentById(ctx context.Context, id int64) (*model.PaymentModel, error) {
	p, err := r.dao.GetPaymentById(ctx, id)
	if err != nil {
		return nil, err
	}
	payment := model.ToPaymentModel(&p)

	items, err := r.dao.GetPaymentItemsByPaymentId(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		payment.Items = append(payment.Items, model.PaymentItemModel(item))
	}

	return payment, nil
}

func (r *repo) UpdatePayment(ctx context.Context, data *dto.UpdatePayment) error {
	params := database.UpdatePaymentParams{
		ID:        data.ID,
		OrderID:   types.StrN(data.OrderId),
		OrderInfo: types.StrN(data.OrderInfo),
		Amount:    types.Int64N(data.Amount),
	}
	if data.Status != nil {
		params.Status = database.NullPAYMENTSTATUS{
			PAYMENTSTATUS: *data.Status,
			Valid:         true,
		}
	}
	return r.dao.UpdatePayment(ctx, params)
}

func (r *repo) CheckPaymentAccessible(ctx context.Context, userId uuid.UUID, id int64) (bool, error) {
	return r.dao.CheckPaymentAccessible(ctx, database.CheckPaymentAccessibleParams{
		UserID: userId,
		ID:     id,
	})
}
