package main

import (
	"encoding/json"
	"fmt"
	"github.com/marni/goigc"
	"net/http"
)

// proccesses GET for "/igcinfo/api/Igc/ID"
func returnID(w http.ResponseWriter, r *http.Request, s string) {
	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")
	// checks if http method and id(s) is valid
	if !validateURL(w, r, s) {
		return
	}
	// return IGC meta in json format for ID 's'
	json.NewEncoder(w).Encode(IgcMap[s])
}

// converts URL string into i Meta struct
func parseFile(URL string) (Meta, error) {
	// extract IGC data form URL
	track, err := igc.ParseLocation(URL)
	if err != nil {
		return Meta{}, err
	}

	// return a Meta struct with relevant data
	temp := track.Task.Distance()
	distance := int(temp)
	return Meta{
			HDate:       track.Date.String(), // alternativ: "track.Header.Date.String()"
			Pilot:       track.Pilot,
			Glider:      track.GliderType,
			GliderID:    track.GliderID,
			TrackLength: distance},
		nil
}

// processes GET for url "/igcinfo/api/Igc/ID/feild"
func returnField(w http.ResponseWriter, r *http.Request, s []string) {
	// checks if http method and id(s) is valid
	if !validateURL(w, r, s[0]) {
		return
	}
	// convert relevant Meta struct to json
	jsonString, err := json.Marshal(IgcMap[s[0]])
	if err != nil {
		http.Error(w, "server Error", http.StatusInternalServerError)
		return
	}
	// convert json to map
	// (map set to string -> interface{} for easy refactoring purposes
	var temp map[string]interface{}
	json.Unmarshal(jsonString, &temp)

	// check if requested field from URL("S[1]") exist,
	// and if so return value to client
	if value, ok := temp[s[1]]; ok {
		http.Header.Add(w.Header(), "content-type", "text/plain")
		fmt.Fprintln(w, value)
		return
	}
	// return key not found to client
	http.NotFound(w, r)
}

// checks if method
func validateURL(w http.ResponseWriter, r *http.Request, s string) bool {
	//  wrong method at this URL
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return false
	}
	// if ID does not exist
	if _, ok := IgcMap[s]; !ok {
		http.Error(w, "", http.StatusNotFound)
		return false
	}
	// return true if everything is ok
	return true
}
