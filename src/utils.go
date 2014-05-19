package main

func hasTulpa(hname, tname string) bool {
	str, _ := database.Get("tulpa." + hname + "." + tname)
	return str != ""
}

func getTulpa(hname string, tname string) (Tulpa, error) {
	var out Tulpa
	err := database.GetGob("tulpa."+hname+"."+tname, &out)
	return out, err
}

func getTulpaListByHost(hname string) ([]string, error) {

	keys, err := database.MatchPrefix("tulpa."+hname+".", 4096)
	if err != nil {
		return nil, err
	}
	return keys, nil
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

func setTulpa(tulpa Tulpa) error {
	return database.SetGob("tulpa."+tulpa.Host+"."+tulpa.Name, tulpa)
}

func getAllTulpas() ([]Tulpa, error) {
	keys, _ := database.MatchPrefix("tulpa.", 4096)

	tulpas := make([]Tulpa, len(keys))
	for i := range keys {
		err := database.GetGob(keys[i], &tulpas[i])
		if err != nil {
			return tulpas, err
		}
	}
	return tulpas, nil
}
