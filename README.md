 README #

### What is this repository for? ###
Assignment 2 in IMT2681-2018

Assignment URL: https://calm-mesa-59678.herokuapp.com/paragliding

### How do I test the remote Heroku api ###

From bash terminal use the following commads:

GET /paragliding/

    curl --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/

GET /api

    curl --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api

POST /api/track

    curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"url": "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}' -X POST https://calm-mesa-59678.herokuapp.com/paragliding/api/track

GET /api/track

    curl --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track

GET /api/track/<id>

    curl --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track/1

GET /api/track/<id>/<field>

    curl --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track/1/pilot

GET /api/ticker/latest

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/latest

GET /api/ticker/

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/

GET /api/ticker/<timestamp>

    !!! write timestamp at the end of command !!!
    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/

POST /api/webhook/new_track/

// TODO make url handler for diplaying webhook content when min_trigger_value is invoced

     curl  --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"web_hook_url": "http://raw.githubusercontent.com/marni/goigc/", "min_trigger_value" : 5 }' -X POST https://calm-mesa-59678.herokuapp.com/paragliding/api/webhook/new_track

GET /api/webhook/new_track/<webhook_id>

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/webhook/new_track/1

DELETE /api/webhook/new_track/<webhook_id>

    curl  --write-out "\n%{http_code} %{content_type}\n"  -X DELETE https://calm-mesa-59678.herokuapp.com/paragliding/api/webhook/new_track/1

GET /admin/api/tracks_count

no auth:

    curl  --write-out "\n%{http_code} %{content_type}\n"  -X GET https://calm-mesa-59678.herokuapp.com/admin/api/tracks_count
with auth:

    curl  -u overlord:pass --write-out "\n%{http_code} %{content_type}\n"  -X GET https://calm-mesa-59678.herokuapp.com/admin/api/tracks_count

DELETE /admin/api/tracks

no auth:

    curl  --write-out "\n%{http_code} %{content_type}\n"  -X DELETE https://calm-mesa-59678.herokuapp.com/admin/api/tracks
with auth:

    curl  -u overlord:pass --write-out "\n%{http_code} %{content_type}\n"  -X DELETE https://calm-mesa-59678.herokuapp.com/admin/api/tracks


### How do I use the program? ###


In my go setup i was using go version: go1.11
I nedded to install "gcc" to run the tests

	sudo apt install gcc
	cd ~/
	git clone git@bitbucket.org:isberg/igcinfo.git
	cd igcinfo/go-getting-started/
	go run .




### code quality ###

code quality checking:

Static code analysis:  

      go tool vet -all .
      golint .
      go fmt .
      gometalinter -- metalinter .

Result:

      everything OK


go test:

     go test -v-cover

Result:

      === RUN   Test_startServer
      --- PASS: Test_startServer (2.00s)
      === RUN   Test_httpConnection
      --- PASS: Test_httpConnection (0.00s)
      === RUN   Test_rubbishURL_local
      --- PASS: Test_rubbishURL_local (0.00s)
      === RUN   Test_igcinfoapi_local
      --- PASS: Test_igcinfoapi_local (0.00s)
      === RUN   Test_PostAtInvalidURL
      --- PASS: Test_PostAtInvalidURL (0.00s)
      === RUN   Test_PostInvalidContent
      --- PASS: Test_PostInvalidContent (1.08s)
      === RUN   Test_PostValidContent
      --- PASS: Test_PostValidContent (2.39s)
      === RUN   Test_getAllIDs
      --- PASS: Test_getAllIDs (0.00s)
      === RUN   Test_getFields
      --- PASS: Test_getFields (0.00s)
      === RUN   Test_apiTtickerLatest
      --- PASS: Test_apiTtickerLatest (0.00s)
      === RUN   Test_apiTicker
      --- PASS: Test_apiTicker (0.01s)
      === RUN   Test_WebhookNewTrack
      --- PASS: Test_WebhookNewTrack (0.00s)
      === RUN   Test_cleanUp
      --- PASS: Test_cleanUp (0.02s)
      PASS
      coverage: 49.9% of statements
      ok  	bitbucket.org/isberg/paragliding	(cached)	coverage: 49.9% of statements
