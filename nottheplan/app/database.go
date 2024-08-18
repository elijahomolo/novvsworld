package app

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DB struct {
	Database *sql.DB
}

// load .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	return os.Getenv(key)
}

func (d *DB) init() error {
	// pass the db credentials into variables
	cfg := mysql.Config{
		User:   goDotEnvVariable("DBUSER"),
		Passwd: goDotEnvVariable("DBPASS"),
		DBName: goDotEnvVariable("DBNAME"),
		Addr:   "127.0.0.1:3306",
		Net:    "tcp",
	}

	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	d.Database = db
	return nil
}

//func dbConn() (db *sql.DB) {
//	// pass the db credentials into variables
//	cfg := mysql.Config{
//		User:   goDotEnvVariable("DBUSER"),
//		Passwd: goDotEnvVariable("DBPASS"),
//		DBName: goDotEnvVariable("DBNAME"),
//		Addr:   "127.0.0.1:3306",
//		Net:    "tcp",
//	}
//
//	db, err := sql.Open("mysql", cfg.FormatDSN())
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

func (d *DB) close() error {
	defer func(Database *sql.DB) {
		err := Database.Close()
		if err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}(d.Database)
	return nil

}

func (d *DB) prepare(query string) (*sql.Stmt, error) {
	stmt, err := d.Database.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}

	return stmt, nil
}

func (d *DB) execute(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}
