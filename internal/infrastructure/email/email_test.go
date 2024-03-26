package email

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/utils"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
)

type TestSendEmailConfig struct {
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME" validate:"required"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS" validate:"required"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD" validate:"required"`
}

var conf TestSendEmailConfig

func TestMain(t *testing.M) {
	viper.AddConfigPath(utils.GetBasePath())
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Fatalf("invalid or missing fields in config file: %v", err)
	}

	os.Exit(t.Run())
}

func TestRenderHtml(t *testing.T) {
	data := struct {
		Date          html_util.HTMLTime
		Name          string
		ApplicationId string
		ListingTitle  string
	}{
		Date:          html_util.NewHTMLTime(time.Now()),
		Name:          "Tehc's School",
		ApplicationId: "123456",
		ListingTitle:  "A test listing",
	}
	templateFile := "templates/test.html"

	buf, err := html_util.RenderHtml(data, templateFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
}

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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
		to, nil, nil, attachFiles)

	require.NoError(t, err)
}
