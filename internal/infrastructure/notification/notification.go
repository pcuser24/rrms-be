package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/misc/model"
)

type UserNotificationDevice struct {
	UserID   uuid.UUID `json:"userId"`
	DeviceID string    `json:"deviceId"`
	// The device ID of the device.
	DeviceType   string    `json:"deviceType"`
	Token        string    `json:"token"`
	LastAccessed time.Time `json:"lastAccessed"`
	CreatedAt    time.Time `json:"createdAt"`
}

type NotificationEmailChannel struct {
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	ReplyTo []string `json:"replyTo"`
}

type NotificationPushChannel struct {
	Tokens []string `json:"tokens"`
}

type NotificationTransport struct {
	NM           model.Notification
	EmailChannel *NotificationEmailChannel
	PushChannel  *NotificationPushChannel
}

type NotificationEndpoint interface {
	SendNotification(ctx context.Context, notification *NotificationTransport) []error
}
