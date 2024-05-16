package email

type EmailSender interface {
	SendEmail(
		subject string,
		contentData any, templateFile string,
		from string,
		to []string,
		cc []string,
		bcc []string,
		replyTo *string,
		attachFiles []string,
	) error
}
