package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	name := r.pa

	fmt.Fprintf(w, "Hello, world!\n")
}

func main() {
	fmt.Println("Hello server!")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("/public"))

	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", indexHandler)

	server := &http.Server{
		Handler: mux,
		Addr: "0.0.0.0:8088",
	}
	server.ListenAndServe()
}