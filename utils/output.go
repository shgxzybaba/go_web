package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

func GenerateHTML(w http.ResponseWriter, data interface{}, fn ...string) {
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

type Data struct {
	Response interface{}
	Err      string
}
