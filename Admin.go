package main

import (
	"fmt"
	"net/http"
)

// returns the size of the IGC Meta tracks database collection
func adminTracksCount(w http.ResponseWriter, r *http.Request) {
	// TODO authentication

	count := MgoTrackDB.count()

	http.Header.Add(w.Header(), "content-type", "text/plain")

	fmt.Fprintln(w, count)
}

// deletes all IGC Meta tracks from database collection
func trackDropTable(w http.ResponseWriter, r *http.Request) {
	// TODO authentication
	http.Header.Add(w.Header(), "content-type", "text/plain")

	count := MgoTrackDB.count()
	if count < 1 {
		fmt.Fprintln(w, count)
		return
	}

	err := MgoTrackDB.dropTable()
	if err != nil {
		http.Error(w, "serverside error, unable to drop collection", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, count)
}
