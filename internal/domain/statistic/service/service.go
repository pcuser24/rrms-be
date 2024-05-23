package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	statistic_repo "github.com/user2410/rrms-backend/internal/domain/statistic/repo"
)

type Service interface {
	GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error)
	GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error)
	GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error)
	GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error)
	GetRentalPaymentIncomes(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentIncomeItem, err error)
}

type service struct {
	authRepo      auth_repo.Repo
	statisticRepo statistic_repo.Repo
	propertyRepo  property_repo.Repo
}

func NewService(authRepo auth_repo.Repo, statisticRepo statistic_repo.Repo, propertyRepo property_repo.Repo) Service {
	return &service{
		authRepo:      authRepo,
		statisticRepo: statisticRepo,
		propertyRepo:  propertyRepo,
	}
}

func (s *service) GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error) {
	managedProperties, err := s.propertyRepo.GetManagedProperties(context.Background(), userId, &property_dto.GetPropertiesQuery{})
	if err != nil {
		return
	}
	for _, p := range managedProperties {
		res.Properties = append(res.Properties, p.PropertyID)
	}

	res.OwnedProperties, err = s.statisticRepo.GetManagedPropertiesByRole(context.Background(), userId, "OWNER")
	if err != nil {
		return
	}

	res.OccupiedProperties, err = s.statisticRepo.GetOccupiedProperties(context.Background(), userId)
	if err != nil {
		return
	}

	res.Units, err = s.statisticRepo.GetManagedUnits(context.Background(), userId)
	if err != nil {
		return
	}

	res.OccupiedUnits, err = s.statisticRepo.GetOccupiedUnits(context.Background(), userId)
	if err != nil {
		return
	}

	res.PropertiesWithActiveListing, err = s.statisticRepo.GetPropertiesWithActiveListing(context.Background(), userId)
	if err != nil {
		return
	}

	res.MostRentedProperties, err = s.statisticRepo.GetMostRentedProperties(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.LeastRentedProperties, err = s.statisticRepo.GetLeastRentedProperties(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.MostRentedUnits, err = s.statisticRepo.GetMostRentedUnits(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	res.LeastRentedUnits, err = s.statisticRepo.GetLeastRentedUnits(context.Background(), userId, query.Limit, query.Offset)
	if err != nil {
		return
	}

	return
}

func (s *service) GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error) {
	now := time.Now()

	res.NewApplicationsThisMonth, err = s.statisticRepo.GetNewApplications(context.Background(), userId, now)
	if err != nil {
		return
	}

	res.NewApplicationsLastMonth, err = s.statisticRepo.GetNewApplications(context.Background(), userId, now.AddDate(0, -1, 0))
	return
}

func (s *service) GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error) {
	now := time.Now()
	res.NewMaintenancesThisMonth, err = s.statisticRepo.GetMaintenanceRequests(context.Background(), userId, now)
	if err != nil {
		return
	}

	res.NewMaintenancesLastMonth, err = s.statisticRepo.GetMaintenanceRequests(context.Background(), userId, now.AddDate(0, -1, 0))
	return
}

func (s *service) GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error) {
	user, err := s.authRepo.GetUserById(context.Background(), userId)
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

		payments, err := s.statisticRepo.GetRentalPaymentArrears(context.Background(), userId, dto.RentalPaymentStatisticQuery{
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
	user, err := s.authRepo.GetUserById(context.Background(), userId)
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

		income, err := s.statisticRepo.GetRentalPaymentIncomes(context.Background(), userId, dto.RentalPaymentStatisticQuery{
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
