package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func AddTulpa(rw http.ResponseWriter, req *http.Request) {
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, "Can't parse form data", 400)
		return
	}

	// Check that all required fields are filled
	vals := []string{"name", "host", "date", "secret"}
	for _, x := range vals {
		if len(req.Form.Get(x)) == 0 {
			http.Error(rw, "Required fields are missing", 400)
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
		http.Error(rw, "There is already a tulpa with that name and that host", 403)
		return
	}

	// Put tulpa record into database
	err = setTulpa(tulpa)

	if err != nil {
		http.Error(rw, "Database error", 500)
		return
	}

	fmt.Fprintf(rw, "Ok")
}

func EditTulpa(rw http.ResponseWriter, req *http.Request) {
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, "Can't parse form data", 400)
		return
	}

	// Check that all required fields are filled
	vals := []string{"name", "host", "date", "secret", "oldname", "oldhost"}
	for _, x := range vals {
		if len(req.Form.Get(x)) == 0 {
			http.Error(rw, "Required fields are missing", 400)
			return
		}
	}

	host := req.Form["host"][0]
	name := req.Form["name"][0]
	born := req.Form["date"][0]
	secret := req.Form["secret"][0]
	oldname := req.Form["oldname"][0]
	oldhost := req.Form["oldhost"][0]

	// Check if provided tulpa exists
	if !hasTulpa(oldhost, oldname) {
		http.Error(rw, "Trying to edit an inexistant tulpa", 404)
		return
	}

	// Check secret code
	tulpa, err := getTulpa(oldhost, oldname)
	if err != nil {
		http.Error(rw, "Database error", 500)
	}
	if tulpa.Secret != secret {
		http.Error(rw, "Wrong secret code", 403)
		return
	}

	// Set new tulpa data
	tulpa.Name = name
	tulpa.Host = host
	tulpa.Born = born

	// Take old tulpa out of the database
	err = removeTulpa(oldhost, oldname)
	// Put new tulpa record into database
	err = setTulpa(tulpa)

	if err != nil {
		http.Error(rw, "Database error", 500)
		return
	}

	fmt.Fprintf(rw, "Ok")
}

func DeleteTulpa(rw http.ResponseWriter, req *http.Request) {
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, "Can't parse form data", 400)
		return
	}

	// Check that all required fields are filled
	vals := []string{"name", "host", "secret"}
	for _, x := range vals {
		if len(req.Form.Get(x)) == 0 {
			http.Error(rw, "Required fields are missing", 400)
			return
		}
	}

	host := req.Form["host"][0]
	name := req.Form["name"][0]
	secret := req.Form["secret"][0]
	// Check if provided tulpa exists
	if !hasTulpa(host, name) {
		http.Error(rw, "Trying to delete an inexistant tulpa", 404)
		return
	}

	// Check secret code
	tulpa, err := getTulpa(host, name)
	if err != nil {
		http.Error(rw, "Database error", 500)
	}
	if tulpa.Secret != secret {
		http.Error(rw, "Wrong secret code", 403)
		return
	}

	removeTulpa(host, name)
	fmt.Fprintf(rw, "Ok")
}

func AdminDelete(rw http.ResponseWriter, req *http.Request) {
	// Check for forwarded ip (if there is a reverse proxy)
	if ipfwd, ok := req.Header["X-Forwarded-For"]; ok && ipfwd[0] != "127.0.0.1" {
		http.Error(rw, "Only the host machine can use admin commands", 403)
		return
	}

	vars := mux.Vars(req)

	err := removeTulpa(vars["host"], vars["tulpa"])
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	fmt.Fprintf(rw, "Ok")
}

func AdminChange(rw http.ResponseWriter, req *http.Request) {
	// Check for forwarded ip (if there is a reverse proxy)
	if val, ok := req.Header["X-Forwarded-For"]; ok && val[0] != "127.0.0.1" {
		http.Error(rw, "Only the host machine can use admin commands", 403)
		return
	}

	vars := mux.Vars(req)

	tulpa, err := getTulpa(vars["host"], vars["tulpa"])
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	tulpa.Secret = vars["key"]

	err = setTulpa(tulpa)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	fmt.Fprintf(rw, "Ok")
}
