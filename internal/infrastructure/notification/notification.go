package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	Data         map[string]interface{} `json:"data"`
	EmailChannel *NotificationEmailChannel
	PushChannel  *NotificationPushChannel
}

type NotificationEndpoint interface {
	SendNotification(ctx context.Context, notification *NotificationTransport) []error
}
