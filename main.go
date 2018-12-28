package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/abbot/go-http-auth"

	"github.com/gorilla/mux"
)

//________________________________________________________________________

const (
	// application description string, used to inform client
	infoSting = "Service for IGC tracks."
	// used as default value for version nr
	unavalabeVersinNr = "Unavalable"
	// default port
	defaultPort = "8080"
	//default Paging number
	defaultPagingNr = "5"
)

//MgoTrackDB is the variable used to access the IGC Meta database collection
var MgoTrackDB = mongoDbStruct{}

//MgoWebHookDB is the variable used to access the webhook database collection
var MgoWebHookDB = mongoDbStruct{}

//StartUpTime registers the startup time for the application, used for calculating application runtime
var StartUpTime = time.Now()

//________________________________________________________________

func main() {

	// local DB
	MgoTrackDB.initTrackCollection("test", "mainCollection", "mongodb://127.0.0.1:27017")
	MgoWebHookDB.initWebHookCollection("test", "WebHook", "mongodb://127.0.0.1:27017")

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", all)
	//GET /paragliding
	r.HandleFunc("/{paragliding:paragliding[/]?}", rederect).Methods("GET")
	//GET /paragliding/api
	r.HandleFunc("/paragliding/{api:api[/]?}", api).Methods("GET")
	//GET /paragliding/api/track
	r.HandleFunc("/paragliding/api/{track:track[/]?}", getFiles).Methods("GET")
	//POST /paragliding/api/track
	r.HandleFunc("/paragliding/api/{track:track[/]?}", postFile).Methods("POST")
	//GET /paragliding/api/track/ID
	r.HandleFunc("/paragliding/api/track/{ID:[1-9]+}", returnID).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{ID:[1-9]+}/", returnID).Methods("GET")
	//GET /paragliding/api/track/ID/field
	r.HandleFunc("/paragliding/api/track/{ID:[1-9]+}/{field:h_date|pilot|glider|glider_id|track_length}", returnField).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{ID:[1-9]+}/{field:h_date|pilot|glider|glider_id|track_length}/", returnField).Methods("GET")
	//GET /paragliding/api/ticker/
	r.HandleFunc("/paragliding/api/{ticker:ticker[/]?}", apiTicker).Methods("GET")
	//GET /paragliding/api/ticker/latest
	r.HandleFunc("/paragliding/api/ticker/{latest:latest[/]?}", apiTtickerLatest).Methods("GET")
	//GET /paragliding/api/ticker/<timestamp>
	r.HandleFunc("/paragliding/api/ticker/{timestamp}", apiTimestamp).Methods("GET")
	//POST /paragliding/api/webhook/new_track/
	r.HandleFunc("/paragliding/api/webhook/{new_track:new_track[/]?}", webhookNewTrack).Methods("POST")
	//GET /api/webhook/new_track/<webhookID>
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhookID}{slash:[/]?}", webhookID).Methods("GET")
	//DELETE /api/webhook/new_track/<webhookID>
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhookID}{slash:[/]?}", deleteWebhook).Methods("DELETE")
	// set up basic http authentication
	authenticator := auth.NewBasicAuthenticator("calm-mesa-59678.herokuapp.com/", secret)
	//GET /admin/api/tracks_count
	r.HandleFunc("/admin/api/{track:tracks_count[/]?}", authenticator.Wrap(adminTracksCount)).Methods("GET")
	//DELETE /admin/api/tracks
	r.HandleFunc("/admin/api/{track:tracks[/]?}", authenticator.Wrap(trackDropTable)).Methods("DELETE")
	// for testing of webhook posts
	r.HandleFunc("/test", printRespons)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}

//________________________________________________________________________
// handels all URL for "/"
func all(w http.ResponseWriter, r *http.Request) {
	// URL not supported
	http.NotFound(w, r)
}

// redirects "/paragliding" to "/paragliding/api"
func rederect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/paragliding/api", http.StatusPermanentRedirect)
}

// stores webhook posts
var webhookStrings []string

// exists for testing of invoking webhook posts
func printRespons(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		a, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("printRespons", err)
		} else {
			webhookStrings = append(webhookStrings, string(a))
		}
	} else if r.Method == http.MethodGet {
		for _, val := range webhookStrings {
			// print webookposts
			fmt.Fprintln(w, val)
			fmt.Println(val)
		}
	}
}
