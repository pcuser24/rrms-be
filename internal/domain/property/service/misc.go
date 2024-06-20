package service

import (
	"context"
	"errors"

	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	"github.com/google/uuid"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
)

func (s *service) GetUnitsOfProperty(id uuid.UUID) ([]unit_model.UnitModel, error) {
	return s.domainRepo.UnitRepo.GetUnitsOfProperty(context.Background(), id)
}

func (s *service) GetListingsOfProperty(id uuid.UUID, query *listing_dto.GetListingsOfPropertyQuery) ([]listing_model.ListingModel, error) {
	ids, err := s.domainRepo.PropertyRepo.GetListingsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.ListingRepo.GetListingsByIds(context.Background(), ids, query.Fields)
}

func (s *service) GetApplicationsOfProperty(id uuid.UUID, query *application_dto.GetApplicationsOfPropertyQuery) ([]application_model.ApplicationModel, error) {
	ids, err := s.domainRepo.PropertyRepo.GetApplicationsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.ApplicationRepo.GetApplicationsByIds(context.Background(), ids, query.Fields)
}

func (s *service) GetRentalsOfProperty(id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]rental_model.RentalModel, error) {
	ids, err := s.domainRepo.PropertyRepo.GetRentalsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.RentalRepo.GetRentalsByIds(context.Background(), ids, query.Fields)
}

var ErrUserIsAlreadyManager = errors.New("user is already a manager of the property")

func (s *service) CreatePropertyManagerRequest(data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error) {
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), data.PropertyID)
	if err != nil {
		return property_model.NewPropertyManagerRequest{}, err
	}
	if exists, err := func() (bool, error) {
		for _, manager := range managers {
			user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), manager.ManagerID)
			if err != nil {
				return false, err
			}
			if user.Email == data.Email {
				return true, nil
			}
		}
		return false, nil
	}(); exists || err != nil {
		if err != nil {
			return property_model.NewPropertyManagerRequest{}, err
		}
		return property_model.NewPropertyManagerRequest{}, ErrUserIsAlreadyManager
	}

	user, err := s.domainRepo.AuthRepo.GetUserByEmail(context.Background(), data.Email)
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return property_model.NewPropertyManagerRequest{}, err
	} else {
		data.UserID = user.ID
	}

	return s.domainRepo.PropertyRepo.CreatePropertyManagerRequest(context.Background(), data)

	// TODO: send email to user, push notification if user is already registered
}

func (s *service) GetNewPropertyManagerRequestsToUser(uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error) {
	return s.domainRepo.PropertyRepo.GetNewPropertyManagerRequestsToUser(context.Background(), uid, limit, offset)
}

var ErrUpdateRequestInfoMismatch = errors.New("request update info mismatch")

func (s *service) UpdatePropertyManagerRequest(pid, uid uuid.UUID, requestId int64, approved bool) error {
	user, err := s.domainRepo.AuthRepo.GetUserById(context.Background(), uid)
	if err != nil {
		return err
	}
	request, err := s.domainRepo.PropertyRepo.GetNewPropertyManagerRequest(context.Background(), requestId)
	if err != nil {
		return err
	}
	if (request.UserID != uuid.Nil && uid != request.UserID) ||
		request.PropertyID != pid ||
		user.Email != request.Email {
		return ErrUpdateRequestInfoMismatch
	}

	return s.domainRepo.PropertyRepo.UpdatePropertyManagerRequest(context.Background(), requestId, user.ID, approved)
}
