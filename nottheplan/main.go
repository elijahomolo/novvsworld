package main

import (
	"github.com/elijahomolo/novvsworld/nottheplan/app"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc("/contest", entryForm)
	//http.HandleFunc("/results", showResults)
	//http.HandleFunc("/test", testTable)
	http.HandleFunc("/", app.HomePage)
	http.HandleFunc("/commissions", app.CommissionsPage)
	http.HandleFunc("/commissions/list-my-submissions", app.ListUserSubmissions)
	http.HandleFunc("/commissions/entry", app.EntryForm)
	http.HandleFunc("/signup", app.Signup)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/welcome", app.Welcome)
	http.HandleFunc("/admin/create-commission", app.CreateCommission)
	http.HandleFunc("/admin/login", app.AdminLogin)
	http.HandleFunc("/admin/create-user", app.CreateAdmin)
	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
