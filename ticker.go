package main

import (
	"net/http"
)

///api/ticker/latest
func apiTtickerLatest(w http.ResponseWriter, r *http.Request) {
	/*
		What: returns the timestamp of the latest added track
	Response type: text/plain
	Response code: 200 if everything is OK, appropriate error code otherwise.
	Response: <timestamp> for the latest added track
	*/

}
