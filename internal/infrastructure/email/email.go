package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		contentData any, templateFile string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

// Render HTML byte slice from template file and data.
// data is a struct to hold data for replacing placeholders
// in the template. htmlTemplateFile is the path to the template file.
func renderHtml(data any, htmlTemplateFile string) ([]byte, error) {
	fileName := path.Base(htmlTemplateFile)
	tmpl, err := template.New(fileName).ParseFiles(htmlTemplateFile)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	err = tmpl.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (sender *GmailSender) SendEmail(
	subject string,
	contentData any, templateFile string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	html, err := renderHtml(contentData, templateFile)
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

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
