package misc

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

	SendNotification(payload *dto.CreateNotification, ntype NOTIFICATIONTYPE) error
	GetNotificationsOfUser(userId uuid.UUID, limit, offset int32) ([]model.Notification, error)
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

func (s *service) GetNotificationsOfUser(userId uuid.UUID, limit, offset int32) ([]model.Notification, error) {
	return s.domainRepo.MiscRepo.GetNotificationsOfUser(context.Background(), userId, limit, offset)
}

type NOTIFICATIONTYPE = string

const (
	NOTIFICATIONTYPE_CREATEAPPLICATION NOTIFICATIONTYPE = "CREATE_APPLICATION"
	NOTIFICATIONTYPE_UPDATEAPPLICATION NOTIFICATIONTYPE = "UPDATE_APPLICATION"
)

func (s *service) SendNotification(payload *dto.CreateNotification, ntype NOTIFICATIONTYPE) error {
	// TODO: get user notification preferences and filter out the ones that are not allowed
	// save to database
	nms, err := s.domainRepo.MiscRepo.CreateNotification(context.Background(), payload)
	if err != nil {
		return err
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
		"notificationType": ntype,
		"notifications":    items,
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
