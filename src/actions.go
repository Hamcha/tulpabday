package main

import (
	"fmt"
	"net/http"
)

func AddTulpa(rw http.ResponseWriter, req *http.Request) {
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, "CAN'T PARSE FORM DATA", 400)
		return
	}

	// Check that all required fields are filled
	vals := []string{"name", "host", "date", "secret"}
	for _, x := range vals {
		if len(req.Form.Get(x)) == 0 {
			http.Error(rw, "REQUIRED FIELDS MISSING", 400)
			return
		}
	}

	// Create tulpa record
	tulpa := Tulpa{
		Name:   req.Form["name"][0],
		Host:   req.Form["host"][0],
		Born:   req.Form["date"][0],
		Secret: req.Form["secret"][0],
	}

	if hasTulpa(tulpa.Host, tulpa.Name) {
		http.Error(rw, "THERE IS ALREADY A TULPA WITH THAT NAME ASSOCIATED TO THAT HOST", 403)
		return
	}

	// Put tulpa record into database
	err = setTulpa(tulpa)

	if err != nil {
		http.Error(rw, "DATABASE ERROR", 500)
		return
	}

	fmt.Fprintf(rw, "OK")
}

func EditTulpa(rw http.ResponseWriter, req *http.Request) {

}

func DeleteTulpa(rw http.ResponseWriter, req *http.Request) {

}
