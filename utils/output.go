package utils

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"html/template"
)

func GenerateHTML(data interface{}, fn ...string) (html string, err error) {
	var files []string
	for _, file := range fn {
		files = append(files, fmt.Sprintf("templates/fragments/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	tpl := bytes.Buffer{}

	err = templates.Execute(&tpl, data)

	if err != nil {
		html = ""
		return
	}
	html = tpl.String()
	return
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
