package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (s *service) CreateRental(data *dto.CreateRental, userId uuid.UUID) (model.RentalModel, error) {
	expiryDate := data.MoveinDate.AddDate(0, int(data.RentalPeriod), 0)
	today := time.Now().Truncate(24 * time.Hour) // time representing today at 00:00:00
	if expiryDate.Before(today) {
		return model.RentalModel{}, ErrInvalidRentalExpired
	}

	// TODO: validate applicationId, propertyId, unitId
	data.CreatorID = userId
	rental, err := s.rRepo.CreateRental(context.Background(), data)
	if err != nil {
		return model.RentalModel{}, err
	}

	// plan rental payments
	_, err = s.rRepo.PlanRentalPayment(context.Background(), rental.ID)
	if err != nil {
		// TODO: log the error
	}
	// TODO: send notification to user
	return rental, nil
}

func (s *service) GetRental(id int64) (model.RentalModel, error) {
	return s.rRepo.GetRental(context.Background(), id)
}

func (s *service) UpdateRental(data *dto.UpdateRental, id int64) error {
	return s.rRepo.UpdateRental(context.Background(), data, id)
}

func (s *service) CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.rRepo.CheckRentalVisibility(context.Background(), id, userId)
}

func (s *service) GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error) {
	rs, err := s.rRepo.GetManagedRentals(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	return s.rRepo.GetRentalsByIds(context.Background(), rs, query.Fields)
}
