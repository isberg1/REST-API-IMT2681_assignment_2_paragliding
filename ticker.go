package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
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

	timeStamp, ok := MgoTrackDB.GetLatest()
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, timeStamp)

}

func apiTicker(w http.ResponseWriter, r *http.Request) {
	/*
		What: returns the JSON struct representing the ticker for the IGC tracks. The first track returned should be the oldest. The array of track ids returned should be capped at 5, to emulate "paging" of the responses. The cap (5) should be a configuration parameter of the application (ie. easy to change by the administrator).
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

	if MgoTrackDB.Count() == 0 {
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

	latestimeStamp, ok := MgoTrackDB.GetLatest()
	if !ok {
		http.NotFound(w, r)
		return
	}
	var nextTimeStamp int64
	nextTimeStamp, ok = MgoTrackDB.GetOldest()
	if !ok {
		http.Error(w, "", http.StatusNoContent)
		return
	}
	nr, err := getPagingNr()
	if err != nil {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	if MgoTrackDB.Count() < nr {
		nr = MgoTrackDB.Count()
	}

	var idArray = make([]ResponsID, 0, 0)

	firsID, ok := MgoTrackDB.GetByTimstamp(nextTimeStamp)
	if !ok {
		http.Error(w, "", http.StatusNoContent)
		return
	}
	idArray = append(idArray, ResponsID{ID: firsID.Id})

	for i := 1; i < nr; i++ {
		temp, err1 := MgoTrackDB.GetBiggerThen(nextTimeStamp)
		if err1 != nil {
			http.Error(w, "serverside error(GetBiggerThen)", http.StatusInternalServerError)
			return
		}
		idArray = append(idArray, ResponsID{ID: temp.Id})
		nextTimeStamp = temp.TimeStamp
	}

	ticker := Ticker{
		TLatest:    latestimeStamp,
		TStart:     firsID.TimeStamp,
		TStorp:     nextTimeStamp,
		Tracks:     idArray,
		Processing: time.Since(startTime).Nanoseconds() / int64(time.Millisecond),
	}

	json.NewEncoder(w).Encode(ticker)

}

func getPagingNr() (int, error) {
	nr := os.Getenv("PAGINGNR")

	if nr == "" {
		nr = defaultPagingNr
	}
	temp, err := strconv.Atoi(nr)

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

	if MgoTrackDB.Count() == 0 || MgoTrackDB.Count() == 1 {
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
	nextTimeStamp, err := strconv.ParseInt(vars["timestamp"], 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	_, ok := MgoTrackDB.GetByTimstamp(nextTimeStamp)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	nr, err := getPagingNr()
	if err != nil {
		http.Error(w, "serverside error", http.StatusInternalServerError)
	}
	if MgoTrackDB.Count() < nr {
		nr = MgoTrackDB.Count()
	}

	var idArray = make([]ResponsID, 0, 0)
	firsID, err1 := MgoTrackDB.GetBiggerThen(nextTimeStamp)
	if err1 != nil {
		http.Error(w, "serverside error(GetBiggerThen111)", http.StatusNoContent)
		return
	}
	idArray = append(idArray, ResponsID{firsID.Id})
	nextTimeStamp = firsID.TimeStamp

	for i := 0; i < nr-1; i++ {
		temp, err1 := MgoTrackDB.GetBiggerThen(nextTimeStamp)
		if err1 != nil {
			continue
		}
		idArray = append(idArray, ResponsID{ID: temp.Id})
		nextTimeStamp = temp.TimeStamp
	}

	latestimeStamp, ok := MgoTrackDB.GetLatest()
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

	json.NewEncoder(w).Encode(ticker)

}

//todo secure id to new files on restart gets max value from DB
