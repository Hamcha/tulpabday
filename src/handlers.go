package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"strings"
)

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	// Get all tulpas in the database that we know of
	tulpalist := make([]Tulpa, len(tulpas))
	hosts := make(map[string]bool)
	hostcount := 0
	for i := range tulpas {
		// Format is host.tulpa, split it!
		tdata := strings.Split(tulpas[i], ".")
		tulpalist[i], _ = getTulpa(tdata[0], tdata[1])

		// Check if host is already in hostlist
		if _, ok := hosts[tdata[0]]; !ok {
			hosts[tdata[0]] = true
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
		tulpalist,
		tulpacount,
		hostcount,
	}
	rendered := mustache.RenderFile("template/home.html", data)
	fmt.Fprintf(rw, rendered)
}
