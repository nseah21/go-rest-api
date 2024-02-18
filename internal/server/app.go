package server

import (
	"log"
	"fmt"
	"database/sql"
	"os"
	"github.com/lpernett/godotenv"
	
	_ "github.com/lib/pq"
)

func Start() {
	log.Println("Starting server...")
}

func getDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	host      := os.Getenv("DB_HOST")
	port      := os.Getenv("DB_PORT")
	user      := os.Getenv("DB_USER")
	password  := os.Getenv("DB_PASSWORD")
	dbname    := os.Getenv("DB_NAME")

	formatStr := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	connStr := fmt.Sprintf(formatStr, host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
