package main

import (
	"fmt"
	"github.com/abbot/go-http-auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

// returns the size of the IGC Meta tracks database collection
func adminTracksCount(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	count := MgoTrackDB.count()
	http.Header.Add(w.Header(), "content-type", "text/plain")
	fmt.Fprint(w, count)
}

// deletes all IGC Meta tracks from database collection
func trackDropTable(w http.ResponseWriter, r *auth.AuthenticatedRequest /*r *http.Request*/) {

	http.Header.Add(w.Header(), "content-type", "text/plain")
	count := MgoTrackDB.count()
	if count < 1 {
		fmt.Fprint(w, count)
		return
	}

	err := MgoTrackDB.dropTable()
	if err != nil {
		http.Error(w, "serverside error, unable to drop collection", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, count)
}

func Secret(user, realm string) string {
	if user == "overlord" {
		//todo chang befor uploding
		// alternatively password could be read from (encrypted)file
		password := os.Getenv("ADMINPASSWORD") // "pass"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err == nil {
			return string(hashedPassword)
		}
	}
	return ""
}
