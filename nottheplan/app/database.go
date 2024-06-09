package app

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// load .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	return os.Getenv(key)
}

func dbConn() (db *sql.DB) {
	// pass the db credentials into variables
	cfg := mysql.Config{
		User:   goDotEnvVariable("DBUSER"),
		Passwd: goDotEnvVariable("DBPASS"),
		DBName: goDotEnvVariable("DBNAME"),
		Addr:   "127.0.0.1:3306",
		Net:    "tcp",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	return db
}
