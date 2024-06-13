package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

type Entry struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Demo      string
	createdOn []uint8
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	err := executeTemplate("templates", "index.html", w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CommisionsPage(w http.ResponseWriter, r *http.Request) {

	commission := &Commission{}

	commissions, err := commission.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list commissions: %v", err), http.StatusInternalServerError)
		return
	}

	err = executeTemplate("templates", "commissions.html", w, commissions)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list commissions: %v", err), http.StatusInternalServerError)
		return
	}
}

//func testTable(w http.ResponseWriter, r *http.Request) {
//	db := dbConn()
//	selDB, err := db.Query("SELECT * FROM entries ORDER BY created_on DESC")
//	if err != nil {
//		panic(err)
//	}
//	//create a slice of Entry objects to store the data
//	entries := []Entry{}
//	//loop through the rows returned by the query
//	for selDB.Next() {
//		//create a new Entry object
//		entry := Entry{}
//		//scan the data into the object
//		err = selDB.Scan(&entry.FirstName, &entry.LastName, &entry.Email, &entry.Demo, &entry.createdOn)
//		if err != nil {
//			panic(err)
//		}
//		//append the object to the slice
//		entries = append(entries, entry)
//	}
//	//close the database connection
//	defer db.Close()
//	//execute the template
//	//tmpl := template.Must(template.ParseFiles("templates/signup/results.html"))
//
//	//t := template.Must(template.ParseFiles("templates/signup/results.html"))
//
//	//if err := t.Execute(w, &entries); err != nil {
//	//	log.Fatal(err)
//	//}
//
//	err = executeTemplate("templates/signup", "results.html", w, &entries)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//func entryForm(w http.ResponseWriter, r *http.Request) {
//	tmpl := template.Must(template.ParseFiles("templates/signup/form.html"))
//
//	if r.Method != http.MethodPost {
//		tmpl.Execute(w, nil)
//		return
//	}
//
//	details := Entry{
//		FirstName: r.FormValue("first_name"),
//		LastName:  r.FormValue("last_name"),
//		Email:     r.FormValue("email"),
//		Demo:      r.FormValue("demo"),
//	}
//
//	// send details to a database
//	//prepare a query to insert the data into the database
//	db := dbConn()
//
//	insForm, err := db.Prepare(`INSERT INTO entries(first_name, last_name, email, demo,  created_on) VALUES (?,?, ?, ?, ?)`)
//	if err != nil {
//		// an error has occurred
//		panic(err)
//	}
//	//execute the query using the form data
//	now := time.Now().UTC()
//	_, err = insForm.Exec(details.FirstName, details.LastName, details.Email, details.Demo, now)
//
//	defer db.Close()
//
//	tmpl.Execute(w, struct{ Success bool }{true})
//
//}

//func showResults(w http.ResponseWriter, r *http.Request) {
//
//	db := dbConn()
//	selDB, err := db.Query("SELECT * FROM entries ORDER BY id DESC")
//	if err != nil {
//		panic(err)
//	}
//	//create a slice of Entry objects to store the data
//	entries := []Entry{}
//	//loop through the rows returned by the query
//	for selDB.Next() {
//		//create a new Entry object
//		entry := Entry{}
//		//scan the data into the object
//		err = selDB.Scan(&entry.ID, &entry.FirstName, &entry.LastName, &entry.Email, &entry.Demo)
//		if err != nil {
//			panic(err)
//		}
//		//append the object to the slice
//		entries = append(entries, entry)
//	}
//	//close the database connection
//	defer db.Close()
//	//execute the template
//	executeTemplate("templates/signup", "results.html", w, entries)
//}

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
