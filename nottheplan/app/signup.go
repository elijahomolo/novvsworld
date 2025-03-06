package app

import (
	"fmt"
	"html/template"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//load the signup page

	tmpl := template.Must(template.ParseFiles("templates/signup/form.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := &Member{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Instagram: r.FormValue("instagram"),
		Twitter:   r.FormValue("twitter"),
		Threads:   r.FormValue("threads"),
		Country:   r.FormValue("country"),
	}

	err := details.CreateUser()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}

}
