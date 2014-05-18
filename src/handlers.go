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
	rendered := mustache.RenderFile("template/home.html", data)
	fmt.Fprintf(rw, rendered)
}
