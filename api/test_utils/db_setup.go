package test_utils

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joelsaunders/blog-go/api/pkg/repository/postgres"
)

func OpenTransaction(t *testing.T) *sqlx.DB {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, err := sqlx.Open("txdb", cName)
	if err != nil {
		t.Fatal("could not open db")
	}
	return db
}

func SetUpTestDB(migrationPath string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", "15432", "root", "root", "blog",
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

	testdb, err := postgres.NewDatabase("localhost", "15432", "root", "root", "test")

	if err != nil {
		log.Print("could not connect to new test db")
		return err
	}

	err = postgres.MigrateDatabase(testdb, migrationPath)

	if err != nil {
		log.Printf("could not migrate test db %s", err)
		return err
	}
	return nil
}
