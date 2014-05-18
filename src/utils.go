package main

import (
	"encoding/json"
)

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

	// Check if we need to initialize the list of tulpas
	// or get it if it already exists
	if !database.Exists("tulpalist") {
		tulpas = make([]string, 0)
		jsontp, err := json.Marshal(tulpas)
		if err != nil {
			panic(err.Error())
		}
		database.Jar("tulpalist", jsontp)
	} else {
		jsontp := database.Unjar("tulpalist")
		err := json.Unmarshal(jsontp, &tulpas)
		if err != nil {
			panic(err.Error())
		}
	}
}

func hasTulpa(hname string, tname string) bool {
	return database.Exists("tulpa." + hname + "." + tname)
}

func getTulpa(hname string, tname string) (Tulpa, error) {
	var out Tulpa
	jsondata := database.Unjar("tulpa." + hname + "." + tname)
	err := json.Unmarshal(jsondata, &out)
	return out, err
}

func getTulpasByHost(hname string) (Tulpa[], error) {
	// Check if the host exists
	if !database.Exists("host."+hname) {
		return nil, new Error("Host not found")
	}

	// Retrieve tulpa list (names)
	var hostval []string
	hostdata := database.Unjar("host."+hname)
	err := json.Unmarshal(hostdata, &hostval)
	if err != nil {
		return nil, err
	}

	// Get all the tulpas
	tulpas := make([]Tulpa, len(hostval))
	for i := range hostval {
		tulpas[i], err = getTulpa(hname, hostval[i])
		if err != nil {
			return tulpas, err
		}
	}

	return tulpas
}
