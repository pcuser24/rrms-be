package seed

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type seedConfig struct {
	DatabaseURL    string `mapstructure:"DB_URL" validate:"required,uri"`
	DBMigrationDir string `mapstructure:"DB_MIGRATION_DIR" validate:"required"`
}

type seedCommand struct {
	*cobra.Command
	config *seedConfig
}

func NewSeedCommand() *seedCommand {
	c := &seedCommand{}
	c.Command = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with initial data",
		Run:   c.run,
	}
	c.config = newSeedConfig(c.Command)
	return c
}

func newSeedConfig(cmd *cobra.Command) *seedConfig {
	var conf seedConfig
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Failed to unmarshal config file:", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Fatal("Invalid or missing fields in config file: ", err)
	}

	return &conf
}

func (c *seedCommand) run(cmd *cobra.Command, args []string) {
}
