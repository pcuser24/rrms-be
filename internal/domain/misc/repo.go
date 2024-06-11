package misc

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/misc/dto"
	"github.com/user2410/rrms-backend/internal/domain/misc/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateNotificationDevice(ctx context.Context, userId, sessionId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error)
	GetNotificationDevice(ctx context.Context, userId, sessionId uuid.UUID, token, platform string) (model.NotificationDevice, error)
	UpdateNotificationDeviceTokenTimestamp(ctx context.Context, userId, sessionId uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context, interval int32) error

	CreateNotification(ctx context.Context, data *dto.CreateNotification) (model.Notification, error)
	GetNotificationsOfUser(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]model.Notification, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateNotificationDevice(ctx context.Context, userId, sessionId uuid.UUID, payload *dto.CreateNotificationDevice) (model.NotificationDevice, error) {
	res, err := r.dao.CreateNotificationDevice(ctx, database.CreateNotificationDeviceParams{
		UserID:    userId,
		SessionID: sessionId,
		Token:     payload.Token,
		Platform:  payload.Platform,
	})
	if err != nil {
		return model.NotificationDevice{}, err
	}

	return model.NotificationDevice{
		UserID:       res.UserID,
		SessionID:    res.SessionID,
		Token:        res.Token,
		Platform:     res.Platform,
		LastAccessed: res.LastAccessed,
		CreatedAt:    res.CreatedAt,
	}, nil
}

func (r *repo) GetNotificationDevice(ctx context.Context, userId, sessionId uuid.UUID, token, platform string) (model.NotificationDevice, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("user_id", "session_id", "token", "platform", "last_accessed", "created_at")
	sb.From("user_notification_devices")
	andExprs := []string{
		sb.Equal("user_id", userId),
		sb.Equal("session_id", sessionId),
	}
	if token != "" {
		andExprs = append(andExprs, sb.Equal("token", token))
	}
	if platform != "" {
		andExprs = append(andExprs, sb.Equal("platform", platform))
	}
	sb.Where(andExprs...)

	query, args := sb.Build()
	row := r.dao.QueryRow(ctx, query, args...)
	var i database.UserNotificationDevice
	err := row.Scan(
		&i.UserID,
		&i.SessionID,
		&i.Token,
		&i.Platform,
		&i.LastAccessed,
		&i.CreatedAt,
	)
	if err != nil {
		return model.NotificationDevice{}, err
	}

	return model.NotificationDevice{
		UserID:       i.UserID,
		SessionID:    i.SessionID,
		Token:        i.Token,
		Platform:     i.Platform,
		LastAccessed: i.LastAccessed,
		CreatedAt:    i.CreatedAt,
	}, nil
}

func (r *repo) UpdateNotificationDeviceTokenTimestamp(ctx context.Context, userId, sessionId uuid.UUID) error {
	return r.dao.UpdateNotificationDeviceTokenTimestamp(ctx, database.UpdateNotificationDeviceTokenTimestampParams{
		UserID:    userId,
		SessionID: sessionId,
	})
}

func (r *repo) DeleteExpiredTokens(ctx context.Context, interval int32) error {
	return r.dao.DeleteExpiredTokens(ctx, interval)
}

func (r *repo) CreateNotification(ctx context.Context, data *dto.CreateNotification) (model.Notification, error) {
	res, err := r.dao.CreateNotification(ctx, data.ToCreateNotificationDB())
	if err != nil {
		return model.Notification{}, err
	}

	return model.ToNotificationModel(res), nil
}

func (r *repo) GetNotificationsOfUser(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]model.Notification, error) {
	res, err := r.dao.GetNotificationsOfUser(ctx, database.GetNotificationsOfUserParams{
		Limit:  limit,
		Offset: offset,
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
	})
	if err != nil {
		return nil, err
	}

	var notifications []model.Notification
	for _, n := range res {
		notifications = append(notifications, model.ToNotificationModel(n))
	}

	return notifications, nil
}
