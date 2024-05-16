package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	require.NotEmpty(t, conf.EmailSenderName)
	require.NotEmpty(t, conf.EmailSenderAddress)
	require.NotEmpty(t, conf.EmailSenderPassword)

	sender := NewGmailSender(conf.EmailSenderName, conf.EmailSenderAddress, conf.EmailSenderPassword)

	subject := "A test email"
	to := []string{conf.EmailSenderAddress}
	attachFiles := []string{"./email.go"}

	err := sender.SendEmail(
		subject,
		struct {
			Name          string
			ApplicationId string
			ListingTitle  string
		}{
			Name:          "Nguyễn Văn A",
			ApplicationId: "123456",
			ListingTitle:  "Tòa nhà Giang Bắc, Số 1 Thái Hà, Đống Đa, Hà Nội",
		},
		"templates/test.html",
		conf.EmailSenderAddress,
		to, nil, nil, nil, attachFiles)

	require.NoError(t, err)
}
