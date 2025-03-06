package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

func executeTemplate(templatePath string, templateName string, w http.ResponseWriter, data interface{}) {

	fp := path.Join(templatePath, templateName)
	lp := path.Join(templatePath, "layout.html")
	t, err := template.ParseFiles(fp, lp)
	if err != nil {
		log.Printf("failed to parse template: %v", err)
		errorPage(w, http.StatusInternalServerError)
	}

	if err = t.ExecuteTemplate(w, "layout", &data); err != nil {
		log.Printf("failed to execute template: %v", err)
		errorPage(w, http.StatusInternalServerError)
	}

	return
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	executeTemplate("templates", "index.html", w, nil)
}

type PageData struct {
	Title       string
	Commissions []Commission
}

func CommissionsPage(w http.ResponseWriter, r *http.Request) {

	commission := &Commission{}
	commissions, err := commission.GetAll()

	data := PageData{
		Title: "Commissions",
	}

	if len(commissions) > 0 {
		data.Commissions = commissions
	} else {
		data.Commissions = []Commission{}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list commissions: %v", err), http.StatusInternalServerError)
		return
	}

	executeTemplate("templates", "results.html", w, data)
}

func errorPage(w http.ResponseWriter, status int) {
	//w.WriteHeader(status)
	msg := ""
	var authError bool
	pd := struct {
		Status    int
		Message   string
		AuthError bool
	}{}
	switch {
	case status == http.StatusNotFound:
		msg = "This page could not be found."
		authError = false
	case status == http.StatusUnauthorized:
		msg = "Please log in or sign up to access this page."
		authError = true
	case status == http.StatusInternalServerError:
		msg = "An error occurred. Please try again later."
		authError = false
	}

	pd.Status = status
	pd.Message = msg
	pd.AuthError = authError

	executeTemplate("templates", "error.html", w, pd)
	return
}

func ListUserSubmissions(w http.ResponseWriter, r *http.Request) {
	auth.ValidateToken(w, r)

	mem := &Member{}
	mem.Username = auth.Username
	err := mem.getUserID()
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		errorPage(w, http.StatusInternalServerError)
	}

	sub := &Submission{}
	//submissions, err := sub.getByUserID(mem.ID)
	//if err != nil {
	//	log.Printf("Failed to get submissions: %v", err)
	//	errorPage(w, http.StatusInternalServerError)
	//}

	comms := []Commission{}
	comms, err = sub.listCommissionsByUser(mem.ID)
	for _, comm := range comms {
		log.Printf("comms: %v", comm.ID)

	}
	if err != nil {
		log.Printf("Failed to get commissions: %v", err)
		errorPage(w, http.StatusInternalServerError)
	}

	pd := struct {
		Title       string
		Commissions []Commission
	}{
		Title: "Your Submissions",
	}

	pd.Commissions = comms

	executeTemplate("templates", "results.html", w, pd)
}

func EntryForm(w http.ResponseWriter, r *http.Request) {

	auth.ValidateToken(w, r)

	comm := &Commission{}

	var err error

	// get commission id from url
	commissionID := r.URL.Query().Get("id")
	comm.ID, err = strconv.Atoi(commissionID)
	if err != nil {
		log.Printf("Failed to get commission id: %v", err)
		errorPage(w, http.StatusInternalServerError)
		return
	}

	err = comm.Get()
	if err != nil {
		log.Printf("Failed to get commission: %v", err)
		errorPage(w, http.StatusInternalServerError)
		return
	}

	pD := struct {
		Title              string
		User               string
		CommissionName     string
		CommissionDesc     string
		CommissionLocation string
		CommissionID       int
		Success            bool
		Message            string
	}{
		Title:              `Entry Submission`,
		User:               auth.Username,
		CommissionName:     comm.Name,
		CommissionDesc:     comm.Description,
		CommissionLocation: comm.Location,
		CommissionID:       comm.ID,
		Success:            false,
	}

	mem := &Member{}
	mem.Username = auth.Username
	err = mem.getUserID()
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		errorPage(w, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		details := &Submission{}
		check, err := details.getByCommissionIDAndUserID(comm.ID, mem.ID)
		if err != nil {
			log.Printf("Failed to get submission: %v", err)
			errorPage(w, http.StatusInternalServerError)
		}

		if len(check) > 0 {
			pD.Success = true
			pD.Message = "You have already submitted an entry for this commission"
			executeTemplate("templates", "submission.html", w, pD)
			return
		}

		details = &Submission{
			UserID:       mem.ID,
			CommissionID: comm.ID,
			Demo:         r.FormValue("demo"),
		}

		checkDemo, err := details.checkDemoExists()
		if err != nil {
			log.Printf("Failed to check demo: %v", err)
			errorPage(w, http.StatusInternalServerError)
			return
		}

		if len(checkDemo) > 0 {
			pD.Success = true
			pD.Message = "This demo has already been submitted"
			executeTemplate("templates", "submission.html", w, pD)
			return
		}

		err = details.create()
		if err != nil {
			log.Printf("Failed to create submission: %v", err)
			errorPage(w, http.StatusInternalServerError)
			return
		}

		pD.Success = true
		pD.Message = "Your entry has been submitted"
		executeTemplate("templates", "submission.html", w, pD)

		log.Printf("details: %v", details)

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
