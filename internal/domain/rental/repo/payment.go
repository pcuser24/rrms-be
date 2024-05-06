package repo

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreateRentalPayment(ctx context.Context, data *dto.CreateRentalPayment) (model.RentalPayment, error) {
	res, err := r.dao.CreateRentalPayment(ctx, data.ToCreateRentalPaymentDB())
	if err != nil {
		return model.RentalPayment{}, err
	}
	return model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayment(ctx context.Context, id int64) (model.RentalPayment, error) {
	res, err := r.dao.GetRentalPayment(ctx, id)
	if err != nil {
		return model.RentalPayment{}, err
	}
	return model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayments(ctx context.Context, ids []int64) ([]model.RentalPayment, error) {
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select("id", "code", "rental_id", "created_at", "updated_at", "expiry_date", "payment_date", "updated_by", "status", "amount", "note")
	ib.From("rental_payments")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	query, args := ib.Build()
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []model.RentalPayment
		i     database.RentalPayment
	)
	for rows.Next() {
		if err := rows.Scan(&i.ID, &i.Code, &i.RentalID, &i.CreatedAt, &i.UpdatedAt, &i.ExpiryDate, &i.PaymentDate, &i.UpdatedBy, &i.Status, &i.Amount, &i.Note); err != nil {
			return nil, err
		}
		items = append(items, model.ToRentalPaymentModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repo) GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]model.RentalPayment, error) {
	res, err := r.dao.GetPaymentsOfRental(ctx, rentalID)
	if err != nil {
		return nil, err
	}

	var (
		rms []model.RentalPayment
		rm  model.RentalPayment
	)
	for i := range res {
		rm = model.ToRentalPaymentModel(&res[i])
		rms = append(rms, rm)
	}
	return rms, nil
}

func (r *repo) UpdateRentalPayment(ctx context.Context, data *dto.UpdateRentalPayment) error {
	return r.dao.UpdateRentalPayment(ctx, data.ToUpdateRentalPaymentDB())
}

func (r *repo) PlanRentalPayments(ctx context.Context) ([]int64, error) {
	return r.dao.PlanRentalPayments(ctx)
}

func (r *repo) PlanRentalPayment(ctx context.Context, rentalId int64) ([]int64, error) {
	return r.dao.PlanRentalPayment(ctx, rentalId)
}
