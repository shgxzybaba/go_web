package main

import (
	"fmt"
	"github.com/shgxzybaba/go_web01/data"
	"github.com/shgxzybaba/go_web01/security"
	"github.com/shgxzybaba/go_web01/utils"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	response := utils.Data{}
	students, err := data.FetchAllStudents()
	if err != nil {
		response.Err = err.Error()
		utils.GenerateHTML(w, response, "layout", "navbar", "error")
		return // Exit the function to prevent further processing
	}
	response.Response, response.Err = students, ""

	utils.GenerateHTML(w, response, "layout", "navbar", "index")
}

func main() {

	fmt.Println("Hello server!")
	defer data.ShutDown()

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", security.LoginHandler)
	mux.HandleFunc("/account", security.BasicSecurity(data.AccountHandler))
	mux.HandleFunc("/dashboard", security.BasicSecurity(data.DashboardHandler))

	server := &http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8088",
	}

	if e := server.ListenAndServe(); e != nil {
		fmt.Println("Unable to start server", e)
		return
	}
}
