package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type TestRepoConfig struct {
	DatabaseURL    string `mapstructure:"DB_TEST_URL" validate:"required,uri"`
	DBMigrationDir string `mapstructure:"DB_MIGRATION_DIR" validate:"required"`

	// Redis
	RedisAddr     string `mapstructure:"REDIS_TEST_ADDR" validate:"required"`
	RedisPassword string `mapstructure:"REDIS_TEST_PASSWORD" validate:"omitempty"`
	RedisDB       int    `mapstructure:"REDIS_TEST_DB" validate:"required"`
}

func NewTestRepoConfig(envPath string) (*TestRepoConfig, error) {
	var conf TestRepoConfig
	viper.AddConfigPath(envPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		return nil, fmt.Errorf("invalid or missing fields in config file: %v", err)
	}

	return &conf, nil
}
