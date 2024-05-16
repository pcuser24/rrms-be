package email

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
)

type ResendSender struct {
	apiKey string
	client *resend.Client
}

func NewResendSender(apiKey string) EmailSender {
	return &ResendSender{
		apiKey: apiKey,
		client: resend.NewClient(apiKey),
	}
}

func (s *ResendSender) SendEmail(
	subject string,
	contentData any, templateFile string,
	from string,
	to []string,
	cc []string,
	bcc []string,
	replyTo *string,
	attachFiles []string,
) error {
	html, err := html_util.RenderHtml(contentData, templateFile)
	if err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Cc:      cc,
		Bcc:     bcc,
		Html:    string(html),
	}
	if replyTo != nil {
		params.ReplyTo = *replyTo
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Println("Email sent: ", sent.Id)
	return nil
}
