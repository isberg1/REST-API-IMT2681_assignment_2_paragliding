package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func webhookNewTrack(w http.ResponseWriter, r *http.Request) {

	var webHook WebHookStruct
	err := json.NewDecoder(r.Body).Decode(&webHook)
	if err != nil {
		http.Error(w, "bad request(json.NewDecoder(r.Body).Decode(&webHook))", http.StatusBadRequest)
		return
	}

	if webHook.MinTriggerValue == 0 {
		webHook.MinTriggerValue = 1
	}

	id, ok := getUniqueWebHookkID()
	if !ok {
		http.Error(w, "serverside error(getUniqueWebHookkID)", http.StatusInternalServerError)
		return
	}
	webHook.ID = id
	webHook.Counter = webHook.MinTriggerValue
	webHook.TimeStamp = getTimestamp()

	err1 := MgoWebHookDB.add(webHook)
	if err1 != nil {
		http.Error(w, "serverside error(MgoWebHookDB.add)", http.StatusInternalServerError)
		return
	}

}

func invokWebHooks(w http.ResponseWriter) {

	processingStartTime := time.Now()

	MgoWebHookDB.counterDown()
	webHook := postToWebHooks(w, processingStartTime)
	MgoWebHookDB.counterReset(webHook)

}

func postToWebHooks(w http.ResponseWriter, processingStartTime time.Time) []WebHookStruct {

	var webHook []WebHookStruct

	webHook, err := MgoWebHookDB.getPostArray()
	if err != nil {
		fmt.Println("unable to get post array", err)
	}

	for _, val := range webHook {
		err1 := postTo(val, w, processingStartTime)
		if err1 != nil {
			fmt.Println("unable to post to ", err1)
		}
	}
	return webHook
}

func postTo(webHook WebHookStruct, w http.ResponseWriter, processingStartTime time.Time) error {
	/*
			{
		   "t_latest": <latest added timestamp of the entire collection>,
		   "tracks": [<id1>, <id2>, ...]
		   "processing": <time in ms of how long it took to process the request>
		}
	*/

	var ids []ResponsID
	ids, err := MgoTrackDB.getLatestMetaIDs(webHook.MinTriggerValue)
	if err != nil {
		fmt.Println("unable to get []ResponsID ", err)
		return err
	}

	latest, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		fmt.Println("unable to post get latest Track ")
		return errors.New("unable to get latest track")
	}
	temp := InvokeWebHookStruct{
		TLatest:    latest,
		Tracks:     ids,
		Processing: time.Since(processingStartTime).Nanoseconds() / int64(time.Millisecond),
	}
	a, err2 := json.Marshal(&temp)
	if err2 != nil {
		http.Error(w, "serverside error(json.Marshal(&temp))", http.StatusInternalServerError)
	} //fmt.Println(temp)

	// Todo post to right host (webHook.WebHookURL )
	_, err1 := http.Post("http://localhost:8080/test", "application/json", bytes.NewBuffer(a))
	if err1 != nil {
		return err1
	}

	return nil
}

func webhookID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	webHook, ok := MgoWebHookDB.getWebHookByID(vars["webhookID"])
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(SimpleWebHookStruct{
		WebHookURL:      webHook.WebHookURL,
		MinTriggerValue: webHook.MinTriggerValue,
	})
}

func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	webHook, err := MgoWebHookDB.deleteWebHook(vars["webhookID"])
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(SimpleWebHookStruct{
		WebHookURL:      webHook.WebHookURL,
		MinTriggerValue: webHook.MinTriggerValue,
	})
}
