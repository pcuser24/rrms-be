package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/application/repo"
	"github.com/user2410/rrms-backend/internal/domain/application/utils"
	chat_dto "github.com/user2410/rrms-backend/internal/domain/chat/dto"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	"github.com/user2410/rrms-backend/internal/domain/notification"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/pkg/ds/set"

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
	CheckApplicationUpdatability(aid int64, uid uuid.UUID) (bool, error)
	CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error)
	GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error)
	CreateReminder(aid int64, userId uuid.UUID, data *dto.CreateReminder) (*model.ReminderModel, error)
	GetRemindersOfUser(userId uuid.UUID, aid int64) ([]model.ReminderModel, error)
	GetRentalByApplicationId(aid int64) (*rental_model.RentalModel, error)
	UpdateReminderStatus(aid int64, userId uuid.UUID, data *dto.UpdateReminderStatus) error
}

type service struct {
	aRepo               repo.Repo
	cRepo               chat_repo.Repo
	lRepo               listing_repo.Repo
	pRepo               property_repo.Repo
	taskDistributor     asynctask.TaskDistributor
	notificationAdapter *notification.WSNotificationAdapter
}

func NewService(
	aRepo repo.Repo,
	cRepo chat_repo.Repo,
	lRepo listing_repo.Repo,
	pRepo property_repo.Repo,
	taskDistributor asynctask.TaskDistributor,
	notificationAdapter *notification.WSNotificationAdapter,
) Service {
	return &service{
		aRepo:               aRepo,
		cRepo:               cRepo,
		lRepo:               lRepo,
		pRepo:               pRepo,
		taskDistributor:     taskDistributor,
		notificationAdapter: notificationAdapter,
	}
}

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
	// Check if there is any application of this user within 30 days
	appIds, err := s.aRepo.GetApplicationsByUserId(context.Background(), data.CreatorID, time.Now().AddDate(0, 0, -30), 1, 0)
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

func (s *service) CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error) {
	a, err := s.aRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return nil, err
	}
	if a.CreatorID == uuid.Nil {
		return nil, ErrAnonymousApplicant
	}

	return s.cRepo.CreateMsgroup(context.Background(), userId, &chat_dto.CreateMsgGroup{
		Name: utils.GetResourceName(aid),
		Members: []chat_dto.CreateMsgGroupMember{
			{UserId: userId},
			{UserId: a.CreatorID},
		},
	})
}

func (s *service) GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error) {
	return s.cRepo.GetMsgGroupByName(context.Background(), userId, utils.GetResourceName(aid))
}

func (s *service) CheckApplicationVisibility(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckVisibility(context.Background(), aid, uid)
}

func (s *service) CheckApplicationUpdatability(aid int64, uid uuid.UUID) (bool, error) {
	return s.aRepo.CheckUpdatability(context.Background(), aid, uid)
}

func (s *service) CreateReminder(aid int64, userId uuid.UUID, data *dto.CreateReminder) (*model.ReminderModel, error) {
	members := set.NewSet[uuid.UUID]()

	// get application
	application, err := s.aRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return nil, err
	}
	if application.CreatorID != uuid.Nil { // not an anonymous applicant
		members.Add(application.CreatorID)
	}
	if application.CreatorID == userId { // current user is the applicant
		pManagers, err := s.pRepo.GetPropertyManagers(context.Background(), application.PropertyID)
		if err != nil {
			return nil, err
		}
		for _, m := range pManagers {
			members.Add(m.ManagerID)
		}
	} else { // current user is a manager
		members.Add(userId)
	}

	for m := range members {
		data.Members = append(data.Members, m)
	}

	res, err := s.aRepo.CreateReminder(context.Background(), aid, userId, data)
	if err != nil {
		return nil, err
	}

	n, err := json.Marshal(*res)
	if err != nil {
		return res, err
	}
	go s.notificationAdapter.PushMessage(notification.Notification{
		UserId:  userId,
		Payload: n,
	})

	return res, err
}

func (s *service) GetRemindersOfUser(userId uuid.UUID, aid int64) ([]model.ReminderModel, error) {
	return s.aRepo.GetRemindersOfUser(context.Background(), aid, userId)
}

var (
	ErrInvalidReminderStatusTransition = fmt.Errorf("invalid reminder status transition")
)

func (s *service) UpdateReminderStatus(aid int64, userId uuid.UUID, data *dto.UpdateReminderStatus) error {
	reminder, err := s.aRepo.GetReminderById(context.Background(), data.ID)
	if err != nil {
		return err
	}

	switch data.Status {
	case database.REMINDERSTATUSINPROGRESS:
		if reminder.Status != database.REMINDERSTATUSPENDING {
			return ErrInvalidReminderStatusTransition
		}
	case database.REMINDERSTATUSCOMPLETED:
	case database.REMINDERSTATUSCANCELLED:
		if reminder.Status != database.REMINDERSTATUSPENDING && reminder.Status != database.REMINDERSTATUSINPROGRESS {
			return ErrInvalidReminderStatusTransition
		}
	default:
		return ErrInvalidReminderStatusTransition
	}

	rowsAffected, err := s.aRepo.UpdateReminderStatus(context.Background(), aid, data.ID, userId, data.Status)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUnauthorizedUpdate
	}

	return nil
}

func (s *service) GetRentalByApplicationId(aid int64) (*rental_model.RentalModel, error) {
	return s.aRepo.GetRentalByApplicationId(context.Background(), aid)
}
