package app

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"log"
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	credentials := AdminUser{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	err := credentials.Verify()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify user: %v", err), http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = createSession(w, credentials.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}
}

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin/createUser.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	admin := &AdminUser{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	err := admin.Create()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create admin: %v", err), http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}
}

func CreateCommission(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/admin/create-commission.html"))

	validateSession(w, r)

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	commission := &Commission{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Budget:      r.FormValue("budget"),
		Currency:    r.FormValue("currency"),
		Location:    r.FormValue("location"),
		Deadline:    r.FormValue("deadline"),
		Status:      r.FormValue("status"),
		Winner:      r.FormValue("winner"),
	}

	log.Printf("Commission: %v", commission)

	err := commission.Create()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create commission: %v", err), http.StatusInternalServerError)
		tmpl.Execute(w, nil)
		return
	}

	err = tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}

	//create a commission
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	//write, err := w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	if err != nil {
		return
	}
}
