package application

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/application/repo"
	chat_dto "github.com/user2410/rrms-backend/internal/domain/chat/dto"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Service interface {
	CreateApplication(data *dto.CreateApplication) (*model.ApplicationModel, error)
	GetApplicationById(id int64) (*model.ApplicationModel, error)
	GetApplicationByIds(ids []int64, fields []string, userId uuid.UUID) ([]model.ApplicationModel, error)
	GetApplicationsByUserId(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error)
	GetApplicationsToUser(uid uuid.UUID, q *dto.GetApplicationsToMeQuery) ([]model.ApplicationModel, error)
	UpdateApplicationStatus(aid int64, userId uuid.UUID, data *dto.UpdateApplicationStatus) error
	CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error)
	CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error)
	GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error)
}

type service struct {
	aRepo           repo.Repo
	cRepo           chat_repo.Repo
	lRepo           listing_repo.Repo
	pRepo           property_repo.Repo
	taskDistributor asynctask.TaskDistributor
}

func NewService(
	aRepo repo.Repo,
	cRepo chat_repo.Repo,
	lRepo listing_repo.Repo,
	pRepo property_repo.Repo,
	taskDistributor asynctask.TaskDistributor,
) Service {
	return &service{
		aRepo:           aRepo,
		cRepo:           cRepo,
		lRepo:           lRepo,
		pRepo:           pRepo,
		taskDistributor: taskDistributor,
	}
}

var (
	ErrListingIsClosed  = fmt.Errorf("listing is not active")
	ErrInvalidApplicant = fmt.Errorf("invalid applicant")
	ErrAlreadyApplied   = fmt.Errorf("user has already applied for this listing")
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
	// Check if there is an application of this user to this property within 30 days
	appIds, err := s.aRepo.GetApplicationsByUserId(
		context.Background(),
		data.CreatorID,
		time.Now().AddDate(0, 0, -30),
		1,
		0,
	)
	if err != nil {
		return nil, err
	}
	if len(appIds) > 0 {
		return nil, ErrAlreadyApplied
	}

	a, err := s.aRepo.CreateApplication(context.Background(), data)
	if err != nil {
		return nil, err
	}
	if err = s.taskDistributor.SendEmailOnNewApplication(
		context.Background(),
		&dto.SendEmailOnNewApplicationPayload{
			Email:         a.Email,
			Username:      a.FullName,
			ApplicationId: a.ID,
			ListingId:     a.ListingID,
		},
	); err != nil {
		log.Errorf("failed to distribute DistributeTaskSendNewApplicationEmail task: %v", err)
	}
	return a, nil
}

func (s *service) GetApplicationById(id int64) (*model.ApplicationModel, error) {
	return s.aRepo.GetApplicationById(context.Background(), id)
}

func (s *service) GetApplicationByIds(ids []int64, fields []string, userId uuid.UUID) ([]model.ApplicationModel, error) {
	var _ids []int64
	for _, id := range ids {
		isVisible, err := s.CheckVisibility(id, userId)
		if err != nil {
			return nil, err
		}
		if isVisible {
			_ids = append(_ids, id)
		}
	}
	return s.aRepo.GetApplicationsByIds(context.Background(), _ids, fields)
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

	// send email to the applicant
	return s.taskDistributor.UpdateApplicationStatus(context.Background(), &dto.UpdateApplicationStatusPayload{
		Email:         a.Email,
		ApplicationId: aid,
		OldStatus:     a.Status,
		NewStatus:     data.Status,
		Message:       data.Message,
	})
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

func (s *service) CheckVisibility(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckVisibility(context.Background(), aid, uid)
}

var (
	ErrAnonymousApplicant = errors.New("anonymous applicant")
)

func getGroupName(aid int64) string {
	return fmt.Sprintf("[APPLICATION_%d]", aid)
}

func (s *service) CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error) {
	a, err := s.aRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return nil, err
	}
	if a.CreatorID == uuid.Nil {
		return nil, ErrAnonymousApplicant
	}

	return s.cRepo.CreateMsgroup(context.Background(), userId, &chat_dto.CreateMsgGroup{
		Name: getGroupName(aid),
		Members: []chat_dto.CreateMsgGroupMember{
			{UserId: userId},
			{UserId: a.CreatorID},
		},
	})
}

func (s *service) GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error) {
	return s.cRepo.GetMsgGroupByName(context.Background(), userId, getGroupName(aid))
}

func (s *service) CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckVisibility(context.Background(), aid, uid)
}
