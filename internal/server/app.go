package server

import (
	"database/sql"
	"example.com/go-rest-api/internal/handlers/registration_handlers"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"

	_ "github.com/lib/pq"
)

func Start() {
	log.Println("Starting server...")

	db := getDB()
	router := gin.Default()
	router.POST("/api/register", registration_handlers.RegisterStudentsHandler(db))
	router.Run()
}

func getDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

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
