package service

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) PreCreateRental(data *dto.PreCreateRental, creatorID uuid.UUID) error {
	ext := filepath.Ext(data.Avatar.Name)
	fname := data.Avatar.Name[:len(data.Avatar.Name)-len(ext)]
	// key = creatorID + "/" + "/property" + filename
	objKey := fmt.Sprintf("%s/rentals/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

	url, err := s.s3Client.GetPutObjectPresignedURL(
		s.imageBucketName, objKey, data.Avatar.Type, data.Avatar.Size, UPLOAD_URL_LIFETIME*time.Minute,
	)
	if err != nil {
		return err
	}
	data.Avatar.Url = url.URL
	return nil
}

func (s *service) CreateRental(data *dto.CreateRental, userId uuid.UUID) (model.RentalModel, error) {
	expiryDate := data.MoveinDate.AddDate(0, int(data.RentalPeriod), 0)
	today := time.Now().Truncate(24 * time.Hour) // time representing today at 00:00:00
	if expiryDate.Before(today) {
		return model.RentalModel{}, ErrInvalidRentalExpired
	}

	// TODO: validate applicationId, propertyId, unitId
	data.CreatorID = userId
	rental, err := s.domainRepo.RentalRepo.CreateRental(context.Background(), data)
	if err != nil {
		return model.RentalModel{}, err
	}

	// plan rental payments
	_, err = s.domainRepo.RentalRepo.PlanRentalPayment(context.Background(), rental.ID)
	if err != nil {
		// TODO: log the error
	}
	// TODO: send notification to user
	return rental, nil
}

func (s *service) GetRental(id int64) (model.RentalModel, error) {
	return s.domainRepo.RentalRepo.GetRental(context.Background(), id)
}

func (s *service) UpdateRental(data *dto.UpdateRental, id int64) error {
	return s.domainRepo.RentalRepo.UpdateRental(context.Background(), data, id)
}

func (s *service) CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.domainRepo.RentalRepo.CheckRentalVisibility(context.Background(), id, userId)
}

func (s *service) GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error) {
	if query.Limit == nil {
		query.Limit = types.Ptr[int32](math.MaxInt32)
	}
	if query.Offset == nil {
		query.Offset = types.Ptr[int32](0)
	}
	rs, err := s.domainRepo.RentalRepo.GetManagedRentals(context.Background(), userId, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.RentalRepo.GetRentalsByIds(context.Background(), rs, query.Fields)
}

func (s *service) GetMyRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]model.RentalModel, error) {
	if query.Limit == nil {
		query.Limit = types.Ptr[int32](math.MaxInt32)
	}
	if query.Offset == nil {
		query.Offset = types.Ptr[int32](0)
	}
	rs, err := s.domainRepo.RentalRepo.GetMyRentals(context.Background(), userId, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.RentalRepo.GetRentalsByIds(context.Background(), rs, query.Fields)
}
