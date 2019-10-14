package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(url, port, user, password, name string) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		url, port, user, password, name,
	)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	log.Println("successfully connected to db")
	db.SetConnMaxLifetime(time.Minute * 30)
	return db, nil
}

func MigrateDatabase(db *sqlx.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
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
