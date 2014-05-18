package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
)

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	data := struct {
		Info SiteInfo
	}{
		info,
	}
	rendered := mustache.RenderFileInLayout("template/home.html", "template/layout.html", data)
	fmt.Fprintf(rw, rendered)
}

func NewTulpaHandler(rw http.ResponseWriter, req *http.Request) {

}

func EditTulpaHandler(rw http.ResponseWriter, req *http.Request) {

}
