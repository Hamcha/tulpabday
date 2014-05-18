package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hamcha/goleg"
	"net/http"
)

var database goleg.Database
var info SiteInfo

func main() {
	// Get command-line flags
	listen := flag.String("bind", "127.0.0.1:9901", "Address:Port or Socket where to listen to")
	dbdir := flag.String("data", "data", "Directory where to store data")
	flag.Parse()

	// Open database
	database = goleg.Open(*dbdir, "tulpa", goleg.F_APPENDONLY|goleg.F_LZ4|goleg.F_SPLAYTREE)
	defer database.CloseSave()

	getInfo()

	r := mux.NewRouter()
	// GET - Full webpages / UI
	g := r.Methods("GET").Subrouter()
	g.HandleFunc("/", HomeHandler)
	g.HandleFunc("/new", NewTulpaHandler)
	g.HandleFunc("/edit", EditTulpaHandler)
	// POST - Actions
	p := r.Methods("POST").Subrouter()
	p.HandleFunc("/add", AddTulpa)
	p.HandleFunc("/edit", EditTulpa)
	p.HandleFunc("/delete", DeleteTulpa)

	// Setup and run http server
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Listening on " + *listen)
	http.ListenAndServe(*listen, nil)
}

func getInfo() {
	// Check if we need to initialize the database
	// and setup SiteInfo
	if !database.Exists("siteinfo") {
		info = SiteInfo{
			TulpaCount: 0,
			HostCount:  0,
		}
		jsoninfo, err := json.Marshal(info)
		if err != nil {
			panic(err.Error())
		}
		database.Jar("siteinfo", jsoninfo)
	} else {
		jsoninfo := database.Unjar("siteinfo")
		err := json.Unmarshal(jsoninfo, &info)
		if err != nil {
			panic(err.Error())
		}
	}
}
