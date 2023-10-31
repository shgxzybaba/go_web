package utils

import (
	"fmt"
	"net/http"
)

func BasicErrorHandle(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "An internal server error occured for %s\n", r.RequestURI)
	if err != nil {
		panic(err)
	}
}
