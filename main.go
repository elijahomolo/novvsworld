package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Title  string
	Author string
}

type Entry struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Demo      string
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/contest", entryForm)
	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

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

func executeTemplate(templatePath string, templateName string, w http.ResponseWriter, data interface{}) error {
	fp := path.Join(templatePath, templateName)
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func entryForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/signup/form.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Entry{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Demo:      r.FormValue("demo"),
	}

	// send details to a database
	//prepare a query to insert the data into the database
	db := dbConn()

	insForm, err := db.Prepare(`INSERT INTO entries(first_name, last_name, email, demo,  created_on) VALUES (?,?, ?, ?, ?)`)
	if err != nil {
		// an error has occurred
		panic(err)
	}
	//execute the query using the form data
	now := time.Now().UTC()
	_, err = insForm.Exec(details.FirstName, details.LastName, details.Email, details.Demo, now)

	defer db.Close()

	tmpl.Execute(w, struct{ Success bool }{true})

}
