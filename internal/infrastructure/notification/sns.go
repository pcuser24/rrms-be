package notification

import (
	"context"
	"maps"

	managed_sns "github.com/user2410/rrms-backend/internal/infrastructure/aws/sns"
)

type snsNotificationEndpoint struct {
	snsClient                 managed_sns.SNSClient
	emailNotificationTopicArn string
	pushNotificationTopicArn  string
}

func NewSNSNotificationEndpoint(
	snsClient managed_sns.SNSClient,
	emailNotificationTopicArn, pushNotificationTopicArn string,
) NotificationEndpoint {
	return &snsNotificationEndpoint{
		snsClient:                 snsClient,
		emailNotificationTopicArn: emailNotificationTopicArn,
		pushNotificationTopicArn:  pushNotificationTopicArn,
	}
}

func (s *snsNotificationEndpoint) SendNotification(ctx context.Context, notification *NotificationTransport) []error {
	var errs []error = make([]error, 0, 3)

	// send notifications
	if notification.EmailChannel != nil {
		err := s.sendEmailNotification(ctx, notification)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if notification.PushChannel != nil {
		err := s.sendPushNotification(ctx, notification)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (s *snsNotificationEndpoint) sendEmailNotification(ctx context.Context, notification *NotificationTransport) error {
	msgAttributes := make(map[string]interface{})
	maps.Copy(msgAttributes, notification.Data)
	maps.Copy(msgAttributes, map[string]interface{}{
		"to":      notification.EmailChannel.To,
		"cc":      notification.EmailChannel.CC,
		"bcc":     notification.EmailChannel.BCC,
		"replyTo": notification.EmailChannel.ReplyTo,
	})

	return s.snsClient.Publish(ctx, notification.Title, notification.Content, s.emailNotificationTopicArn, msgAttributes)
	// log.Printf("Email notification sent %v, with attributes: %v", *notification, msgAttributes)
	// return nil
}

func (s *snsNotificationEndpoint) sendPushNotification(ctx context.Context, notification *NotificationTransport) error {
	msgAttributes := make(map[string]interface{})
	maps.Copy(msgAttributes, notification.Data)
	maps.Copy(msgAttributes, map[string]interface{}{
		"tokens": notification.PushChannel.Tokens,
	})

	return s.snsClient.Publish(ctx, notification.Title, notification.Content, s.pushNotificationTopicArn, msgAttributes)
	// log.Printf("Push notification sent %v, with attributes: %v", *notification, msgAttributes)
	// return nil
}
