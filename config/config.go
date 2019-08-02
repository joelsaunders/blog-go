package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	viper "github.com/spf13/viper"
)

type Constants struct {
	PORT     string
	Postgres struct {
		URL        string
		DBNAME     string
		DBPORT     string
		DBUSER     string
		DBPASSWORD string
	}
}

type Config struct {
	Constants
	Database *sql.DB
}

func initViper() (Constants, error) {
	viper.SetConfigName("bilbo.config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		return Constants{}, err
	}

	viper.SetDefault("PORT", "8000")
	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, err
}

func NewConfig() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants

	if err != nil {
		return &config, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Postgres.URL, config.Postgres.DBPORT,
		config.Postgres.DBUSER, config.Postgres.DBUSER, config.Postgres.DBNAME)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return &config, err
	}

	err = db.Ping()

	if err != nil {
		return &config, err
	}

	log.Println("successfully connected to db")

	config.Database = db
	return &config, nil
}

func MigrateDatabase(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		log.Println("No migrations applied")
		return nil
	} else if err != nil {
		return err
	}

	log.Println("Migrations successful")
	return nil
}
