package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//________________________________________________________________________

const (
	// application description string, used to inform client
	infoSting = "Service for IGC tracks."
	// used to make uniqe Id for IgcMap
	idPrefix = "IGC_file_"
	// used as default value for version nr
	unavalabeVersinNr = "Unavalable"
	// default port
	defaultPort = "8080"
	//default Paging number
	defaultPagingNr = "5"
)

// MgnDB is the global
var MgoTrackDB = MongoDbStruct{}

//
var MgoWebHookDB = MongoDbStruct{}

// StartUpTime registers the startup time for the application, used for calculating application runtime
var StartUpTime = time.Now()

//________made basic admin handlers
//tested remote database ________________________________________________________________

func main() {

	MgoTrackDB.InitTrackCollection("test", "mainCollection", "mongodb://127.0.0.1:27017")
	MgoWebHookDB.InitWebHookCollection("test", "WebHook", "mongodb://127.0.0.1:27017")

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", all)

	r.HandleFunc("/{paragliding:paragliding[/]?}", rederect).Methods("GET")
	r.HandleFunc("/paragliding/{api:api[/]?}", api).Methods("GET")
	r.HandleFunc("/paragliding/api/{track:track[/]?}", getFiles).Methods("GET")
	r.HandleFunc("/paragliding/api/{track:track[/]?}", postFile).Methods("POST")
	r.HandleFunc("/paragliding/api/track/{Id:[1-9]+}", returnID).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{Id:[1-9]+}/", returnID).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{Id:[1-9]+}/{field:h_date|pilot|glider|glider_id|track_length}", returnField).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{Id:[1-9]+}/{field:h_date|pilot|glider|glider_id|track_length}/", returnField).Methods("GET")
	//GET /paragliding/api/ticker/
	r.HandleFunc("/paragliding/api/{ticker:ticker[/]?}", apiTicker).Methods("GET")
	//GET /paragliding/api/ticker/latest
	r.HandleFunc("/paragliding/api/ticker/{latest:latest[/]?}", apiTtickerLatest).Methods("GET")
	//GET /paragliding/api/ticker/<timestamp>
	r.HandleFunc("/paragliding/api/ticker/{timestamp}", apiTimestamp).Methods("GET")
	//POST /paragliding/api/webhook/new_track/
	r.HandleFunc("/paragliding/api/webhook/{new_track:new_track[/]?}", WebhookNewTrack).Methods("POST")
	//GET /api/webhook/new_track/<webhook_id>
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}{slash:[/]?}", Webhook_id).Methods("GET")
	//DELETE /api/webhook/new_track/<webhook_id>
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}{slash:[/]?}", deleteWebhook).Methods("DELETE")
	//GET /admin/api/tracks_count
	r.HandleFunc("/admin/api/{track:tracks_count[/]?}", adminTrackscount).Methods("GET")
	//DELETE /admin/api/tracks
	r.HandleFunc("/admin/api/{track:tracks[/]?}", trackDropTable).Methods("DELETE")

	r.HandleFunc("/test", printRespons).Methods("POST")

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

func rederect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/paragliding/api", http.StatusPermanentRedirect)
}

func printRespons(w http.ResponseWriter, r *http.Request) {
	a, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("printRespons", err)
	}
	fmt.Println(string(a))
}
