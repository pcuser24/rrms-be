package email

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/internal/utils"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
)

type TestSendEmailConfig struct {
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME" validate:"omitempty"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS" validate:"omitempty"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD" validate:"omitempty"`

	ResendAPIKey string `mapstructure:"RESEND_API_KEY" validate:"omitempty"`
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
