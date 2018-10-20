package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"net/http"
	"sync"
	"time"
)

// proccesses GET for "paragliding/api/Igc/ID"
func returnID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")

	igcStruct, ok := MgoTrackDB.Get(vars["Id"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	simpleIgcStruct := SimpleMeta{
		HDate:       igcStruct.HDate,
		Pilot:       igcStruct.Pilot,
		Glider:      igcStruct.Glider,
		GliderID:    igcStruct.GliderID,
		TrackLength: igcStruct.TrackLength,
	}

	// return IGC meta in json format for ID 's'
	json.NewEncoder(w).Encode(simpleIgcStruct)
}

// converts URL string into i Meta struct
func parseFile(URLfile string) (Meta, error) {
	// extract IGC data form URL
	track, err := igc.ParseLocation(URLfile)
	if err != nil {
		return Meta{}, err
	}
	// return a Meta struct with relevant data
	temp := track.Task.Distance()
	distance := int(temp)

	id, ok := getUniqueTrackID()
	if !ok {
		return Meta{}, errors.New("unable to get getUniqueTrackID")
	}
	return Meta{
			Id:          id,
			TimeStamp:   getTimestamp(),
			URL:         URLfile,
			HDate:       track.Date.String(), // alternativ: "track.Header.Date.String()"
			Pilot:       track.Pilot,
			Glider:      track.GliderType,
			GliderID:    track.GliderID,
			TrackLength: distance},
		nil
}

// processes GET for url "paragliding/api/Igc/ID/feild"
func returnField(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")

	igcStruct, ok := MgoTrackDB.Get(vars["Id"])
	if !ok {
		fmt.Println("from not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch vars["field"] {
	case "pilot":
		fmt.Fprintln(w, igcStruct.Pilot)
	case "h_date":
		fmt.Fprintln(w, igcStruct.HDate)
	case "glider":
		fmt.Fprintln(w, igcStruct.Glider)
	case "glider_id":
		fmt.Fprintln(w, igcStruct.GliderID)
	case "track_length":
		fmt.Fprintln(w, igcStruct.TrackLength)
	}

}

func getTimestamp() int64 {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	timeStamp := time.Now().UnixNano() / 1000000
	mutex.Unlock()

	return timeStamp
}
