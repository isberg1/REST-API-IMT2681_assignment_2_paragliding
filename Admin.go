package main

import (
	"fmt"
	"net/http"
)

//GET /admin/api/tracks_count
func adminTrackscount(w http.ResponseWriter, r *http.Request) {
	// TODO authentication

	count := MgoTrackDB.Count()

	http.Header.Add(w.Header(), "content-type", "text/plain")

	fmt.Fprintln(w, count)
}

func trackDropTable(w http.ResponseWriter, r *http.Request) {
	// TODO authentication
	http.Header.Add(w.Header(), "content-type", "text/plain")

	count := MgoTrackDB.Count()
	if count < 1 {
		fmt.Fprintln(w, count)
		return
	}

	err := MgoTrackDB.DropTable()
	if err != nil {
		http.Error(w, "serverside error, unable to drop collection", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, count)
}
