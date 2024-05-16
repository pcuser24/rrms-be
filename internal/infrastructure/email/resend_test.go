package email

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
)

func TestSendEmailWithResend(t *testing.T) {
	require.NotEmpty(t, conf.ResendAPIKey)

	sender := NewResendSender(conf.ResendAPIKey)

	subject := "A test email"
	require.NotEmpty(t, conf.EmailSenderAddress)
	to := []string{conf.EmailSenderAddress}
	attachFiles := []string{"./email.go"}

	err := sender.SendEmail(
		subject,
		struct {
			Date          html_util.HTMLTime
			Name          string
			ApplicationId string
			ListingTitle  string
		}{
			Date:          html_util.NewHTMLTime(time.Now()),
			Name:          "Nguyễn Văn A",
			ApplicationId: "123456",
			ListingTitle:  "Tòa nhà Giang Bắc, Số 1 Thái Hà, Đống Đa, Hà Nội",
		},
		"templates/test.html",
		"rrms@resend.dev",
		to, nil, nil, nil, attachFiles)

	require.NoError(t, err)
}
