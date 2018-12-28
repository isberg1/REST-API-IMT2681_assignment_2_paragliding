package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

/*

const localProjectURLroot = "http://localhost:8080/"
const localProjectURLbase = "http://localhost:8080/paragliding/"
const localProjectURLinfo1 = "http://localhost:8080/paragliding/api"
const localProjectURLinfo2 = "http://localhost:8080/paragliding/api/"
const localProjectURLarray1 = "http://localhost:8080/paragliding/api/track/"
const localProjectURLarray2 = "http://localhost:8080/paragliding/api/track"
const validIgcURL1 = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
const validIgcURL2 = "https://raw.githubusercontent.com/marni/goigc/master/testdata/optimize-long-flight-1.igc"

const contentType = "application/json"

*/

func setup(t *testing.T, method, url string, handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {

	// set up http test

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Error("unable to create http new request")
	}

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(url, handler)
	router.ServeHTTP(recorder, request)

	return recorder
}

// test the root of the URL: "/"
func Test_all(t *testing.T) {

	recorder := setup(t, "GET", "/", all)

	// checks the http status code
	expectedStatusCode := http.StatusNotFound
	if recorder.Code != expectedStatusCode {
		t.Error(
			"incorrect status code received, expected: " +
				string(expectedStatusCode) +
				" got: " +
				string(recorder.Code))
	}
}

// tests the redirected URL: "/paragliding"
func Test_redirect(t *testing.T) {

	recorder := setup(t, "GET", "/paragliding", rederect)

	// checks the http status code
	expectedStatusCode := http.StatusPermanentRedirect
	if recorder.Code != expectedStatusCode {
		t.Error(
			"incorrect status code received, expected: " +
				string(expectedStatusCode) +
				" got: " +
				string(recorder.Code))
	}
}

// Test_igcinfoapi check if the responding json is as expected
func Test_igcinfoapi(t *testing.T) {

	recorder := setup(t, "GET", "/paragliding/api", api)

	res := recorder.Body.Bytes()

	if res == nil {
		t.Error("error reading body")
	}

	var appInfo GetIgcinfoAPI
	// check if values are correct
	err3 := json.Unmarshal(res, &appInfo)
	if err3 != nil {
		t.Error("error umarshaling ", err3)
	}
	// check http status code
	expectedStatusCode := http.StatusOK
	if recorder.Code != expectedStatusCode {
		t.Error(
			"incorrect status code received, expected: " +
				string(expectedStatusCode) +
				" got: " +
				string(recorder.Code))
	}
	// check if values is correct
	expectedVersion := unavalabeVersinNr
	if appInfo.Version != expectedVersion {
		t.Error("error invalid version nr ")
	}
	// check if values is correct
	expectedInfo := infoSting
	if appInfo.Info != expectedInfo {
		t.Error("error invalid information string ")
	}
	// check if format is correct
	ok, err4 := regexp.Match(
		"P[0-9]*Y?[0-9]*M?[0-9]*W?[0-9]*D?T[0-9]*H?[0-9]*M?[0-9]*S",
		[]byte(appInfo.Uptime))
	if err4 != nil {
		t.Error("error unable to run regex check")
	}
	if !ok {
		t.Error("error incorrect uptime value/format", appInfo.Uptime)
	}
}
