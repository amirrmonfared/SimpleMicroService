package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var counts int64

func TestMain(m *testing.M) {
	log.Println("starting authentication service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

func openDB(dsn string) (*sql.DB, error) {
	testDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = testDB.Ping()
	if err != nil {
		return nil, err
	}

	return testDB, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
