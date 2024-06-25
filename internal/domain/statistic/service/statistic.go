package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error) {
	managedProperties, err := s.domainRepo.PropertyRepo.GetManagedProperties(context.Background(), userId, &property_dto.GetPropertiesQuery{})
	if err != nil {
		return
	}
	for _, p := range managedProperties {
		res.Properties = append(res.Properties, p.PropertyID)
	}

	res.OwnedProperties, err = s.domainRepo.StatisticRepo.GetManagedPropertiesByRole(context.Background(), userId, "OWNER")
	if err != nil {
		return
	}

	res.OccupiedProperties, err = s.domainRepo.StatisticRepo.GetOccupiedProperties(context.Background(), userId)
	if err != nil {
		return
	}

	res.Units, err = s.domainRepo.StatisticRepo.GetManagedUnits(context.Background(), userId)
	if err != nil {
		return
	}

	res.OccupiedUnits, err = s.domainRepo.StatisticRepo.GetOccupiedUnits(context.Background(), userId)
	if err != nil {
		return
	}

	res.PropertiesWithActiveListing, err = s.domainRepo.StatisticRepo.GetPropertiesWithActiveListing(context.Background(), userId)
	if err != nil {
		return
	}

	res.MostRentedProperties, err = s.domainRepo.StatisticRepo.GetMostRentedProperties(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.LeastRentedProperties, err = s.domainRepo.StatisticRepo.GetLeastRentedProperties(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.MostRentedUnits, err = s.domainRepo.StatisticRepo.GetMostRentedUnits(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.LeastRentedUnits, err = s.domainRepo.StatisticRepo.GetLeastRentedUnits(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	return
}

func (s *service) GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error) {
	now := time.Now()

	res.NewApplicationsThisMonth, err = s.domainRepo.StatisticRepo.GetNewApplications(context.Background(), userId, now)
	if err != nil {
		return
	}

	res.NewApplicationsLastMonth, err = s.domainRepo.StatisticRepo.GetNewApplications(context.Background(), userId, now.AddDate(0, -1, 0))
	return
}

func (s *service) GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error) {
	now := time.Now()
	res.NewMaintenancesThisMonth, err = s.domainRepo.StatisticRepo.GetMaintenanceRequests(context.Background(), userId, now)
	if err != nil {
		return
	}

	res.NewMaintenancesLastMonth, err = s.domainRepo.StatisticRepo.GetMaintenanceRequests(context.Background(), userId, now.AddDate(0, -1, 0))
	return
}

func (s *service) GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error) {
	user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), userId)
	if err != nil {
		return
	}
	if query.StartTime.Before(user.CreatedAt) {
		query.StartTime = user.CreatedAt
	}

	// for example: endDate = May 18th 2024, startDate = March 10th 2024, then derived intervals are [March 10th - April 1st], [April 2nd - May 1st], [May 2nd - May 18th]

	// Start from the first day of the start month
	current := query.StartTime

	for current.Before(query.EndTime) {
		// Calculate the end of the current month
		endOfMonth := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, query.StartTime.Location()).AddDate(0, 1, -1)

		// Choose the end of the current month or the end date, whichever is earlier
		intervalEnd := endOfMonth
		if endOfMonth.After(query.EndTime) {
			intervalEnd = query.EndTime
		}

		payments, err := s.domainRepo.StatisticRepo.GetRentalPaymentArrears(context.Background(), userId, dto.RentalPaymentStatisticQuery{
			StartTime: current,
			EndTime:   intervalEnd,
			Limit:     query.Limit,
			Offset:    query.Offset,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, dto.RentalPaymentArrearsItem{
			StartTime: current,
			EndTime:   intervalEnd,
			Payments:  payments,
		})

		// Move to the next month
		current = endOfMonth.AddDate(0, 0, 1)
	}

	return
}

func (s *service) GetRentalPaymentIncomes(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentIncomeItem, err error) {
	user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), userId)
	if err != nil {
		return
	}
	if query.StartTime.Before(user.CreatedAt) {
		query.StartTime = user.CreatedAt
	}

	// for example: endDate = May 18th 2024, startDate = March 10th 2024, then derived intervals are [March 10th - April 1st], [April 2nd - May 1st], [May 2nd - May 18th]

	// Start from the first day of the start month
	current := query.StartTime

	for current.Before(query.EndTime) {
		// Calculate the end of the current month
		endOfMonth := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, query.StartTime.Location()).AddDate(0, 1, -1)

		// Choose the end of the current month or the end date, whichever is earlier
		intervalEnd := endOfMonth
		if endOfMonth.After(query.EndTime) {
			intervalEnd = query.EndTime
		}

		income, err := s.domainRepo.StatisticRepo.GetRentalPaymentIncomes(context.Background(), userId, dto.RentalPaymentStatisticQuery{
			StartTime: current,
			EndTime:   intervalEnd,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, dto.RentalPaymentIncomeItem{
			StartTime: current,
			EndTime:   intervalEnd,
			Amount:    income,
		})

		// Move to the next month
		current = endOfMonth.AddDate(0, 0, 1)
	}

	return
}

func (s *service) GetPaymentsStatistic(userId uuid.UUID, query dto.PaymentsStatisticQuery) (res []dto.PaymentsStatisticItem, err error) {
	user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), userId)
	if err != nil {
		return
	}
	if query.StartTime.Before(user.CreatedAt) {
		query.StartTime = user.CreatedAt
	}

	// Start from the first day of the start month
	current := query.StartTime

	for current.Before(query.EndTime) {
		// Calculate the end of the current month
		endOfMonth := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, query.StartTime.Location()).AddDate(0, 1, -1)

		// Choose the end of the current month or the end date, whichever is earlier
		intervalEnd := endOfMonth
		if endOfMonth.After(query.EndTime) {
			intervalEnd = query.EndTime
		}

		payment, err := s.domainRepo.StatisticRepo.GetPaymentsStatistic(context.Background(), userId, dto.PaymentsStatisticQuery{
			StartTime: current,
			EndTime:   intervalEnd,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, dto.PaymentsStatisticItem{
			StartTime: current,
			EndTime:   intervalEnd,
			Amount:    payment,
		})

		// Move to the next month
		current = endOfMonth.AddDate(0, 0, 1)
	}
	return nil, nil
}

func (s *service) GetTenantRentalStatistic(userId uuid.UUID) (res dto.TenantRentalStatisticResponse, err error) {
	query := rental_dto.GetRentalsQuery{
		Limit:   types.Ptr[int32](math.MaxInt32),
		Offset:  types.Ptr[int32](0),
		Expired: false,
		Fields:  []string{"property_id,unit_id", "start_date", "rental_period"},
	}

	currentRentalIds, err := s.domainRepo.RentalRepo.GetMyRentals(context.Background(), userId, &query)
	if err != nil {
		return
	}
	query.Expired = true
	res.EndedRentalIds, err = s.domainRepo.RentalRepo.GetMyRentals(context.Background(), userId, &query)
	if err != nil {
		return
	}

	res.CurrentRentals, err = s.domainRepo.RentalRepo.GetRentalsByIds(context.Background(), currentRentalIds, query.Fields)
	if err != nil {
		return
	}

	return
}

func (s *service) GetTenantMaintenanceStatistic(userId uuid.UUID) (res dto.TenantMaintenanceStatisticResponse, err error) {
	res.Pending, err = s.domainRepo.StatisticRepo.GetRentalComplaintStatistics(context.Background(), userId, database.RENTALCOMPLAINTSTATUSPENDING)
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return
	}
	res.Resolved, err = s.domainRepo.StatisticRepo.GetRentalComplaintStatistics(context.Background(), userId, database.RENTALCOMPLAINTSTATUSRESOLVED)
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return
	}
	res.Closed, err = s.domainRepo.StatisticRepo.GetRentalComplaintStatistics(context.Background(), userId, database.RENTALCOMPLAINTSTATUSCLOSED)
	if errors.Is(err, database.ErrRecordNotFound) {
		return res, nil
	}
	return
}

func (s *service) GetTenantExpenditureStatistic(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) ([]dto.TenantExpenditureStatisticItem, error) {
	user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	if query.StartTime.Before(user.CreatedAt) {
		query.StartTime = user.CreatedAt
	}

	// for example: endDate = May 18th 2024, startDate = March 10th 2024, then derived intervals are [March 10th - April 1st], [April 2nd - May 1st], [May 2nd - May 18th]

	// Start from the first day of the start month
	current := query.StartTime
	res := make([]dto.TenantExpenditureStatisticItem, 0)
	for current.Before(query.EndTime) {
		// Calculate the end of the current month
		endOfMonth := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, query.StartTime.Location()).AddDate(0, 1, -1)

		// Choose the end of the current month or the end date, whichever is earlier
		intervalEnd := endOfMonth
		if endOfMonth.After(query.EndTime) {
			intervalEnd = query.EndTime
		}

		expenditure, err := s.domainRepo.StatisticRepo.GetTenantExpenditure(context.Background(), userId, dto.RentalPaymentStatisticQuery{
			StartTime: current,
			EndTime:   intervalEnd,
		})
		if err != nil {
			return nil, err
		}

		res = append(res, dto.TenantExpenditureStatisticItem{
			StartTime:   current,
			EndTime:     intervalEnd,
			Expenditure: expenditure,
		})

		// Move to the next month
		current = endOfMonth.AddDate(0, 0, 1)
	}

	return res, nil
}

func (s *service) GetTenantArrearsStatistic(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res dto.TenantArrearsStatistic, err error) {
	res.Payments, err = s.domainRepo.StatisticRepo.GetTenantPendingPayments(context.Background(), userId, *query)
	if err != nil {
		return
	}

	res.Total, err = s.domainRepo.StatisticRepo.GetTotalTenantPendingPayments(context.Background(), userId)
	return
}
