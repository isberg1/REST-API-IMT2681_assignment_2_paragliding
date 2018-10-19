package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

//processes http Methods for url "/paragliding/api/track/"
func selectMethodForAPIIgc(w http.ResponseWriter, r *http.Request) {
	// process method type and call appropriate response function
	switch r.Method {
	case http.MethodGet: // if GET method
		//getFiles(w)
	case http.MethodPost: // if POST method
		postFile(w, r)
	default: // for anything else, respond with error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// processes POST content for url "/paragliding/api/Igc/"
func postFile(w http.ResponseWriter, r *http.Request) {
	//read the POST body content
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// convert POST content from json and validate content
	var urlStruckt InputURL
	err2 := json.Unmarshal(content, &urlStruckt)
	if err2 != nil || urlStruckt.URL == "" {
		http.Error(w, "bad request "+strconv.Itoa(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//check that URL exits and that URL is valid IGC file
	trackStruct, err3 := parseFile(urlStruckt.URL)
	if err3 != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// put the IGC file string in global data structure,
	// and respond to client with conformation message

	//addToMap(trackStruct)
	respondToClient(w, trackStruct.Id)
	MgoDB.Add(trackStruct)

}

// writes conformation message back to client verifying that POST was successful
func respondToClient(w http.ResponseWriter, s string) {
	//set http header content-type
	http.Header.Add(w.Header(), "content-type", "application/json")
	// set correct http status code
	w.WriteHeader(http.StatusCreated)
	// write conformation massage back to client
	response := ResponsID{ID: s} // empty struct needed to make empty json array
	json.NewEncoder(w).Encode(response)
}

// makes a unique ID for Posted content to be stored in IgcMap
func getUniqueID() string {
	counter++ // increments global counter
	return strconv.Itoa(counter)
}

// processes GET content for url "/paragliding/api/Igc/"
func getFiles(w http.ResponseWriter, r *http.Request) {
	//set http header content-type
	http.Header.Add(w.Header(), "content-type", "application/json")

	var keySlice []ResponsID
	// transfer all IgcMap key to its own slice"keySlice"
	// and put the keys into a slice "keySlice"

	ids, ok := MgoDB.GetAllKeys()
	if !ok {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	for _, val := range ids {
		ids = append(ids, val)
		temp := ResponsID{ID: val}
		keySlice = append(keySlice, temp)
	}
	// special case for no IGC file registered
	if MgoDB.Count() < 1 {
		// make an empty array
		keySlice = make([]ResponsID, 0)
		// write empty json array back to client
		json.NewEncoder(w).Encode(keySlice)
		return
	}

	// general case
	w.WriteHeader(http.StatusOK)
	// write all keys for all registered IGC files back to client
	json.NewEncoder(w).Encode(keySlice)
}
