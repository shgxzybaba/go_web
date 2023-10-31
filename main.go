package main

import (
	"fmt"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/security"
	"html/template"
	"net/http"
)

type Data struct {
	Response interface{}
	Err      string
}

type User struct {
	Username string
	Id       int
}

func (u *User) user(s data.Student) {
	u.Username = s.Username
	u.Id = int(s.Id)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	response := Data{}
	students, err := data.FetchAllStudents()
	if err != nil {
		response.Err = err.Error()
		generateHTML(w, response, "layout", "navbar", "content", "error")
		return // Exit the function to prevent further processing
	}
	response.Response, response.Err = students, ""

	generateHTML(w, response, "layout", "navbar", "content", "error")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		response := Data{nil, ""}
		generateHTML(w, response, "layout", "navbar", "content", "login") //todo: make "content" dynamic
	} else if method == "POST" {
		user, password := r.FormValue("username"), r.FormValue("password")
		student, sess, err := security.Login(user, password)
		var response = Data{}
		if err != nil {
			response.Response = nil
			response.Err = err.Error()
		} else {
			u := User{}
			u.user(student)
			response.Response = u
			response.Err = ""
			w.Header().Set("session-id", sess.Uuid)
		}

		generateHTML(w, response, "layout", "navbar", "content", "dashboard", "error")
	}

}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

func generateHTML(w http.ResponseWriter, data interface{}, fn ...string) {
	var files []string
	for _, file := range fn {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	if err1 := templates.ExecuteTemplate(w, "layout", data); err1 != nil {
		if _, err := fmt.Fprintln(w, "An error occurred while loading template", err1); err != nil {
			fmt.Println("Could not write to response writer", err1)
		}
	}
}

//todo: create a session validation filter

func main() {
	fmt.Println("Hello server!")
	defer data.ShutDown()

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))

	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/account", security.BasicSecurity(accountHandler))

	server := &http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8088",
	}
	if e := server.ListenAndServe(); e != nil {
		fmt.Println("Unable to start server", e)
		return
	}
}
