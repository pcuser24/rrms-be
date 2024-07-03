package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	statistic_dto "github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Repo interface {
	GetManagedPropertiesByRole(ctx context.Context, userId uuid.UUID, role string) ([]uuid.UUID, error)
	GetPropertiesWithActiveListing(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetOccupiedProperties(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetManagedUnits(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetOccupiedUnits(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetMostRentedProperties(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedPropertyItem, error)
	GetMostRentedUnits(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedUnitItem, error)
	GetLeastRentedProperties(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedPropertyItem, error)
	GetLeastRentedUnits(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedUnitItem, error)
	GetNewApplications(ctx context.Context, userId uuid.UUID, month time.Time) ([]int64, error)
	GetRentalPaymentArrears(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) ([]statistic_dto.RentalPayment, error)
	GetRentalPaymentIncomes(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) (float32, error)
	GetMaintenanceRequests(ctx context.Context, userId uuid.UUID, month time.Time) ([]int64, error)
	GetPaymentsStatistic(ctx context.Context, userId uuid.UUID, query statistic_dto.PaymentsStatisticQuery) (float32, error)
	GetRecentListings(ctx context.Context, limit int32) ([]uuid.UUID, error)
	GetTotalTenantPendingPayments(ctx context.Context, userId uuid.UUID) (float32, error)
	GetTenantPendingPayments(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) ([]statistic_dto.RentalPayment, error)
	GetTenantExpenditure(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) (float32, error)
	GetTotalTenantsManagedByUserStatistic(ctx context.Context, userId uuid.UUID, query *statistic_dto.RentalStatisticQuery) (int32, error)
	GetTotalTenantsOfUnitStatistic(ctx context.Context, unitId uuid.UUID) (int32, error)
	GetRentalComplaintStatistics(ctx context.Context, userId uuid.UUID, status database.RENTALCOMPLAINTSTATUS) (int64, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) GetManagedPropertiesByRole(ctx context.Context, userId uuid.UUID, role string) ([]uuid.UUID, error) {
	return r.dao.GetManagedPropertiesByRole(ctx, database.GetManagedPropertiesByRoleParams{
		ManagerID: userId,
		Role:      role,
	})
}

func (r *repo) GetPropertiesWithActiveListing(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return r.dao.GetPropertiesWithActiveListing(ctx, userId)
}

func (r *repo) GetOccupiedProperties(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return r.dao.GetOccupiedProperties(ctx, userId)
}

func (r *repo) GetMostRentedProperties(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedPropertyItem, error) {
	resDB, err := r.dao.GetMostRentedProperties(ctx, database.GetMostRentedPropertiesParams{
		ManagerID: userId,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	res := make([]statistic_dto.ExtremelyRentedPropertyItem, len(resDB))
	for i, v := range resDB {
		res[i] = statistic_dto.ExtremelyRentedPropertyItem{
			PropertyID: v.ID,
			Count:      v.Count,
		}
	}
	return res, nil
}

func (r *repo) GetLeastRentedProperties(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedPropertyItem, error) {
	resDB, err := r.dao.GetLeastRentedProperties(ctx, database.GetLeastRentedPropertiesParams{
		ManagerID: userId,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	res := make([]statistic_dto.ExtremelyRentedPropertyItem, len(resDB))
	for i, v := range resDB {
		res[i] = statistic_dto.ExtremelyRentedPropertyItem{
			PropertyID: v.ID,
			Count:      v.Count,
		}
	}
	return res, nil
}

func (r *repo) GetMostRentedUnits(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedUnitItem, error) {

	resDB, err := r.dao.GetMostRentedUnits(ctx, database.GetMostRentedUnitsParams{
		ManagerID: userId,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	res := make([]statistic_dto.ExtremelyRentedUnitItem, len(resDB))
	for i, v := range resDB {
		res[i] = statistic_dto.ExtremelyRentedUnitItem{
			PropertyID: v.PropertyID,
			UnitID:     v.ID,
			Count:      v.Count,
		}
	}
	return res, nil
}

func (r *repo) GetLeastRentedUnits(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]statistic_dto.ExtremelyRentedUnitItem, error) {

	resDB, err := r.dao.GetLeastRentedUnits(ctx, database.GetLeastRentedUnitsParams{
		ManagerID: userId,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	res := make([]statistic_dto.ExtremelyRentedUnitItem, len(resDB))
	for i, v := range resDB {
		res[i] = statistic_dto.ExtremelyRentedUnitItem{
			PropertyID: v.PropertyID,
			UnitID:     v.ID,
			Count:      v.Count,
		}
	}
	return res, nil
}

func (r *repo) GetManagedUnits(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return r.dao.GetManagedUnits(ctx, userId)
}

// GetNewApplications returns the new applications of the given user in the given month.
func (r *repo) GetNewApplications(ctx context.Context, userId uuid.UUID, month time.Time) ([]int64, error) {
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, &time.Location{})

	subSB := sqlbuilder.NewSelectBuilder()
	subSB.Select("1").
		From("property_managers").
		Where(
			subSB.Equal("property_managers.manager_id", userId),
			"property_managers.property_id = applications.property_id",
		)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id").
		From("applications").
		Where(
			sb.Exists(subSB),
			sb.Equal("DATE_TRUNC('month', applications.created_at)", pgtype.Date{Time: month, Valid: !month.IsZero()}),
		)

	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *repo) GetOccupiedUnits(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	return r.dao.GetOccupiedUnits(ctx, userId)
}

// GetMaintenanceRequests returns the maintenance requests of the given user in the given month.
func (r *repo) GetMaintenanceRequests(ctx context.Context, userId uuid.UUID, month time.Time) ([]int64, error) {
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, &time.Location{})

	subSB1 := sqlbuilder.PostgreSQL.NewSelectBuilder()
	subSB1.Select("1").
		From("property_managers").
		Where(
			subSB1.Equal("property_managers.manager_id", userId),
			"property_managers.property_id = rentals.property_id",
		)

	subSB := sqlbuilder.PostgreSQL.NewSelectBuilder()
	subSB.Select("1").
		From("rentals").
		Where(
			"rental_complaints.rental_id = rentals.id",
			subSB.Exists(subSB1),
		)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id").
		From("rental_complaints").
		Where(
			sb.Exists(subSB),
			sb.Equal("DATE_TRUNC('month', rental_complaints.created_at)", pgtype.Date{Time: month, Valid: !month.IsZero()}),
		)

	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *repo) GetRentalPaymentArrears(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) ([]statistic_dto.RentalPayment, error) {
	res, err := r.dao.GetRentalPaymentArrears(ctx, database.GetRentalPaymentArrearsParams{
		ManagerID: userId,
		StartDate: pgtype.Date{
			Time:  query.StartTime,
			Valid: !query.StartTime.IsZero(),
		},
		EndDate: pgtype.Date{
			Time:  query.EndTime,
			Valid: !query.EndTime.IsZero(),
		},
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]statistic_dto.RentalPayment, len(res))
	for i, v := range res {
		items[i] = statistic_dto.RentalPayment{
			RentalPayment: model.RentalPayment{
				ID:          v.ID,
				Code:        v.Code,
				RentalID:    v.RentalID,
				CreatedAt:   v.CreatedAt,
				UpdatedAt:   v.UpdatedAt,
				StartDate:   v.StartDate.Time,
				EndDate:     v.EndDate.Time,
				ExpiryDate:  v.ExpiryDate.Time,
				PaymentDate: v.PaymentDate.Time,
				UpdatedBy:   v.UpdatedBy.Bytes,
				Status:      v.Status,
				Amount:      v.Amount,
				Discount:    types.PNFloat32(v.Discount),
				Note:        types.PNStr(v.Note),
			},
			ExpiryDuration: v.ExpiryDuration,
			TenantId:       v.TenantID.Bytes,
			TenantName:     v.TenantName,
			PropertyID:     v.PropertyID,
			UnitID:         v.UnitID,
		}
	}
	return items, nil
}

func (r *repo) GetRentalPaymentIncomes(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) (float32, error) {
	res, err := r.dao.GetRentalPaymentIncomes(ctx, database.GetRentalPaymentIncomesParams{
		ManagerID: userId,
		StartDate: pgtype.Date{
			Time:  query.StartTime,
			Valid: !query.StartTime.IsZero(),
		},
		EndDate: pgtype.Date{
			Time:  query.EndTime,
			Valid: !query.EndTime.IsZero(),
		},
	})
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *repo) GetPaymentsStatistic(ctx context.Context, userId uuid.UUID, query statistic_dto.PaymentsStatisticQuery) (float32, error) {
	res, err := r.dao.GetPaymentsStatistic(ctx, database.GetPaymentsStatisticParams{
		UserID:    userId,
		StartDate: query.StartTime,
		EndDate:   query.EndTime,
	})
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *repo) GetTotalTenantPendingPayments(ctx context.Context, userId uuid.UUID) (float32, error) {
	return r.dao.GetTotalTenantPendingPayments(ctx, pgtype.UUID{
		Bytes: userId,
		Valid: userId != uuid.Nil,
	})
}

func (r *repo) GetTenantPendingPayments(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) ([]statistic_dto.RentalPayment, error) {
	res, err := r.dao.GetTenantPendingPayments(ctx, database.GetTenantPendingPaymentsParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]statistic_dto.RentalPayment, len(res))
	for i, v := range res {
		items[i] = statistic_dto.RentalPayment{
			RentalPayment: model.RentalPayment{
				ID:          v.ID,
				Code:        v.Code,
				RentalID:    v.RentalID,
				CreatedAt:   v.CreatedAt,
				UpdatedAt:   v.UpdatedAt,
				StartDate:   v.StartDate.Time,
				EndDate:     v.EndDate.Time,
				ExpiryDate:  v.ExpiryDate.Time,
				PaymentDate: v.PaymentDate.Time,
				UpdatedBy:   v.UpdatedBy.Bytes,
				Status:      v.Status,
				Amount:      v.Amount,
				Discount:    types.PNFloat32(v.Discount),
				Note:        types.PNStr(v.Note),
			},
			ExpiryDuration: v.ExpiryDuration,
			TenantId:       v.TenantID.Bytes,
			TenantName:     v.TenantName,
			PropertyID:     v.PropertyID,
			UnitID:         v.UnitID,
		}
		servicedb, err := r.dao.GetRentalServicesByRentalID(ctx, items[i].RentalID)
		if err != nil {
			return nil, err
		}
		services := make([]model.RentalService, 0, len(servicedb))
		for _, item := range servicedb {
			services = append(services, model.ToRentalService(&item))
		}
		items[i].Services = services
	}
	return items, nil
}

func (r *repo) GetTenantExpenditure(ctx context.Context, userId uuid.UUID, query statistic_dto.RentalPaymentStatisticQuery) (float32, error) {
	res, err := r.dao.GetTenantExpenditure(ctx, database.GetTenantExpenditureParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		StartDate: pgtype.Date{
			Time:  query.StartTime,
			Valid: !query.StartTime.IsZero(),
		},
		EndDate: pgtype.Date{
			Time:  query.EndTime,
			Valid: !query.EndTime.IsZero(),
		},
	})
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *repo) GetRentalComplaintStatistics(ctx context.Context, userId uuid.UUID, status database.RENTALCOMPLAINTSTATUS) (int64, error) {
	return r.dao.GetRentalComplaintStatistics(ctx, database.GetRentalComplaintStatisticsParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Status: status,
	})
}

func (r *repo) GetTotalTenantsManagedByUserStatistic(ctx context.Context, userId uuid.UUID, query *statistic_dto.RentalStatisticQuery) (int32, error) {
	return r.dao.GetTotalTenantsManagedByUserStatistic(ctx, database.GetTotalTenantsManagedByUserStatisticParams{
		UserID: userId,
		StartTime: pgtype.Date{
			Time:  query.StartTime,
			Valid: !query.StartTime.IsZero(),
		},
		EndTime: pgtype.Date{
			Time:  query.EndTime,
			Valid: !query.EndTime.IsZero(),
		},
	})
}

func (r *repo) GetTotalTenantsOfUnitStatistic(ctx context.Context, unitId uuid.UUID) (int32, error) {
	return r.dao.GetTotalTenantsOfUnitStatistic(ctx, unitId)
}
