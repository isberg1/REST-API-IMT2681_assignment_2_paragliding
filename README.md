 README Assignment 2 in IMT2681-2018


### What is this repository for? ###
An online service that will allow users to browse information about IGC files(international file format used by paragliders). The program will store IGC files metadata in a mongoDB Database. The system will generate events and it will monitor for new events happening from the outside services. The project will make use of Heroku and OpenStack.

Click [here](https://github.com/isberg1/REST-API-IMT2681_assignment_2_paragliding/blob/master/Assignment%202%20spesifications.md) for assignment requirements




Assignment URL: https://calm-mesa-59678.herokuapp.com/paragliding (no longer supported)

### How do I test the remote Heroku api ###

From bash terminal use the following commads:

GET /paragliding/

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/

GET /api

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api

POST /api/track

    curl  --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"url": "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}' -X POST https://calm-mesa-59678.herokuapp.com/paragliding/api/track

GET /api/track

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track

GET /api/track/<id>

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track/1

GET /api/track/<id>/<field>

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/track/1/pilot

GET /api/ticker/latest

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/latest

GET /api/ticker/

    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/

GET /api/ticker/<timestamp>

    !!! write timestamp at the end of command !!!
    curl --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/

POST /api/webhook/new_track/


     curl  --write-out "\n%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"web_hook_url": " https://calm-mesa-59678.herokuapp.com/test", "min_trigger_value" : 2 }' -X POST https://calm-mesa-59678.herokuapp.com/paragliding/api/webhook/new_track

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


for my setup i was using go version: go1.11 on Ubuntu 18.04

in order to use mongoDB i had to install it first se link(you may have to register first)

    https://docs.mongodb.com/manual/tutorial/install-mongodb-enterprise-on-ubuntu/


I needed to install "gcc" to run the tests

	sudo apt install gcc
	cd ~/
	git clone https://isberg@bitbucket.org/isberg/paragliding.git
	cd paragliding/
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

      $ go test . -v -cover
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
     --- PASS: Test_PostInvalidContent (0.99s)
     === RUN   Test_PostValidContent
     --- PASS: Test_PostValidContent (2.78s)
     === RUN   Test_getAllIDs
     --- PASS: Test_getAllIDs (0.01s)
     === RUN   Test_getFields
     --- PASS: Test_getFields (0.00s)
     === RUN   Test_apiTtickerLatest
     --- PASS: Test_apiTtickerLatest (0.00s)
     === RUN   Test_apiTicker
     --- PASS: Test_apiTicker (0.02s)
     === RUN   Test_WebhookNewTrack
     --- PASS: Test_WebhookNewTrack (0.01s)
     === RUN   Test_getWebhookByID
     --- PASS: Test_getWebhookByID (0.00s)
     === RUN   Test_adminCount
     --- PASS: Test_adminCount (0.00s)
     === RUN   Test_adminTrackDropTable
     --- PASS: Test_adminTrackDropTable (0.00s)
     === RUN   Test_cleanUp
     --- PASS: Test_cleanUp (0.02s)
     PASS
     coverage: 50.7% of statements
     ok  	bitbucket.org/isberg/paragliding	(cached)	coverage: 50.7% of statements


# choices and decisions

i choose the globalsign/mgo driver because i found better documentation for it


in order to check what webhocks should be posted to, i made a counter in the webhook document
i works like this:

* for every new IGC track post, all webhook counters are decremented by 1
* for all webhook document where counter == 0, post to them
* for all webhook document where counter == 0, reset counter to counter = min_trigger_value

i was unsure if we where to implement authentication for the admin endpoints, but decided to try, so i
found a library that allows for simple authentication.

my clocktrigger app is running in a tmux session in openstack. it checks a config file to se the timestamp
check interval. the check interval can be altered at runtime by editing the config file. the app posts to slack


in order to test webhook funtionality set URL subscription address to be https://calm-mesa-59678.herokuapp.com/test

post as many new igc track as your minimal_trigger_value

open website https://calm-mesa-59678.herokuapp.com/test to se webhook post

the file with the main function is called app.go
