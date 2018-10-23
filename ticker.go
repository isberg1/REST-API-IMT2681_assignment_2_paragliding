package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

///api/ticker/latest
func apiTtickerLatest(w http.ResponseWriter, r *http.Request) {
	/*
			What: returns the timestamp of the latest added track
		Response type: text/plain
		Response code: 200 if everything is OK, appropriate error code otherwise.
		Response: <timestamp> for the latest added track
	*/
	w.Header().Add("Content-Type", "text/plain")

	timeStamp, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, timeStamp)

}

func apiTicker(w http.ResponseWriter, r *http.Request) {
	/*
		What: returns the JSON struct representing the ticker for the IGC tracks.
		The first track returned should be the oldest. The array of track ids returned should be capped at 5,
		to emulate "paging" of the responses. The cap (5) should be a configuration parameter of the
		application (ie. easy to change by the administrator).

		Response type: application/json
		Response code: 200 if everything is OK, appropriate error code otherwise.
		Response
		{
		"t_latest": <latest added timestamp>,
		"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
		"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
		"tracks": [<id1>, <id2>, ...]
		"processing": <time in ms of how long it took to process the request>
		}
	*/
	startTime := time.Now()
	w.Header().Add("Content-Type", "application/json")
	// if DB is empty
	if MgoTrackDB.count() == 0 {
		/* alternative sulution:
		ticker := Ticker{
			TLatest: -1,
			TStart:-1,
			TStorp:-1,
			Tracks:[]ResponsID{},
			Processing:time.Since(startTime).Nanoseconds() / int64(time.Millisecond),
		}
		json.NewEncoder(w).Encode(ticker)
		*/
		http.Error(w, "", http.StatusNoContent)
		return
	}

	latestimeStamp, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		http.NotFound(w, r)
		return
	}
	var nextTimeStamp int64
	nextTimeStamp, ok = MgoTrackDB.getOldestMetaByTimeStamp()
	if !ok {
		http.Error(w, "", http.StatusNoContent)
		return
	}
	// get the nr of  ID entries to be added in the response
	nr, err := getPagingNr()
	if err != nil {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	// if there is less document entries in the database then spesified from  getPagingNr(), then adjust "nr"
	if MgoTrackDB.count() < nr {
		nr = MgoTrackDB.count()
	}

	var idArray = make([]ResponsID, 0, 0)

	firsID, ok := MgoTrackDB.getMetaByTimstamp(nextTimeStamp)
	if !ok {
		http.Error(w, "", http.StatusNoContent)
		return
	}
	idArray = append(idArray, ResponsID{ID: firsID.ID})

	for i := 1; i < nr; i++ {
		temp, err1 := MgoTrackDB.getMetaBiggerThen(nextTimeStamp)
		if err1 != nil {
			http.Error(w, "serverside error(getMetaBiggerThen)", http.StatusInternalServerError)
			return
		}
		idArray = append(idArray, ResponsID{ID: temp.ID})
		nextTimeStamp = temp.TimeStamp
	}

	ticker := Ticker{
		TLatest:    latestimeStamp,
		TStart:     firsID.TimeStamp,
		TStorp:     nextTimeStamp,
		Tracks:     idArray,
		Processing: time.Since(startTime).Nanoseconds() / int64(time.Millisecond),
	}

	err2 := json.NewEncoder(w).Encode(ticker)
	if err2 != nil {
		http.Error(w, "serverside error(json.NewEncoder(w).Encode(ticker))", http.StatusNoContent)
		return
	}

}

func getPagingNr() (int, error) {
	nr := os.Getenv("PAGINGNR")

	if nr == "" {
		nr = defaultPagingNr
	}
	temp, err := strconv.Atoi(nr)

	if MgoTrackDB.count() < temp {
		temp = MgoTrackDB.count()
	}

	return temp, err
}

func apiTimestamp(w http.ResponseWriter, r *http.Request) {
	/*
			What: returns the JSON struct representing the ticker for the IGC tracks. The first returned track should have the timestamp HIGHER than the one provided in the query. The array of track IDs returned should be capped at 5, to emulate "paging" of the responses. The cap (5) should be a configuration parameter of the application (ie. easy to change by the administrator).
		Response type: application/json
		Response code: 200 if everything is OK, appropriate error code otherwise.
		Response:
		{
		   "t_latest": <latest added timestamp of the entire collection>,
		   "t_start": <the first timestamp of the added track>, this must be higher than the parameter provided in the query
		   "t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
		   "tracks": [<id1>, <id2>, ...]
		   "processing": <time in ms of how long it took to process the request>
		}
	*/

	startTime := time.Now()

	if MgoTrackDB.count() == 0 || MgoTrackDB.count() == 1 {
		/* alternative sulution:
		ticker := Ticker{
			TLatest: -1,
			TStart:-1,
			TStorp:-1,
			Tracks:[]ResponsID{},
			Processing:time.Since(startTime).Nanoseconds() / int64(time.Millisecond),
		}
		json.NewEncoder(w).Encode(ticker)
		*/
		http.Error(w, "", http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)

	// set response type for http header
	http.Header.Add(w.Header(), "content-type", "application/json")
	// extract timestamp from URL
	nextTimeStamp, err := strconv.ParseInt(vars["timestamp"], 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// if timestamp exits
	_, ok := MgoTrackDB.getMetaByTimstamp(nextTimeStamp)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// get the nr of  ID entries to be added in the response
	nr, err := getPagingNr()
	if err != nil {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	// if there is less document entries in the database then spesified from  getPagingNr(), then adjust "nr"
	if MgoTrackDB.count() < nr {
		nr = MgoTrackDB.count()
	}

	var idArray = make([]ResponsID, 0, 0)

	firsID, err1 := MgoTrackDB.getMetaBiggerThen(nextTimeStamp)
	if err1 != nil {
		http.Error(w, "serverside error(getMetaBiggerThen)", http.StatusNoContent)
		return
	}
	idArray = append(idArray, ResponsID{firsID.ID})
	nextTimeStamp = firsID.TimeStamp

	for i := 0; i < nr-1; i++ {
		temp, err1 := MgoTrackDB.getMetaBiggerThen(nextTimeStamp)
		if err1 != nil {
			continue
		}
		idArray = append(idArray, ResponsID{ID: temp.ID})
		nextTimeStamp = temp.TimeStamp
	}

	latestimeStamp, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		http.NotFound(w, r)
		return
	}

	ticker := Ticker{
		TLatest:    latestimeStamp,
		TStart:     firsID.TimeStamp,
		TStorp:     nextTimeStamp,
		Tracks:     idArray,
		Processing: time.Since(startTime).Nanoseconds() / int64(time.Millisecond),
	}

	err2 := json.NewEncoder(w).Encode(ticker)
	if err2 != nil {
		http.Error(w, "serverside error(json.NewEncoder(w).Encode(ticker))", http.StatusNoContent)
		return
	}
}
