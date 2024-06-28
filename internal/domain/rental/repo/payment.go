package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreateRentalPayment(ctx context.Context, data *dto.CreateRentalPayment) (rental_model.RentalPayment, error) {
	res, err := r.dao.CreateRentalPayment(ctx, data.ToCreateRentalPaymentDB())
	if err != nil {
		return rental_model.RentalPayment{}, err
	}
	return rental_model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayment(ctx context.Context, id int64) (rental_model.RentalPayment, error) {
	res, err := r.dao.GetRentalPayment(ctx, id)
	if err != nil {
		return rental_model.RentalPayment{}, err
	}
	return rental_model.ToRentalPaymentModel(&res), nil
}

func (r *repo) GetRentalPayments(ctx context.Context, ids []int64) ([]rental_model.RentalPayment, error) {
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select("id", "code", "rental_id", "created_at", "updated_at", "start_date", "end_date", "expiry_date", "payment_date", "updated_by", "status", "amount", "discount", "note")
	ib.From("rental_payments")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	query, args := ib.Build()
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items []rental_model.RentalPayment
		i     database.RentalPayment
	)
	for rows.Next() {
		if err := rows.Scan(&i.ID, &i.Code, &i.RentalID, &i.CreatedAt, &i.UpdatedAt, &i.ExpiryDate, &i.PaymentDate, &i.UpdatedBy, &i.Status, &i.Amount, &i.Note); err != nil {
			return nil, err
		}
		items = append(items, rental_model.ToRentalPaymentModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repo) GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]rental_model.RentalPayment, error) {
	res, err := r.dao.GetPaymentsOfRental(ctx, rentalID)
	if err != nil {
		return nil, err
	}

	var (
		rms []rental_model.RentalPayment
		rm  rental_model.RentalPayment
	)
	for i := range res {
		rm = rental_model.ToRentalPaymentModel(&res[i])
		rms = append(rms, rm)
	}
	return rms, nil
}

func (r *repo) GetManagedRentalPayments(ctx context.Context, uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) (res []rental_model.RentalPayment, err error) {
	subSB1 := sqlbuilder.PostgreSQL.NewSelectBuilder()
	subSB1.Select("1")
	subSB1.From("property_managers")
	subSB1.Where(
		"rentals.property_id = property_managers.property_id",
		subSB1.Equal("property_managers.manager_id", uid),
	)
	subSB := sqlbuilder.PostgreSQL.NewSelectBuilder()
	subSB.Select("1")
	subSB.From("rentals")
	subSB.Where(
		"rental_payments.rental_id = rentals.id",
		subSB.Exists(subSB1),
	)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "code", "rental_id", "created_at", "updated_at", "start_date", "end_date", "expiry_date", "payment_date", "updated_by", "status", "amount", "discount", "note")
	sb.From("rental_payments")
	sb.Where(
		sb.In("rental_payments.status", sqlbuilder.List(query.Status)),
		sb.Exists(subSB),
	)
	sb.Limit(int(*query.Limit))
	sb.Offset(int(*query.Offset))

	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []rental_model.RentalPayment
	for rows.Next() {
		var i database.RentalPayment
		if err = rows.Scan(
			&i.ID,
			&i.Code,
			&i.RentalID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StartDate,
			&i.EndDate,
			&i.ExpiryDate,
			&i.PaymentDate,
			&i.UpdatedBy,
			&i.Status,
			&i.Amount,
			&i.Discount,
			&i.Note,
		); err != nil {
			return nil, err
		}
		payments = append(payments, rental_model.ToRentalPaymentModel(&i))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
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

func (r *repo) UpdateFinePayments(ctx context.Context) error {
	return r.dao.UpdateFinePayments(ctx)
}

func (r *repo) UpdateFinePaymentsOfRental(ctx context.Context, rentalId int64) error {
	return r.dao.UpdateFinePaymentsOfRental(ctx, rentalId)
}
