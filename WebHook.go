package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func WebhookNewTrack(w http.ResponseWriter, r *http.Request) {

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
	webHook.Id = id
	webHook.Counter = webHook.MinTriggerValue
	webHook.TimeStamp = getTimestamp()

	err1 := MgoWebHookDB.Add(webHook)
	if err1 != nil {
		http.Error(w, "serverside error(MgoWebHookDB.Add)", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, id)
}

func InvokWebHooks(w http.ResponseWriter) {

	processingStartTime := time.Now()

	MgoWebHookDB.counterDown()
	webHook := postToWebHooks(w, processingStartTime)
	MgoWebHookDB.counterReset(webHook)

}

func postToWebHooks(w http.ResponseWriter, processingStartTime time.Time) []WebHookStruct {

	var webHook []WebHookStruct

	webHook, err := MgoWebHookDB.GetPostArray()
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
	ids, err := MgoTrackDB.GetXLatest(webHook.MinTriggerValue)
	if err != nil {
		fmt.Println("unable to get []ResponsID ", err)
	}

	latest, ok := MgoTrackDB.GetLatest()
	if !ok {
		fmt.Println("unable to post get latest Track ")
	}
	temp := InvokeWebHookStruct{
		TLatest:    latest,
		Tracks:     ids,
		Processing: time.Since(processingStartTime).Nanoseconds() / int64(time.Millisecond),
	}
	a, _ := json.Marshal(&temp)
	//fmt.Println(temp)

	// Todo post to right host (webHook.WebHookURL )
	http.Post("http://localhost:8080/test", "application/json", bytes.NewBuffer(a))

	return nil
}

func Webhook_id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	webHook, ok := MgoWebHookDB.GetWebHook(vars["webhook_id"])
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

	webHook, err := MgoWebHookDB.DeleteWebHook(vars["webhook_id"])
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
