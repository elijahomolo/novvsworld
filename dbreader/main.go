package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	//"net/smtp"
	"os"
	_ "github.com/lib/pq"
	_ "github.com/lib/pq"
)

// define a user model
type User struct {
	Id    int
	Email string
	SocialNetwork string
	Handle string
}

// load .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load("prod.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// connect to the database and return it as an object
func dbConn() (db *sql.DB) {
	// pass the db credentials into variables
	host := goDotEnvVariable("DBHOST")
	port := goDotEnvVariable("DBPORT")
	dbUser := goDotEnvVariable("DBUSER")
	dbPass := goDotEnvVariable("DBPASS")
	dbname := goDotEnvVariable("DBNAME")
//	sslmode := goDotEnvVariable("SSLMODE")
//	caCert := "sslrootcert = " + goDotEnvVariable("CA_CERT")
	// create a connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		host, port, dbUser, dbPass, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}


func readDB()  {
	db := dbConn()
	rows, err := db.Query(`SELECT "email", "social_network", "social_handle" FROM public."users"`)
	CheckError(err)
	usr := User{}
	res := []User{}
	for rows.Next() {
		var email, social_network, social_handle string
		err = rows.Scan(&email, &social_network, &social_handle)
		CheckError(err)

		usr.Email = email
		usr.SocialNetwork = social_network
		usr.Handle = social_handle
		res = append(res, usr)
	}
	fmt.Println(res)
	defer db.Close()
}


func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}


func main() {
	readDB()
}
