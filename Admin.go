package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abbot/go-http-auth"
	"golang.org/x/crypto/bcrypt"
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

// check to see if username and password is correct
func secret(user, realm string) string {
	if user == "overlord" {

		// alternatively password could be read from (encrypted)file
		password := os.Getenv("ADMINPASSWORD")
		if password == "" {
			password = "pass"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err == nil {
			return string(hashedPassword)
		}
	}
	return ""
}
