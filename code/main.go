package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "<h1>Hello World</h1>")
}

func main(){
	http.HandleFunc("/", index)
	fmt.Println("server starting on port 3000...")
	http.ListenAndServe(":3000", nil)
}