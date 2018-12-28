package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const localProjectURLroot = "http://localhost:8080/"
const localProjectURLbase = "http://localhost:8080/paragliding/"
const localProjectURLinfo1 = "http://localhost:8080/paragliding/api"
const localProjectURLinfo2 = "http://localhost:8080/paragliding/api/"
const localProjectURLarray1 = "http://localhost:8080/paragliding/api/track/"
const localProjectURLarray2 = "http://localhost:8080/paragliding/api/track"
const validIgcURL1 = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
const validIgcURL2 = "https://raw.githubusercontent.com/marni/goigc/master/testdata/optimize-long-flight-1.igc"

const contentType = "application/json"

// start the server so other tests can be run later
func Test_startServer(t *testing.T) {
	// start local server
	go main()
	time.Sleep(2 * time.Second)
	err := MgoTrackDB.dropTable()
	if err != nil {
		t.Error("unable to drop collection", err)
	}
	err1 := MgoWebHookDB.dropTable()
	if err1 != nil {
		t.Error("unable to drop collection", err)
	}
}

// wait for server to be established
func Test_httpConnection(t *testing.T) {
	for i := 0; i < 5; i++ {
		_, err := http.Get(localProjectURLroot)
		// if connection was successful
		if err == nil {
			i = 5 // abort loop
		} else {
			time.Sleep(2 * time.Second)
		}
		// stop trying to connect
		if i == 4 {
			t.Error("error unable to connect to  ", localProjectURLroot, err)
		}
	}
}

// checks if correct http status code 404 is received from expected URL's
func Test_rubbishURL_local(t *testing.T) {
	expected := http.StatusNotFound

	temp := make([]string, 0)
	temp = append(
		temp,
		localProjectURLroot,
		localProjectURLbase+"rubbish",
		localProjectURLinfo1+"rubbish",
		localProjectURLarray1+"rubbish",
	)
	// iterates over URL's and checks responds value
	for _, res := range temp {
		get, err := http.Get(res)
		if err != nil {
			t.Error("error getting content from URL", res, err)
		}
		defer get.Body.Close()
		if get.StatusCode != expected {
			t.Error("error incorrect status code from: ", res, err)
		}
	}
}

// tries to post at an invalid URL
func Test_PostAtInvalidURL(t *testing.T) {
	expected := make([]int, 0)

	expected = append(expected,
		http.StatusMethodNotAllowed,
		http.StatusMethodNotAllowed,
		http.StatusNotFound,
		http.StatusNotFound,
	)

	temp := make([]string, 0)
	temp = append(
		temp,
		localProjectURLinfo1,
		localProjectURLinfo2,
		localProjectURLarray1+"rubbish",
		localProjectURLarray1+"rubbish/pilot",
	)

	jsonVar, err := json.Marshal(InputURL{URL: validIgcURL1})
	if err != nil {
		t.Error("error marshaling into json")
	}
	// iterate over slice with URL's
	for nr, res := range temp {
		post, err := http.Post(res, contentType, bytes.NewBuffer(jsonVar))
		if err != nil {
			t.Error("error unable to POST for:", res)
		}
		defer post.Body.Close()
		if post.StatusCode != expected[nr] {
			t.Error("error illegal POST permitted form: ", res, post.StatusCode)
		}
	}
}

// tries to post invalid content to correct URL
func Test_PostInvalidContent1(t *testing.T) {
	expected := http.StatusBadRequest

	invalidIgcFile := make([]string, 0)
	invalidIgcFile = append(
		invalidIgcFile,
		"https://github.com/",
		"https://raw.githubusercontent.com/marni/goigc/master/testdata/parse-0-invalid-record.0.igc",
		"https://raw.githubusercontent.com/marni/goigc/master/testdata/parse-c-invalid-finish.0.igc",
	)
	// iterate over slice with IGC file's
	for _, res := range invalidIgcFile {
		jsonVar, err := json.Marshal(InputURL{URL: res})
		if err != nil {
			t.Error("error marshaling into json")
		}
		post, err2 := http.Post(localProjectURLarray1, contentType, bytes.NewBuffer(jsonVar))
		if err2 != nil {
			t.Error("error unsuccessful post attempt for ", res)
		}
		defer post.Body.Close()
		if post.StatusCode != expected {
			t.Error("error illegal POST permitted for: ", res)
		}
	} // test if IgcMap length is as expected
	if MgoTrackDB.count() != 0 {
		t.Error("error invalid content posted in data structure: IgcMap")
	}
}

// tries to post valid content at valid URL
func Test_PostValidContent(t *testing.T) {
	expected := http.StatusCreated
	// IGC files URL's to be posted
	igcULR := make([]string, 0)
	igcULR = append(
		igcULR,
		validIgcURL1,
		validIgcURL2,
	)
	// URL to send POST content to
	postURL := make([]string, 0)
	postURL = append(
		postURL,
		localProjectURLarray1,
		localProjectURLarray2,
	)
	// for all IGC files
	for _, res := range igcULR {
		// make json object to be posted
		jsonVar, err := json.Marshal(InputURL{URL: res})
		if err != nil {
			t.Error("error marshaling into json")
		} // for all  URL to send POST content to
		for _, pURL := range postURL {
			// try posing
			post, err2 := http.Post(pURL, contentType, bytes.NewBuffer(jsonVar))
			if err2 != nil {
				t.Error("error unable to POST from: ", res)
			}
			defer post.Body.Close()
			if post.StatusCode != expected {
				t.Error("error legal post not permitted for: ", res)
			}
			// more validation tests
			err3 := validatePostResponse(post)
			if err3 != nil {
				t.Error("error ", err3, " from: ", res)
			}
		}
	}
	// check that number of entries in IgcMap is correct
	if MgoTrackDB.count() != (len(postURL) + len(igcULR)) {
		t.Error("error data structure does not contain expected nr of values: ", MgoTrackDB.count())
	}
}

// checks that ioutil can read content, content can be unmarshaled and that content string is correct via regex
func validatePostResponse(p *http.Response) error {
	read, err := ioutil.ReadAll(p.Body)
	if err != nil {
		return errors.New("can not read content body ")
	}

	strc := ResponsID{}
	err2 := json.Unmarshal(read, &strc)
	if err2 != nil {
		return errors.New("can not unmarshal ")
	}

	test, err4 := regexp.MatchString("[1-9]*", strc.ID)
	if err4 != nil {
		return errors.New("defective regex compilation ")
	}
	if !test {
		return errors.New("incorrect return ID string based on regex match ")
	}
	return nil
}

// tries to get the json array of all stored Igc file ID's
func Test_getAllIDs(t *testing.T) {
	expected := http.StatusOK

	get, err := http.Get(localProjectURLarray1)
	if err != nil {
		t.Error("error getting content from URL", localProjectURLarray1)
	}
	defer get.Body.Close()
	if get.StatusCode != expected {
		t.Error("error incorrect status code", get.StatusCode)
	}
	res, err2 := ioutil.ReadAll(get.Body)
	if err2 != nil {
		t.Error("error reading body", err2)
	}
	slice := make([]ResponsID, 0)
	err3 := json.Unmarshal(res, &slice)
	if err3 != nil {
		t.Error("error unable to unmarshal json array", err3)
	}
	if len(slice) != MgoTrackDB.count() {
		t.Error("error not the same nr of objects in local and global data structures")
	}
}

func Test_returnID(t *testing.T) {
	expected := http.StatusOK

	if MgoTrackDB.count() < 1 {
		t.Error("error no tracks in DB")
	}

	get, err := http.Get(localProjectURLarray1 + "/" + fmt.Sprintf("%v", 1))
	if err != nil {
		t.Error("error getting content from URL", localProjectURLarray1+"/"+fmt.Sprintf("%v", 1))
	}
	defer get.Body.Close()
	if get.StatusCode != expected {
		t.Error("error incorrect status code", get.StatusCode)
	}

	simpleMetaStruct := SimpleMeta{}

	err2 := json.NewDecoder(get.Body).Decode(&simpleMetaStruct)
	if err2 != nil {
		t.Error("error reading body", err2)
	}

	_, err3 := http.Get(simpleMetaStruct.URL)
	if err3 != nil {
		t.Error("error unable to get content form URL")
	}
}

// tries to get single fields form URL's
func Test_getFields(t *testing.T) {
	expected := http.StatusOK

	metaKey := make([]string, 0)
	metaKey = append(
		metaKey,
		"h_date",
		"pilot",
		"glider",
		"glider_id",
		"track_length",
	)

	if MgoTrackDB.count() < 1 {
		t.Error("error no entries in DB")
	}

	for _, field := range metaKey {

		strKey := strconv.Itoa(1)

		myURL := localProjectURLarray1 + strKey + "/" + field
		get, err := http.Get(myURL)
		if err != nil {
			t.Error("error getting from URL", err)
		}
		defer get.Body.Close()
		if get.StatusCode != expected {
			t.Error("error invalid status code from", myURL)
		}
		res, err2 := ioutil.ReadAll(get.Body)
		if err2 != nil {
			t.Error("error reading body", err2)
		}
		str := string(res)
		if str == "" {
			t.Error("error illegal content has been successfully posted")
		}
		// if the field is "track_length"
		if field == metaKey[4] {
			// try to convert the field "track_length" into int
			str = strings.TrimSpace(str)
			_, err3 := strconv.Atoi(str)
			if err3 != nil {
				t.Error("error illegal value(not int) for", err3, field, str)
			}
		}
	}
}

// tests the "http://localhost:8080/paragliding/api/ticker/latest" handler and
// checks the response
func Test_apiTtickerLatest(t *testing.T) {

	expectedContentType := "text/plain"
	expectedStatusCode := 200
	biggerThenExpected := getTimestamp()

	get, err := http.Get(localProjectURLroot + "paragliding/api/ticker/latest")
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error("error invalid status code")
	}
	if get.Header.Get("Content-Type") != expectedContentType {
		t.Error("error invalid Content-Type ")
	}

	res, err2 := ioutil.ReadAll(get.Body)
	if err2 != nil {
		t.Error("error reading body", err2)
	}

	intRes, err3 := strconv.ParseInt(string(res), 10, 64)
	if err3 != nil {
		t.Error("error convering to int", err3)
	}
	if biggerThenExpected <= intRes {
		t.Error("error reading body", biggerThenExpected, intRes)
	}
}

// tests the "http://localhost:8080/paragliding/api/ticker" handler and
// checks the respons
func Test_apiTicker(t *testing.T) {

	expectedContentType := "application/json"
	expectedStatusCode := 200

	get, err := http.Get(localProjectURLroot + "paragliding/api/ticker/")
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error("error invalid status code")
	}
	if get.Header.Get("Content-Type") != expectedContentType {
		t.Error("error invalid Content-Type ", get.Header.Get("Content-Type"))
	}

	res, err2 := ioutil.ReadAll(get.Body)
	if err2 != nil {
		t.Error("error reading body", err2)
	}
	var ticker Ticker
	err3 := json.Unmarshal(res, &ticker)
	if err3 != nil {
		t.Error("error unable to unmarshal json array", err3)
	}
	timeStamp, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		t.Error("error unable to get latest timestamp")
	}

	if ticker.TLatest != timeStamp {
		t.Error("error wrong timestamp")
	}
	pagingNr, err4 := getPagingNr()
	if err4 != nil {
		t.Error("error unable to get pagingNr")
	}
	if len(ticker.Tracks) != pagingNr {
		t.Error("error wrong len of tincker.Tracks pagingNr", len(ticker.Tracks), pagingNr)
	}
}

func Test_apiTimestamp(t *testing.T) {

	expectedContentType := "application/json"
	expectedStatusCode := 200

	time, ok := MgoTrackDB.getOldestMetaByTimeStamp()
	if !ok {
		fmt.Println(time)
		fmt.Println(fmt.Sprintf("%v", time))

		t.Error("error no entries in DB")
	}

	get, err := http.Get(localProjectURLroot + "paragliding/api/ticker/" + fmt.Sprintf("%v", time))
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	// check if response status code is correct
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error(
			"incorrect status code received, expected: " +
				string(expectedStatusCode) +
				" got: " +
				string(get.StatusCode))
	}
	// check if response has correct header type
	if get.Header.Get("Content-Type") != expectedContentType {
		t.Error("error invalid Content-Type ", get.Header.Get("Content-Type"))
	}
	res, err2 := ioutil.ReadAll(get.Body)
	if err2 != nil {
		t.Error("error reading body", err2)
	}
	// check if response is in correct format
	var ticker Ticker
	err3 := json.Unmarshal(res, &ticker)
	if err3 != nil {
		t.Error("error unable to unmarshal json array", err3)
	}

	latestTime, ok := MgoTrackDB.getLatestMetaTimestamp()
	if !ok {
		t.Error("error unable to get latest timestamp")
	}
	// check if received value is correct
	if latestTime != ticker.TLatest {
		t.Error(
			"error wrong timestamp received, expected: " +
				fmt.Sprintf("%v", latestTime) +
				" got " +
				fmt.Sprintf("%v", ticker.TLatest))
	}

}

// tests the "http://localhost:8080/paragliding/api/webhook/new_track" handler and
// checks the respons
func Test_WebhookNewTrack(t *testing.T) {

	expectedStatusCode := 200
	//dbCountPriorToPost := MgoWebHookDB.count()

	subscriberWebHook := localProjectURLroot + "/test"
	postURL := localProjectURLroot + "paragliding/api/webhook/new_track"

	jsonString, err := json.Marshal(SimpleWebHookStruct{WebHookURL: subscriberWebHook, MinTriggerValue: 3})
	if err != nil {
		t.Error("error marshaling into json")
	} // for all  URL to send POST content to

	for i := 0; i < 2; i++ {

		// try posing
		post, err2 := http.Post(postURL, contentType, bytes.NewBuffer(jsonString))
		if err2 != nil {
			t.Error("error unable to POST from: ", postURL)
		}
		defer post.Body.Close()
		if post.StatusCode != expectedStatusCode {
			t.Error("error legal post not permitted for: ", postURL)
		}

		id, err3 := ioutil.ReadAll(post.Body)
		if err3 != nil {
			t.Error("error unable to read Post.Body : ", err3)
		}

		responce, ok := MgoWebHookDB.getWebHookByID(string(id))
		if !ok {
			t.Error("error was unsuccessful at posting to webHook Subscription", responce, "--", string(id))
		}

	}

}

// /paragliding/api/webhook/new_track/{webhookID}{slash:[/]?}", webhookID
func Test_getWebhookByID(t *testing.T) {
	expectedContentType := "application/json"
	expectedStatusCode := 200

	get, err := http.Get(localProjectURLroot + "paragliding/api/webhook/new_track/1")
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error("error invalid status code")
	}
	if get.Header.Get("Content-Type") != expectedContentType {
		t.Error("error invalid Content-Type ", get.Header.Get("Content-Type"))
	}
}

// tests the handler authenticator.Wrap(adminTracksCount)/admin/api/{track:tracks_count[/]?}
func Test_adminCount(t *testing.T) {
	expectedStatusCode := http.StatusUnauthorized

	get, err := http.Get(localProjectURLroot + "admin/api/tracks_count")
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error("error invalid status code", get.StatusCode)
	}

}

// tests the handeler authenticator.authenticator.Wrap(trackDropTable) /admin/api/tracks
func Test_adminTrackDropTable(t *testing.T) {
	expectedStatusCode := http.StatusMethodNotAllowed

	var dummy io.Reader
	get, err := http.Post(localProjectURLroot+"admin/api/tracks", "", dummy)
	if err != nil {
		t.Error("error getting from URL", err)
	}
	defer get.Body.Close()
	if get.StatusCode != expectedStatusCode {
		fmt.Println(get.StatusCode)
		t.Error("error invalid status code", get.StatusCode)
	}

}

func Test_adminTracksCount(t *testing.T) {

	username := "overlord"
	passwd := "pass"
	myURL := localProjectURLroot + "admin/api/tracks_count"
	// set up http test

	client := &http.Client{}
	req, err := http.NewRequest("GET", myURL, nil)
	if err != nil {
		t.Error("error unable to make new http request")
	}
	req.SetBasicAuth(username, passwd)
	resp, err2 := client.Do(req)
	if err2 != nil {
		t.Error("error unable to run client.Doo")
	}
	bodyText, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		t.Error("error unable to read response body")
	}
	str := string(bodyText)

	if str == "" {
		t.Error("error no content in response body")
	}

	_, err3 := strconv.Atoi(str)
	if err3 != nil {
		t.Error("error response not an int", str)
	}
}

func Test_trackDropTable(t *testing.T) {

	username := "overlord"
	passwd := "pass"
	myURL := localProjectURLroot + "admin/api/tracks"

	// set up http test

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", myURL, nil)
	if err != nil {
		t.Error("error unable to make new http request")
	}
	req.SetBasicAuth(username, passwd)
	resp, err2 := client.Do(req)
	if err2 != nil {
		t.Error("error unable to run client.Doo")
	}
	bodyNumber, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		t.Error("error unable to read response body")
	}
	nr, err3 := strconv.Atoi(string(bodyNumber))
	if err3 != nil {
		t.Error("error received item not a nr ", nr)
	}

	if MgoTrackDB.count() != 0 {
		t.Error("error DB tracks not dropped")
	}

}

// drops the collection
func Test_cleanUp(t *testing.T) {
	err := MgoTrackDB.dropTable()
	if err != nil {
		t.Error("unable to drop collection", err)
	}
	err1 := MgoWebHookDB.dropTable()
	if err1 != nil {
		t.Error("unable to drop collection", err)
	}
}
