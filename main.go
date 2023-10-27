package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Person struct {
	Name string
	Age uint16
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	person := Person{
		Name: "Akinduro",
		Age: 30,
	}

	generateHTML(w, person, "layout", "navbar", "content")
}

func generateHTML(w http.ResponseWriter, data interface{}, fn ...string) {
	var files []string
	for _, file := range fn {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

func main() {
	fmt.Println("Hello server!")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))

	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", indexHandler)

	server := &http.Server{
		Handler: mux,
		Addr: "0.0.0.0:8088",
	}
	server.ListenAndServe()
}