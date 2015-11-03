package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
)

func DumpHandler(rw http.ResponseWriter, req *http.Request) {
	// Get all tulpas in the database that we know of
	tulpas, err := getAllTulpas()
	if err != nil {
		http.Error(rw, "Database error", 500)
		return
	}

	for i := range tulpas {
		// Strip secret keys
		tulpas[i].Secret = ""
	}

	// Put tulpas into JSON format
	err = json.NewEncoder(rw).Encode(tulpas)
	if err != nil {
		http.Error(rw, "Encoding error", 500)
		return
	}
}

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	// Get all tulpas in the database that we know of
	tulpas, err := getAllTulpas()
	if err != nil {
		http.Error(rw, "Database error", 500)
		return
	}

	hosts := make(map[string]bool)
	hostcount := 0
	for i := range tulpas {
		// Check if host is already in hostlist
		if _, ok := hosts[tulpas[i].Host]; !ok {
			hosts[tulpas[i].Host] = true
			hostcount += 1
		}
		// Strip secret keys
		tulpas[i].Secret = ""
	}

	// Get total tulpas count
	tulpacount := len(tulpas)

	// Put tulpas into JSON format
	jsondata, err := json.Marshal(tulpas)
	if err != nil {
		http.Error(rw, "Encoding error", 500)
		return
	}

	data := struct {
		Tulpas     string
		TulpaCount int
		HostCount  int
	}{
		string(jsondata),
		tulpacount,
		hostcount,
	}
	rendered := mustache.RenderFile("template/home.html", data)
	fmt.Fprintln(rw, rendered)
}
