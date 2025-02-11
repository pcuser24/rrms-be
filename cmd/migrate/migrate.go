package migrate

import (
	"fmt"
	"log"
	"os"

	"github.com/user2410/rrms-backend/cmd/version"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type migrateConfig struct {
	DatabaseURL    string `mapstructure:"DB_URL" validate:"required,uri"`
	DBMigrationDir string `mapstructure:"DB_MIGRATION_DIR" validate:"required"`
}

type migrateCommand struct {
	*cobra.Command
}

func NewMigrateCommand() *migrateCommand {
	c := &migrateCommand{}
	c.Command = &cobra.Command{
		Use:   "migrate",
		Short: fmt.Sprintf("Database migration manager for %s", version.ReadableName),
		Long: fmt.Sprintf(`%s
			Manage the database migrations for %s from the command line`, version.Art(), version.ReadableName),
		Run: c.run,
	}
	c.Command.AddCommand(
		newMigrateUpCommand().Command,
		newMigrateDownCommand().Command,
	)
	return c
}

func (c *migrateCommand) run(cmd *cobra.Command, args []string) {
	c.Help()
}

func initMigrator(c *migrateConfig) (database.Migrator, error) {
	m, err := database.NewMigrator(c.DBMigrationDir, c.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create the migration manager: %w", err)
	}
	return m, nil
}

func newMigrationConfig(cmd *cobra.Command) *migrateConfig {
	var conf migrateConfig
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

// migrate up command

type migrateUpCommand struct {
	*cobra.Command
	config *migrateConfig
}

func newMigrateUpCommand() *migrateUpCommand {
	c := &migrateUpCommand{}
	c.Command = &cobra.Command{
		Use:   "up",
		Short: "Upgrade a migration on the permission API database",
		Run:   c.run,
	}
	c.config = newMigrationConfig(c.Command)
	return c
}

func (c *migrateUpCommand) run(cmd *cobra.Command, args []string) {
	m, err := initMigrator(c.config)
	if err != nil {
		log.Fatal(err.Error())
	}
	res, err := m.Upgrade()
	if err != nil {
		log.Fatal(err.Error())
	}
	if !res {
		log.Println("The database is already up to date")
	}
	log.Println("The database has been upgraded successfully")
}

// migrate down command

type migrateDownCommand struct {
	*cobra.Command
	config *migrateConfig
}

func newMigrateDownCommand() *migrateDownCommand {
	c := &migrateDownCommand{}
	c.Command = &cobra.Command{
		Use:   "down",
		Short: "Downgrade a migration on the permission API database",
		Run:   c.run,
	}
	c.config = newMigrationConfig(c.Command)
	return c
}

func (c *migrateDownCommand) run(cmd *cobra.Command, args []string) {
	m, err := initMigrator(c.config)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	res, err := m.Downgrade()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if !res {
		log.Println("The database is already up to date")
	}
	log.Println("The database has been downgraded successfully")
}
