package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/shgxzybaba/go_web01/data"
)

type Data struct {
	Response interface{}
	err string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	response := Data{}
	students, err := data.FetchAllStudents()
	if err != nil {
		response.err = err.Error()
		generateHTML(w, response, "layout", "navbar", "content", "error")
        return // Exit the function to prevent further processing
	}
	response.Response, response.err = students, ""

	generateHTML(w, response, "layout", "navbar", "content", "error")
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
	fmt.Println("Setting up connection to database")
	err := data.Setup()
	defer data.ShutDown()
	if err != nil {
		fmt.Println("Could not open database!", err)
		return
	}

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