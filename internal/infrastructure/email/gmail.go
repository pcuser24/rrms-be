package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
)

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

const (
	gmailSMTPAuthAddress   = "smtp.gmail.com"
	gmailSMTPServerAddress = "smtp.gmail.com:587"
)

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	contentData any, templateFile string,
	from string,
	to []string,
	cc []string,
	bcc []string,
	replyTo *string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	html, err := html_util.RenderHtml(contentData, templateFile)
	if err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}
	e.HTML = html

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, gmailSMTPAuthAddress)
	return e.Send(gmailSMTPServerAddress, smtpAuth)
}
