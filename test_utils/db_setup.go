package test_utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joelsaunders/bilbo-go/config"
)

func SetUpTestDB(migrationPath string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", "15432", "root", "root", "bilbo",
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Print("db setup failed")
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("db ping fail: %s\n", err)
		return err
	}

	_, err = db.Exec("drop database test")
	_, err = db.Exec("create database test")

	if err != nil {
		log.Printf("could not create test db error: %s \n", err)
		return err
	}

	testdb, err := config.NewDatabase("localhost", "15432", "root", "root", "test")

	if err != nil {
		log.Print("could not connect to new test db")
		return err
	}

	err = config.MigrateDatabase(testdb, migrationPath)

	if err != nil {
		log.Printf("could not migrate test db %s", err)
		return err
	}
	return nil
}
