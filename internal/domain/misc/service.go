package misc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/user2410/rrms-backend/internal/domain/misc/dto"
	"github.com/user2410/rrms-backend/internal/domain/misc/model"
)

type Service interface {
	CreateNotificationDevice(userId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error)
	GetNotificationDevice(userId, sessionId uuid.UUID, token, platform string) (*model.NotificationDevice, error)

	GetNotificationsOfUser(userId uuid.UUID, limit, offset int32) ([]model.Notification, error)
}

type service struct {
	repo Repo

	cronEntries []cron.EntryID
}

func NewService(repo Repo, c *cron.Cron) Service {
	s := &service{
		repo: repo,

		cronEntries: []cron.EntryID{},
	}

	s.setupCronjob(c)

	return s
}

func (s *service) CreateNotificationDevice(userId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error) {
	sessionId := uuid.New()

	return s.repo.CreateNotificationDevice(context.Background(), userId, sessionId, payload)
}

func (s *service) GetNotificationDevice(userId, sessionId uuid.UUID, token, platform string) (*model.NotificationDevice, error) {
	res, err := s.repo.GetNotificationDevice(context.Background(), userId, sessionId, token, platform)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpdateNotificationDeviceTokenTimestamp(context.Background(), userId, sessionId)

	return &res, err
}

func (s *service) CreateNotification(userId uuid.UUID, payload *dto.CreateNotification) (model.Notification, error) {
	return s.repo.CreateNotification(context.Background(), payload)
}

func (s *service) setupCronjob(c *cron.Cron) ([]cron.EntryID, error) {
	var (
		entryID cron.EntryID
		err     error
	)
	entryID, err = c.AddFunc("@hourly", func() {
		err := s.repo.DeleteExpiredTokens(context.Background(), 60)
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
	return s.repo.GetNotificationsOfUser(context.Background(), userId, limit, offset)
}
