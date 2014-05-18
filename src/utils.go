package main

import (
	"encoding/json"
	"errors"
)

func getInfo() {
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

func hasTulpa(hname, tname string) bool {
	return database.Exists("tulpa." + hname + "." + tname)
}

func getTulpa(hname string, tname string) (Tulpa, error) {
	var out Tulpa
	jsondata := database.Unjar("tulpa." + hname + "." + tname)
	err := json.Unmarshal(jsondata, &out)
	return out, err
}

func getTulpaListByHost(hname string) ([]string, error) {

	// Check if the host exists
	if !database.Exists("host." + hname) {
		return nil, errors.New("Host not found")
	}

	// Retrieve tulpa list (names)
	var hostval []string
	hostdata := database.Unjar("host." + hname)
	err := json.Unmarshal(hostdata, &hostval)
	if err != nil {
		return nil, err
	}

	return hostval, nil
}

func getTulpasByHost(hname string) ([]Tulpa, error) {

	hostval, err := getTulpaListByHost(hname)
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

	return tulpas, nil
}

func setNewTulpa(tulpa Tulpa) error {

	// Put tulpa record into database
	jsonTulpa, err := json.Marshal(tulpa)
	if err != nil {
		return err
	}
	database.Jar("tulpa."+tulpa.Host+"."+tulpa.Name, jsonTulpa)

	// Put tulpa into tulpalist
	tulpas = append(tulpas, tulpa.Host+"."+tulpa.Name)
	jsontulpa, err := json.Marshal(tulpas)
	if err != nil {
		return err
	}
	database.Jar("tulpalist", jsontulpa)

	// Put tulpa into host's list
	host := tulpa.Host
	hostlist := make([]string, 0)
	if database.Exists("host." + host) {
		hostlist, err = getTulpaListByHost(host)
		if err != nil {
			return err
		}
	}
	hostlist = append(hostlist, tulpa.Name)

	// Update host's list on the database
	jsonhost, err := json.Marshal(hostlist)
	if err != nil {
		return err
	}
	database.Jar("host."+host, jsonhost)

	return nil
}
