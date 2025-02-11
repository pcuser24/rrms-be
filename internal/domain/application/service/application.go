package service

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/pkg/ds/set"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

const (
	MAX_IMAGE_SIZE      = 10 * 1024 * 1024 // 10MB
	UPLOAD_URL_LIFETIME = 5                // 5 minutes
)

var (
	ErrListingIsClosed  = fmt.Errorf("listing is not active")
	ErrInvalidApplicant = fmt.Errorf("invalid applicant")
	ErrAlreadyApplied   = fmt.Errorf("user has already applied to this property within 30 days")

	basePath = "internal/domain/application/service/templates"
)

func (s *service) PreCreateApplication(data *dto.PreCreateApplication, creatorID uuid.UUID) error {
	ext := filepath.Ext(data.Avatar.Name)
	fname := data.Avatar.Name[:len(data.Avatar.Name)-len(ext)]
	// key = creatorID + "/" + "/property" + filename
	objKey := fmt.Sprintf("%s/applications/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

	url, err := s.s3Client.GetPutObjectPresignedURL(
		s.imageBucketName, objKey, data.Avatar.Type, data.Avatar.Size, UPLOAD_URL_LIFETIME*time.Minute,
	)
	if err != nil {
		return err
	}
	data.Avatar.Url = url.URL
	return nil
}

func (s *service) CreateApplication(data *dto.CreateApplication) (*model.ApplicationModel, error) {
	// Check eligibility of the user to apply for this listing
	// Check if the listing is still open
	if data.ListingID != uuid.Nil {
		expired, err := s.domainRepo.ListingRepo.CheckListingExpired(context.Background(), data.ListingID)
		if err != nil {
			return nil, err
		}
		if expired {
			return nil, ErrListingIsClosed
		}
	}
	// Check if the current user is a manager of the property
	pManagers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), data.PropertyID)
	if err != nil {
		return nil, err
	}
	if slices.IndexFunc(pManagers, func(m property_model.PropertyManagerModel) bool { return m.ManagerID == data.CreatorID }) != -1 {
		return nil, ErrInvalidApplicant
	}

	am, err := s.domainRepo.ApplicationRepo.CreateApplication(context.Background(), data)
	if err != nil {
		return nil, err
	}

	// Send notifications
	// err = s.SendNotificationOnNewApplication(am)
	err = s.asynctaskDistributor.DistributeTaskJSON(context.Background(), asynctask.APPLICATION_NEW, am)

	return am, err
}

func (s *service) GetApplicationById(id int64) (*model.ApplicationModel, error) {
	return s.domainRepo.ApplicationRepo.GetApplicationById(context.Background(), id)
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
	return s.domainRepo.ApplicationRepo.GetApplicationsByIds(context.Background(), _ids.ToSlice(), fields)
}

var (
	ErrInvalidStatusTransition = fmt.Errorf("invalid status transition")
	ErrUnauthorizedUpdate      = fmt.Errorf("unauthorized update")
)

func (s *service) UpdateApplicationStatus(aid int64, userId uuid.UUID, data *dto.UpdateApplicationStatus) error {
	a, err := s.domainRepo.ApplicationRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return err
	}

	willNotify := true

	switch data.Status {
	case database.APPLICATIONSTATUSWITHDRAWN:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSCONDITIONALLYAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING {
			return ErrInvalidStatusTransition
		}
		willNotify = false
	case database.APPLICATIONSTATUSAPPROVED:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	case database.APPLICATIONSTATUSREJECTED:
		if a.Status != database.APPLICATIONSTATUSPENDING && a.Status != database.APPLICATIONSTATUSCONDITIONALLYAPPROVED {
			return ErrInvalidStatusTransition
		}
	}

	rowsAffected, err := s.domainRepo.ApplicationRepo.UpdateApplicationStatus(context.Background(), aid, userId, data.Status)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUnauthorizedUpdate
	}

	// Send notifications
	if willNotify {
		data := dto.NotificationOnUpdateApplication{
			Application: a,
			Status:      data.Status,
		}
		err = s.asynctaskDistributor.DistributeTaskJSON(context.Background(), asynctask.APPLICATION_UPDATE, data)
	}

	return err
}

func (s *service) GetApplicationsByUserId(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error) {
	ids, err := s.domainRepo.ApplicationRepo.GetApplicationsByUserId(
		context.Background(),
		uid,
		q.Limit,
		q.Offset,
	)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.ApplicationRepo.GetApplicationsByIds(
		context.Background(),
		ids,
		q.Fields,
	)
}

func (s *service) GetApplicationsToUser(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error) {
	ids, err := s.domainRepo.ApplicationRepo.GetApplicationsToUser(
		context.Background(),
		uid,
		q.Limit,
		q.Offset,
	)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.ApplicationRepo.GetApplicationsByIds(
		context.Background(),
		ids,
		q.Fields,
	)
}

func (s *service) CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error) {
	return s.domainRepo.ApplicationRepo.CheckVisibility(context.Background(), aid, uid)
}

func (s *service) CheckApplicationUpdatability(aid int64, uid uuid.UUID) (bool, error) {
	return s.domainRepo.ApplicationRepo.CheckUpdatability(context.Background(), aid, uid)
}
