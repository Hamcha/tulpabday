package main

import (
	"flag"
	"fmt"
	"github.com/cloudflare/gokabinet/kc"
	"github.com/gorilla/mux"
	"net/http"
)

var database *kc.DB

func main() {
	// Get command-line flags
	listen := flag.String("bind", "127.0.0.1:9901", "Address:Port or Socket where to listen to")
	dbdir := flag.String("data", "data.kch", "Directory where to store data")
	flag.Parse()

	// Open database
	var err error
	database, err = kc.Open(*dbdir, kc.WRITE)

	if err != nil {
		panic(err.Error())
	}
	defer database.Close()

	r := mux.NewRouter()
	// GET - Full webpages / UI / JSON Read API
	g := r.Methods("GET").Subrouter()
	g.HandleFunc("/", HomeHandler)
	g.HandleFunc("/json", DumpHandler)
	// POST - Actions
	p := r.Methods("POST").Subrouter()
	p.HandleFunc("/add", AddTulpa)
	p.HandleFunc("/edit", EditTulpa)
	p.HandleFunc("/delete", DeleteTulpa)
	// Admin reserved actions (restricted to localhost)
	g.HandleFunc("/admin/delete/{host}/{tulpa}", AdminDelete)
	g.HandleFunc("/admin/change/{host}/{tulpa}/{key}", AdminChange)

	// Setup and run http server
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Listening on " + *listen)
	http.ListenAndServe(*listen, nil)
}
