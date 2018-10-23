package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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
	respondToClient(w, trackStruct.ID)
	err4 := MgoTrackDB.add(trackStruct)
	if err4 != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	invokWebHooks(w)

}

// writes conformation message back to client verifying that POST was successful
func respondToClient(w http.ResponseWriter, s string) {
	//set http header content-type
	http.Header.Add(w.Header(), "content-type", "application/json")
	// set correct http status code
	w.WriteHeader(http.StatusCreated)
	// write conformation massage back to client
	response := ResponsID{ID: s} // empty struct needed to make empty json array
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "serverside error(respondToClient)", http.StatusInternalServerError)
	}
}

// makes a unique ID for Posted content to be stored in track collection
func getUniqueTrackID() (string, bool) {
	count := MgoTrackDB.count()
	if count == -1 {
		return "", false
	}
	if count == 0 {
		return "1", true
	}

	timeStamp, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		return "", false
	}

	meta, ok := MgoTrackDB.getWebHookByTimstamp(timeStamp)
	if !ok {
		return "", false
	}

	strID, err := strconv.Atoi(meta.ID)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	strID++
	id := strconv.Itoa(strID)
	return id, true
}

// makes a unique ID for Posted content to be stored in webhook collection
func getUniqueWebHookkID() (string, bool) {
	count := MgoWebHookDB.count()
	if count == -1 {
		return "", false
	}
	if count == 0 {
		return "1", true
	}

	timeStamp, ok := MgoWebHookDB.getLatestMetaTimestamp()
	if !ok {
		return "", false
	}

	webHook, ok := MgoWebHookDB.getWebHookByTimstamp(timeStamp)
	if !ok {
		return "", false
	}

	strID, err := strconv.Atoi(webHook.ID)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	strID++
	id := strconv.Itoa(strID)
	return id, true
}

// processes GET content for url "/paragliding/api/tracks/"
func getFiles(w http.ResponseWriter, r *http.Request) {
	//set http header content-type
	http.Header.Add(w.Header(), "content-type", "application/json")

	var keySlice []ResponsID
	// transfer all IgcMap key to its own slice"keySlice"
	// and put the keys into a slice "keySlice"

	ids, ok := MgoTrackDB.getAllKeys()
	if !ok {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	for _, val := range ids {
		ids = append(ids, val)
		temp := ResponsID{ID: val}
		keySlice = append(keySlice, temp)
	}
	// special case for no IGC file registered
	if MgoTrackDB.count() < 1 {
		// make an empty array
		keySlice = make([]ResponsID, 0)
		// write empty json array back to client
		err := json.NewEncoder(w).Encode(keySlice)
		if err != nil {
			http.Error(w, "serverside error(json.NewEncoder(w).Encode(keySlice))", http.StatusInternalServerError)
		}
		return
	}

	// general case
	w.WriteHeader(http.StatusOK)
	// write all keys for all registered IGC files back to client
	err := json.NewEncoder(w).Encode(keySlice)
	if err != nil {
		http.Error(w, "serverside error(json.NewEncoder(w).Encode(keySlice))", http.StatusInternalServerError)
		return
	}
}
