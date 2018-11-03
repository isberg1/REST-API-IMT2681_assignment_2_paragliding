package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// add a new webhook track to the DB
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
	// print the webhook id
	fmt.Fprint(w, webHook.ID)

}

// when new ICG tracks are posted  this function is called. it finds
// witch webhooks should be posted to based on a counter in the whbhook documents, in the DB
func invokWebHooks(w http.ResponseWriter) {

	processingStartTime := time.Now()
	// count down for all webhooks
	MgoWebHookDB.counterDown()
	// get and post to all the webhooks where the counter == 0
	webHook := postToWebHooks(w, processingStartTime)
	// reset the counter(back to minTriggerValue) for all webbhocks where counter == 0
	MgoWebHookDB.counterReset(webHook)

}

//  gets relevant array of webhooks and iterates over them in order to post to each one of them
func postToWebHooks(w http.ResponseWriter, processingStartTime time.Time) []WebHookStruct {

	var webHook []WebHookStruct

	webHook, err := MgoWebHookDB.getPostArray()
	if err != nil {
		fmt.Println("unable to get post array", err)
	}

	for _, val := range webHook {
		err1 := postTo(val, w, processingStartTime)
		if err1 != nil {
			fmt.Println("unable to post to ", val, err1)
		}
	}
	return webHook
}

// posts message to the URL stored in the webhook struckt
func postTo(webHook WebHookStruct, w http.ResponseWriter, processingStartTime time.Time) error {

	var ids []ResponsID
	ids, err := MgoTrackDB.getLatestMetaIDs(webHook.MinTriggerValue)
	if err != nil {
		fmt.Println("unable to get []ResponsID ", err)
		return err
	}

	latest, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		return errors.New("unable to get latest track")
	}

	//_________________________________________

	// old way of posting before specs where redefied

	temp := InvokeWebHookStruct{
		TLatest:    latest,
		Tracks:     ids,
		Processing: time.Since(processingStartTime).Nanoseconds() / int64(time.Millisecond),
	}
	a, err2 := json.Marshal(&temp)
	if err2 != nil {
		http.Error(w, "serverside error(json.Marshal(&temp))", http.StatusInternalServerError)
	}
	// had som weird problems when i wasn't ensuring content was a string
	str := fmt.Sprint(webHook.WebHookURL)
	str = strings.TrimSpace(str)

	_, err1 := http.Post(str, "application/json", bytes.NewBuffer(a))
	if err1 != nil {
		//fmt.Printf("%v %T", str, str)
		return err1
	}

	return nil
}

// prints a webhook struct as a json string
func webhookID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	webHook, ok := MgoWebHookDB.getWebHookByID(vars["webhookID"])
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(SimpleWebHookStruct{
		WebHookURL:      webHook.WebHookURL,
		MinTriggerValue: webHook.MinTriggerValue,
	})
}

// deletes a webhook from the collection based on a id found in the url
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
