package utils

import (
	"fmt"
	"github.com/google/uuid"
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

func (response *Data) ErrorResponse(err error) {
	response.Response = nil
	response.Err = err.Error()
}

func (response *Data) DataResponse(data any) {
	response.Response = data
	response.Err = ""
}

func GenerateUUID() string {

	uuidObj := uuid.New()
	uuidString := uuidObj.String()

	return uuidString
}
