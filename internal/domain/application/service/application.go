package service

import (
	"context"
	"fmt"
	"slices"

	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/pkg/ds/set"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

var (
	ErrListingIsClosed  = fmt.Errorf("listing is not active")
	ErrInvalidApplicant = fmt.Errorf("invalid applicant")
	ErrAlreadyApplied   = fmt.Errorf("user has already applied to this property within 30 days")
)

func (s *service) CreateApplication(data *dto.CreateApplication) (*model.ApplicationModel, error) {
	// Check eligibility of the user to apply for this listing
	// Check if the listing is still open
	if data.ListingID != uuid.Nil {
		expired, err := s.lRepo.CheckListingExpired(context.Background(), data.ListingID)
		if err != nil {
			return nil, err
		}
		if expired {
			return nil, ErrListingIsClosed
		}
	}
	// Check if the current user is a manager of the property
	pManagers, err := s.pRepo.GetPropertyManagers(context.Background(), data.PropertyID)
	if err != nil {
		return nil, err
	}
	if slices.IndexFunc(pManagers, func(m property_model.PropertyManagerModel) bool { return m.ManagerID == data.CreatorID }) != -1 {
		return nil, ErrInvalidApplicant
	}

	return s.aRepo.CreateApplication(context.Background(), data)
}

func (s *service) GetApplicationById(id int64) (*model.ApplicationModel, error) {
	return s.aRepo.GetApplicationById(context.Background(), id)
}

func (s *service) GetApplicationByIds(ids []int64, fields []string, userId uuid.UUID) ([]model.ApplicationModel, error) {
	_ids := set.NewSet[int64]()
	for _, id := range ids {
		isVisible, err := s.CheckApplicationVisibility(id, userId)
		if err != nil {
			return nil, err
		}
		if isVisible {
			_ids.Add(id)
		}
	}
	return s.aRepo.GetApplicationsByIds(context.Background(), _ids.ToSlice(), fields)
}

var (
	ErrInvalidStatusTransition = fmt.Errorf("invalid status transition")
	ErrUnauthorizedUpdate      = fmt.Errorf("unauthorized update")
)

func (s *service) UpdateApplicationStatus(aid int64, userId uuid.UUID, data *dto.UpdateApplicationStatus) error {
	a, err := s.aRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return err
	}

	switch data.Status {
	case database.APPLICATIONSTATUSWITHDRAWN:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSCONDITIONALLYAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSREJECTED:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	}

	rowsAffected, err := s.aRepo.UpdateApplicationStatus(context.Background(), aid, userId, data.Status)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUnauthorizedUpdate
	}

	return nil
}

func (s *service) GetApplicationsByUserId(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error) {
	ids, err := s.aRepo.GetApplicationsByUserId(
		context.Background(),
		uid,
		q.CreatedBefore,
		q.Limit,
		q.Offset,
	)
	if err != nil {
		return nil, err
	}

	return s.aRepo.GetApplicationsByIds(
		context.Background(),
		ids,
		q.Fields,
	)
}

func (s *service) GetApplicationsToUser(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error) {
	ids, err := s.aRepo.GetApplicationsToUser(
		context.Background(),
		uid,
		q.CreatedBefore,
		q.Limit,
		q.Offset,
	)
	if err != nil {
		return nil, err
	}

	return s.aRepo.GetApplicationsByIds(
		context.Background(),
		ids,
		q.Fields,
	)
}

func (s *service) CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckVisibility(context.Background(), aid, uid)
}

func (s *service) CheckApplicationUpdatability(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckUpdatability(context.Background(), aid, uid)
}
