package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WebhookNewTrack(w http.ResponseWriter, r *http.Request) {

	/*
		What: Registration of new webhook for notifications about tracks being added to the system.
		Returns the details about the registration. The webhookURL is required parameter of the request.
		The minTriggerValue is optional integer, that defaults to 1 if ommited.
		It indicated the frequency of updates - after how many new tracks the webhook should be called.

	Response type: application/json
	Response code: 200 or 201 if everything is OK, appropriate error code otherwise.

	Request
	{
	    "webhookURL": {
	      "type": "string"
	    },
	    "minTriggerValue": {
	      "type": "number"
	    }
	}
	Example, that registers a webhook that should be trigger for every two new tracks added to the system.

	Response
	The response body should contain the id of the created resource (aka webhook registration), as string. Note,
		the response body will contain only the created id, as string, not the entire path; no json encoding.
		Response code upon success should be 200 or 201.
	*/

	var webHook WebHookStruct
	err := json.NewDecoder(r.Body).Decode(&webHook)
	if err != nil {
		http.Error(w, "bad request(json.NewDecoder(r.Body).Decode(&webHook))", http.StatusBadRequest)
		return
	}

	if webHook.MinTriggerValue == 0 {
		webHook.MinTriggerValue = 1
	}
	fmt.Fprintln(w, webHook)
	id, ok := getUniqueWebHookkID()
	if !ok {
		http.Error(w, "serverside error(getUniqueWebHookkID)", http.StatusInternalServerError)
		return
	}
	webHook.ID = id
	webHook.Counter = webHook.MinTriggerValue

	err1 := MgoWebHookDB.Add(webHook)
	if err1 != nil {
		http.Error(w, "serverside error(MgoWebHookDB.Add)", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, webHook.ID)
}

func InvokWebHooks() {

	MgoWebHookDB.counterDown()
	postToWebHooks()
	//MgoWebHookDB.counterReset()

}

func postToWebHooks() {

	var webHook []WebHookStruct

	webHook, err := MgoWebHookDB.GetPostArray()
	if err != nil {
		fmt.Println("unable to get post array", err)
	}

	for _, val := range webHook {
		err1 := postTo(val)
		if err1 != nil {
			fmt.Println("unable to post to ", err1)
		}
	}

}

func postTo(webHook WebHookStruct) error {
	fmt.Println(webHook)

	MgoTrackDB.G

	return nil
}
