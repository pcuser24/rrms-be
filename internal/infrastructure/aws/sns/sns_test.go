package sns

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	managed_aws "github.com/user2410/rrms-backend/internal/infrastructure/aws"
	"github.com/user2410/rrms-backend/internal/utils"
)

type TestNotificationConfig struct {
	AWSRegion                       string  `mapstructure:"AWS_REGION" validate:"required"`
	AWSSNSEmailNotificationTopicArn string  `mapstructure:"AWS_SNS_EMAIL_NOTIFICATION_TOPIC_ARN" validate:"required"`
	AWSSNSPushNotificationTopicArn  string  `mapstructure:"AWS_SNS_PUSH_NOTIFICATION_TOPIC_ARN" validate:"required"`
	AWSEndpoint                     *string `mapstructure:"AWS_ENDPOINT" validate:"omitempty"`
}

var (
	basePath = utils.GetBasePath()
	conf     TestNotificationConfig

	awsConf *aws.Config
	c       SNSClient
)

func TestMain(m *testing.M) {
	viper.AddConfigPath(basePath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("failed to read config file:", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalln("failed to unmarshal config file:", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Fatalln("invalid or missing fields in config file:", err)
	}

	awsConf, err = managed_aws.LoadConfig(conf.AWSRegion, conf.AWSEndpoint)
	if err != nil {
		log.Fatalln("failed to load aws config:", err)
	}

	c = NewSNSClient(awsConf)

	os.Exit(m.Run())
}
