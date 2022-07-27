package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Index", nil)
}

var tmpl = template.Must(template.ParseGlob("assets/templates/*"))

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	//	http.Handle("/assets/static", http.StripPrefix("/assets/static", http.FileServer(http.Dir("assets/static"))))
	http.HandleFunc("/", index)
	fmt.Println("server starting on port 3000...")
	http.ListenAndServe(":3000", nil)
}
