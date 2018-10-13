package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

//________________________________________________________________________

// application description string, used to inform client
const infoSting = "Service for IGC tracks."

// used to make uniqe id for IgcMap
const idPrefix = "IGC_file_"

// used as default value for version nr
const unavalabeVersinNr = "Unavalable"

// default port
const defaultPort = "8080"

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
	// handlerfuctions for general URL's
	http.HandleFunc("/", all)
	// the 2 almost identical handlerfunctions below are needed to  fix som buggy
	// behavior regarding http POST being permitted when it shouldn't.
	http.HandleFunc("/igcinfo/api", api)
	http.HandleFunc("/igcinfo/api/", api)

	http.HandleFunc("/igcinfo/api/igc", apiIgc)
	http.HandleFunc("/igcinfo/api/igc/", apiIgc)

	http.HandleFunc("/igcinfo/api/drop_table", dropTable)
	http.HandleFunc("/igcinfo/api/drop_table/", dropTable)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
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

// handels URL for "/igcinfo/api/drop_table/"
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
