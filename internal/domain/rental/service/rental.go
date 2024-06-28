package service

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/pkg/ds/set"
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

func (s *service) CreatePreRental(data *dto.CreateRental, userId uuid.UUID) (rental_model.PreRental, error) {
	expiryDate := data.MoveinDate.AddDate(0, int(data.RentalPeriod), 0)
	today := time.Now().Truncate(24 * time.Hour) // time representing today at 00:00:00
	if expiryDate.Before(today) {
		return rental_model.RentalModel{}, ErrInvalidRentalExpired
	}

	// TODO: validate applicationId, propertyId, unitId
	data.CreatorID = userId
	rental, err := s.domainRepo.RentalRepo.CreatePreRental(context.Background(), data)
	if err != nil {
		return rental_model.RentalModel{}, err
	}

	err = s.notifyCreatePreRental(&rental, s.secret)
	if err != nil {
		// TODO: log the error

	}

	return rental, nil
}

// func (s *service) CreateRental(data *dto.CreateRental, userId uuid.UUID) (rental_model.RentalModel, error) {
// 	expiryDate := data.MoveinDate.AddDate(0, int(data.RentalPeriod), 0)
// 	today := time.Now().Truncate(24 * time.Hour) // time representing today at 00:00:00
// 	if expiryDate.Before(today) {
// 		return rental_model.RentalModel{}, ErrInvalidRentalExpired
// 	}

// 	// TODO: validate applicationId, propertyId, unitId
// 	data.CreatorID = userId
// 	rental, err := s.domainRepo.RentalRepo.CreateRental(context.Background(), data)
// 	if err != nil {
// 		return rental_model.RentalModel{}, err
// 	}

// 	// plan rental payments
// 	_, err = s.domainRepo.RentalRepo.PlanRentalPayment(context.Background(), rental.ID)
// 	if err != nil {
// 		// TODO: log the error
// 	}

// 	err = s.notifyCreatePreRental(&rental)
// 	if err != nil {
// 		// TODO: log the error
// 	}

// 	return rental, nil
// }

func (s *service) GetRental(id int64) (rental_model.RentalModel, error) {
	return s.domainRepo.RentalRepo.GetRental(context.Background(), id)
}

func (s *service) UpdateRental(data *dto.UpdateRental, id int64) error {
	return s.domainRepo.RentalRepo.UpdateRental(context.Background(), data, id)
}

func (s *service) FilterVisibleRentals(userId uuid.UUID, ids []int64) ([]int64, error) {
	return s.domainRepo.RentalRepo.FilterVisibleRentals(context.Background(), userId, ids)
}

func (s *service) CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.domainRepo.RentalRepo.CheckRentalVisibility(context.Background(), id, userId)
}

func (s *service) CheckPreRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.domainRepo.RentalRepo.CheckPreRentalVisibility(context.Background(), id, userId)
}

func (s *service) GetManagedRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]rental_model.RentalModel, error) {
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

func (s *service) GetMyRentals(userId uuid.UUID, query *dto.GetRentalsQuery) ([]rental_model.RentalModel, error) {
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

func (s *service) GetPreRentalExtended(id int64, userId uuid.UUID, key string) (dto.GetPreRentalResponse, error) {
	// TODO: access control
	pr, err := s.domainRepo.RentalRepo.GetPreRental(context.Background(), id)
	if err != nil {
		return dto.GetPreRentalResponse{}, err
	}
	property, err := s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), pr.PropertyID)
	if err != nil {
		return dto.GetPreRentalResponse{}, err
	}
	unit, err := s.domainRepo.UnitRepo.GetUnitById(context.Background(), pr.UnitID)
	if err != nil {
		return dto.GetPreRentalResponse{}, err
	}

	return dto.GetPreRentalResponse{
		PreRental: &pr,
		Property:  property,
		Unit:      unit,
	}, nil
}

func (s *service) GetPreRentalByID(id int64) (rental_model.PreRental, error) {
	return s.domainRepo.RentalRepo.GetPreRental(context.Background(), id)
}

func (s *service) UpdatePreRentalState(id int64, payload *dto.UpdatePreRental) (int64, error) {
	var (
		rental rental_model.PreRental
		err    error
	)

	preRental, err := s.domainRepo.RentalRepo.GetPreRental(context.Background(), id)
	if err != nil {
		return 0, err
	}

	if payload.State == "APPROVED" {
		rental, err = s.domainRepo.RentalRepo.MovePreRentalToRental(context.Background(), id)
		if err != nil {
			return 0, err
		}
		// plan rental payments
		_, err = s.domainRepo.RentalRepo.PlanRentalPayment(context.Background(), rental.ID)
		if err != nil {
			// TODO: log the error
		}
		// TODO: send notification about the new rental payment plan
		// s.notifyCreateRentalPayment(&rental, )
		// send notification to managers about the acceptance
		err = s.notifyUpdatePreRental(&preRental, &rental, payload)
	} else if payload.State == "REVIEW" {
		// send notifcation to managers to request a review
		err = s.notifyUpdatePreRental(&preRental, nil, payload)
	} else if payload.State == "REJECTED" {
		err := s.domainRepo.RentalRepo.RemovePreRental(context.Background(), id)
		if err != nil {
			return 0, err
		}
		// send notification to managers about the rejection
		err = s.notifyUpdatePreRental(&preRental, nil, payload)
	}

	if err != nil {
		// TODO: log the error
	}

	return rental.ID, nil
}

func (s *service) GetPreRentalsToMe(userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]rental_model.PreRental, error) {
	// TODO: access control
	return s.domainRepo.RentalRepo.GetPreRentalsToTenant(context.Background(), userId, query)
}

func (s *service) GetManagedPreRentals(userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]rental_model.PreRental, error) {
	// TODO: access control
	return s.domainRepo.RentalRepo.GetManagedPreRentals(context.Background(), userId, query)
}

func (s *service) GetRentalByIds(userId uuid.UUID, ids []int64, fields []string) ([]rental_model.RentalModel, error) {
	idSet := set.NewSet[int64]().AddAll(ids...)
	visibleIds, err := s.domainRepo.RentalRepo.FilterVisibleRentals(context.Background(), userId, idSet.ToSlice())
	if err != nil {
		return nil, err
	}
	return s.domainRepo.RentalRepo.GetRentalsByIds(context.Background(), visibleIds, fields)
}
