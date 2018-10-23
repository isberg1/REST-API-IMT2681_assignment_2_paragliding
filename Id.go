package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

// proccesses GET for "paragliding/api/Igc/ID"
func returnID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")

	igcStruct, ok := MgoTrackDB.getMetaByID(vars["ID"])
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
	err := json.NewEncoder(w).Encode(simpleIgcStruct)
	if err != nil {
		http.Error(w, "severside error(returnID)", http.StatusInternalServerError)
	}
}

// calculate the track_lenght in km
func calculateTrackLenght(track igc.Track) int {

	floatDistance := 0.0
	// sums up the total distance
	for i := 0; i < len(track.Points)-1; i++ {
		floatDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return int(floatDistance)
}

// converts URL string into i Meta struct
func parseFile(URLfile string) (Meta, error) {
	// extract IGC data form URL
	track, err := igc.ParseLocation(URLfile)
	if err != nil {
		return Meta{}, err
	}

	id, ok := getUniqueTrackID()
	if !ok {
		return Meta{}, errors.New("unable to get getUniqueTrackID")
	}
	// return a Meta struct with relevant data
	return Meta{
			ID:          id,
			TimeStamp:   getTimestamp(),
			URL:         URLfile,
			HDate:       track.Date.String(),
			Pilot:       track.Pilot,
			Glider:      track.GliderType,
			GliderID:    track.GliderID,
			TrackLength: calculateTrackLenght(track)},
		nil
}

// processes GET for url "paragliding/api/Igc/ID/feild"
func returnField(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")

	igcStruct, ok := MgoTrackDB.getMetaByID(vars["ID"])
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
