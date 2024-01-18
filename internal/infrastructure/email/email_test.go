package email

import (
	"os"
	"testing"
)

type Data struct {
	Name          string
	ApplicationId string
	ListingTitle  string
}

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	EmailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	EmailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	EmailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")
	t.Log("EMAIL_SENDER_NAME", EmailSenderName)
	t.Log("EMAIL_SENDER_ADDRESS", EmailSenderAddress)
	t.Log("EMAIL_SENDER_PASSWORD", EmailSenderPassword)

	sender := NewGmailSender(EmailSenderName, EmailSenderAddress, EmailSenderPassword)

	subject := "A test email"
	to := []string{EmailSenderAddress}
	attachFiles := []string{"./email_test.go"}

	err := sender.SendEmail(
		subject,
		Data{
			Name:          "Nguyễn Văn A",
			ApplicationId: "123456",
			ListingTitle:  "Tòa nhà Giang Bắc, Số 1 Thái Hà, Đống Đa, Hà Nội",
		},
		"templates/test.html",
		to, nil, nil, attachFiles)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}

	t.Log("Email sent successfully")
}
