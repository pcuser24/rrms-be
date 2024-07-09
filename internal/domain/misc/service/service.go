package service

import (
	"context"
	"errors"
	"log"
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	"github.com/user2410/rrms-backend/internal/domain/misc/dto"
	"github.com/user2410/rrms-backend/internal/domain/misc/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/notification"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

type Service interface {
	CreateNotificationDevice(userId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error)
	GetNotificationDevice(userId, sessionId uuid.UUID, token, platform string) ([]model.NotificationDevice, error)

	SendNotification(payload *dto.CreateNotification) error
	GetNotificationsOfUser(userId uuid.UUID, query dto.GetNotificationsOfUserQuery) ([]model.Notification, error)
	UpdateNotification(data *dto.UpdateNotification) error

	GetNotificationManagersTargets(propertyID uuid.UUID) ([]dto.CreateNotificationTarget, error)
	GetNotificationTenantTargets(tenantID uuid.UUID, tenantEmail string) (dto.CreateNotificationTarget, error)
}

type service struct {
	domainRepo           repos.DomainRepo
	notificationEndpoint notification.NotificationEndpoint

	cronEntries []cron.EntryID
}

func NewService(domainRepo repos.DomainRepo, notificationEndpoint notification.NotificationEndpoint, c *cron.Cron) Service {
	s := &service{
		domainRepo:           domainRepo,
		notificationEndpoint: notificationEndpoint,
		cronEntries:          []cron.EntryID{},
	}

	s.setupCronjob(c)

	return s
}

func (s *service) CreateNotificationDevice(userId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error) {
	nds, err := s.GetNotificationDevice(userId, uuid.Nil, payload.Token, string(payload.Platform))
	if err != nil {
		return model.NotificationDevice{}, err
	}
	if len(nds) > 0 {
		return nds[0], nil
	}

	sessionId := uuid.New()
	return s.domainRepo.MiscRepo.CreateNotificationDevice(context.Background(), userId, sessionId, payload)
}

func (s *service) GetNotificationDevice(userId, sessionId uuid.UUID, token, platform string) ([]model.NotificationDevice, error) {
	res, err := s.domainRepo.MiscRepo.GetNotificationDevice(context.Background(), userId, sessionId, token, platform)
	if err != nil {
		return nil, err
	}
	err = s.domainRepo.MiscRepo.UpdateNotificationDeviceTokenTimestamp(context.Background(), userId, sessionId)

	return res, err
}

func (s *service) setupCronjob(c *cron.Cron) ([]cron.EntryID, error) {
	var (
		entryID cron.EntryID
		err     error
	)
	entryID, err = c.AddFunc("@hourly", func() {
		err := s.domainRepo.MiscRepo.DeleteExpiredTokens(context.Background(), 60)
		if err != nil {
			log.Println("failed to delete expired tokens:", err)
		}
	})
	if err != nil {
		return nil, err
	}
	s.cronEntries = append(s.cronEntries, entryID)

	return s.cronEntries, nil
}

func (s *service) GetNotificationsOfUser(userId uuid.UUID, query dto.GetNotificationsOfUserQuery) ([]model.Notification, error) {
	return s.domainRepo.MiscRepo.GetNotificationsOfUser(context.Background(), userId, query)
}

type NOTIFICATIONTYPE = string

const (
	NOTIFICATIONTYPE_CREATEAPPLICATION NOTIFICATIONTYPE = "CREATE_APPLICATION"
	NOTIFICATIONTYPE_UPDATEAPPLICATION NOTIFICATIONTYPE = "UPDATE_APPLICATION"

	NOTIFICATIONTYPE_CREATEPRERENTAL NOTIFICATIONTYPE = "CREATE_PRERENTAL"
	NOTIFICATIONTYPE_UPDATEPRERENTAL NOTIFICATIONTYPE = "UPDATE_PRERENTAL"

	NOTIFICATIONTYPE_CREATERENTALPAYMENT NOTIFICATIONTYPE = "CREATE_RENTALPAYMENT"
	NOTIFICATIONTYPE_UPDATERENTALPAYMENT NOTIFICATIONTYPE = "UPDATE_RENTALPAYMENT"

	NOTIFICATIONTYPE_CREATECONTRACT NOTIFICATIONTYPE = "CREATE_CONTRACT"
	NOTIFICATIONTYPE_UPDATECONTRACT NOTIFICATIONTYPE = "UPDATE_CONTRACT"

	NOTIFICATIONTYPE_CREATERENTALCOMPLAINT       NOTIFICATIONTYPE = "CREATE_RENTALCOMPLAINT"
	NOTIFICATIONTYPE_UPDATERENTALCOMPLAINTSTATUS NOTIFICATIONTYPE = "UPDATE_RENTALCOMPLAINTSTATUS"
	NOTIFICATIONTYPE_CREATERENTALCOMPLAINTREPLY  NOTIFICATIONTYPE = "CREATE_RENTALCOMPLAINTREPLY"

	NOTIFICATIONTYPE_CREATEPROPERTYVERIFICATIONSTATUS NOTIFICATIONTYPE = "CREATE_PROPERTYVERIFICATIONSTATUS"
	NOTIFICATIONTYPE_UPDATEPROPERTYVERIFICATIONSTATUS NOTIFICATIONTYPE = "UPDATE_PROPERTYVERIFICATIONSTATUS"
)

func (s *service) SendNotification(payload *dto.CreateNotification) error {
	// TODO: get user notification preferences and filter out the ones that are not allowed
	// save to database
	nms, err := s.domainRepo.MiscRepo.CreateNotification(context.Background(), payload)
	if err != nil {
		return err
	}
	if nms == nil {
		return nil
	}

	type metadata struct {
		ID        int64                        `json:"id"`
		UserID    uuid.UUID                    `json:"userId"`
		Seen      bool                         `json:"seen"`
		Target    string                       `json:"target"`
		Channel   database.NOTIFICATIONCHANNEL `json:"channel"`
		CreatedAt time.Time                    `json:"createdAt"`
		UpdatedAt time.Time                    `json:"updatedAt"`
	}
	var items []metadata
	for _, nm := range nms {
		items = append(items, metadata{
			ID:        nm.ID,
			UserID:    nm.UserID,
			Target:    nm.Target,
			Channel:   nm.Channel,
			Seen:      nm.Seen,
			CreatedAt: nm.CreatedAt,
			UpdatedAt: nm.UpdatedAt,
		})
	}
	// send to notification endpoint
	nt := notification.NotificationTransport{
		Title:   payload.Title,
		Content: payload.Content,
		Data:    maps.Clone(payload.Data),
	}
	maps.Copy(nt.Data, map[string]interface{}{
		"notifications": items,
	})

	emails := set.NewSet[string]()
	tokens := set.NewSet[string]()
	for _, t := range payload.Targets {
		emails.AddAll(t.Emails...)
		tokens.AddAll(t.Tokens...)
	}
	if len(emails) > 0 {
		nt.EmailChannel = &notification.NotificationEmailChannel{
			To:      emails.ToSlice(),
			CC:      []string{},
			BCC:     []string{},
			ReplyTo: []string{},
		}
	}
	if len(tokens) > 0 {
		nt.PushChannel = &notification.NotificationPushChannel{
			Tokens: tokens.ToSlice(),
		}
	}

	errs := s.notificationEndpoint.SendNotification(context.Background(), &nt)
	return errors.Join(errs...)
}

func (s *service) UpdateNotification(data *dto.UpdateNotification) error {
	return s.domainRepo.MiscRepo.UpdateNotification(context.Background(), data)
}

func (s *service) GetNotificationManagersTargets(propertyID uuid.UUID) ([]dto.CreateNotificationTarget, error) {
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), propertyID)
	if err != nil {
		return nil, err
	}
	managerIds := make([]uuid.UUID, 0)
	for _, m := range managers {
		managerIds = append(managerIds, m.ManagerID)
	}

	us, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), managerIds, []string{"email"})
	if err != nil {
		return nil, err
	}

	targets := make([]dto.CreateNotificationTarget, 0)
	for _, u := range us {
		target := dto.CreateNotificationTarget{
			UserId: u.ID,
			Emails: []string{u.Email},
		}
		ds, err := s.GetNotificationDevice(u.ID, uuid.Nil, "", "")
		if err != nil {
			return nil, err
		}
		for _, d := range ds {
			target.Tokens = append(target.Tokens, d.Token)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (s *service) GetNotificationTenantTargets(tenantID uuid.UUID, tenantEmail string) (dto.CreateNotificationTarget, error) {
	emailTargets := set.NewSet[string]()
	pushTargets := set.NewSet[string]()
	emailTargets.Add(tenantEmail)
	if tenantID != uuid.Nil {
		ts, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), []uuid.UUID{tenantID}, []string{"email"})
		if err != nil {
			return dto.CreateNotificationTarget{}, err
		}
		if len(ts) > 0 {
			emailTargets.Add(ts[0].Email)
			ds, err := s.GetNotificationDevice(tenantID, uuid.Nil, "", "")
			if err != nil {
				return dto.CreateNotificationTarget{}, err
			}
			for _, d := range ds {
				pushTargets.Add(d.Token)
			}
		}
	}
	return dto.CreateNotificationTarget{
		UserId: tenantID,
		Emails: emailTargets.ToSlice(),
		Tokens: pushTargets.ToSlice(),
	}, nil
}
