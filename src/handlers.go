package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
)

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
	}
	// Get total tulpas count
	tulpacount := len(tulpas)

	data := struct {
		Tulpas     []Tulpa
		TulpaCount int
		HostCount  int
	}{
		tulpas,
		tulpacount,
		hostcount,
	}
	rendered := mustache.RenderFile("template/home.html", data)
	fmt.Fprintf(rw, rendered)
}
