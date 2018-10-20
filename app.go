package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//________________________________________________________________________

// application description string, used to inform client
const infoSting = "Service for IGC tracks."

// used to make uniqe Id for IgcMap
const idPrefix = "IGC_file_"

// used as default value for version nr
const unavalabeVersinNr = "Unavalable"

// default port
const defaultPort = "8080"

//default Paging number
const defaultPagingNr = "5"

// MgnDB is the global
var MgoTrackDB = MongoDbStruct{}

//
var MgoWebHookDB = MongoDbStruct{}

// IgcMap global variable to store all IGC files
var IgcMap = make(map[string]Meta)

// global counter, used to make uniqe IGC file ID
var counter int

// StartUpTime registers the startup time for the application, used for calculating application runtime
var StartUpTime = time.Now()

// GlobalDebug used in debugging
var GlobalDebug = false

//________________________________________________________________________

func main() {

	MgoTrackDB.InitTrackCollection("test", "mainCollection", "mongodb://127.0.0.1:27017")
	MgoWebHookDB.InitWebHookCollection("test", "WebHook", "mongodb://127.0.0.1:27017")

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", all)

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

	r.HandleFunc("/test", printRespons).Methods("POST")
	/*
		//


		//GET /api/webhook/new_track/<webhook_id>
		r.HandleFunc("/api/webhook/new_track/{webhook_id}{slash:[/]?}", apiWebhookNew_trackWebhook_id).Methods("GET")
		//DELETE /api/webhook/new_track/<webhook_id>
		r.HandleFunc("/api/webhook/new_track/{webhook_id}{slash:[/]?}", delApiWebhookNew_trackWebhook_id).Methods("DELETE")

		//GET /admin/api/tracks_count
		r.HandleFunc("/admin/api/{track:tracks_count[/]?}", adminApiTracks_count).Methods("GET")
		//DELETE /admin/api/tracks
		r.HandleFunc("/admin/api/{track:tracks[/]?}", apiTicker).Methods("DELETE")
		/*

			http.HandleFunc("/paragliding/api/drop_table", dropTable)
			http.HandleFunc("/paragliding/api/drop_table/", dropTable)
	*/

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

func debug(w http.ResponseWriter, s string) {
	if GlobalDebug == true {
		fmt.Fprintln(w, "debug from "+s)
	}
}

// handels URL for "/paragliding/api/drop_table/"
func dropTable(w http.ResponseWriter, r *http.Request) {
	// process url
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	a := strings.Split(message, "/")

	//check if there are rubbis URL section after /drop_table/
	if len(a) > 3 && a[3] != "" {
		http.NotFound(w, r)
		return
	}
	// process http method
	// if DELETE method is used
	if r.Method == http.MethodDelete {
		if len(IgcMap) > 0 {
			IgcMap = make(map[string]Meta)
			counter = 0
		}
	} else {
		// if method is anything else
		http.Error(w, "illegal method", http.StatusMethodNotAllowed)
	}
}

func printRespons(w http.ResponseWriter, r *http.Request) {

	a, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("printRespons", err)
	}

	fmt.Println(string(a))

}
