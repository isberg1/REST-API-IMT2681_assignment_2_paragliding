 README #

### What is this repository for? ###
Assignment 1 in IMT2681-2018

Assignment URL: https://sheltered-garden-37170.herokuapp.com/igcinfo/api/

### How do I test the remote Heroku api ###

From bash terminal use the following commads:

Get information about application:

     $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://sheltered-garden-37170.herokuapp.com/igcinfo/api/
     
Output
	     
         {"Uptime":"PT6M0S","Info":"Service for IGC tracks.","Version":"1.0.0"}
          200 application/json

post content:

      curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"url": "https://raw.githubusercontent.com/marni/goigc/master/testdata/optimize-long-flight-1.igc"}' -X POST https://sheltered-garden-37170.herokuapp.com/igcinfo/api/igc

Output:
     
     {"id":"IGC_file_1"}
     201 application/json
     
Get array with ID's of all strored objects:

    $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://sheltered-garden-37170.herokuapp.com/igcinfo/api/igc

Output:
     
     [{"id":"IGC_file_1"},]
     200 application/json
     
Get meta information about a object:

     $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://sheltered-garden-37170.herokuapp.com/igcinfo/api/igc/IGC_file_1

Output:

    {"h_date":"2017-08-07 00:00:00 +0000 UTC","pilot":"Pascal GENIN","glider":"LS 6","glider_id":"D-5860","track_length":507}
    200 application/json
     
Get a spesified meta infromation field from a object:

     curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET https://sheltered-garden-37170.herokuapp.com/igcinfo/api/igc/IGC_file_1/glider_id

Output:

    D-5860
    200 text/plain


Delete all content 

    $ curl  --write-out "%{http_code}"  -X DELETE https://sheltered-garden-37170.herokuapp.com/igcinfo/api/drop_table
    
Output
    
     200

### How do I use the program? ###
The following has been tested on 2 ubuntu 18.04, one dualbooted and one VM, and 1 ubuntu 16.04 VM

In my go setup i was using go version: go1.11
I nedded to install "gcc" to run the tests

	sudo apt install gcc
	cd ~/
	git clone git@bitbucket.org:isberg/igcinfo.git
	cd igcinfo/go-getting-started/
	go run .
	
wait until "go run ." is done. this may take some time
	
Open an other bash terminal and use the following commads:
The ouput will be the same as when testing the remote Heroku URL exept for the "Version" will say unavalable


Get information about application:

     $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/igcinfo/api/
     
Output
	     
         {"Uptime":"PT6M0S","Info":"Service for IGC tracks.","Version":"1.0.0"}
          200 application/json

post content:

      curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -d '{"url": "http://raw.githubusercontent.com/marni/goigc/master/testdata/optimize-long-flight-1.igc"}' -X POST http://localhost:8080/igcinfo/api/igc

Output:
     
     {"id":"IGC_file_1"}
     201 application/json
     
Get array with ID's of all strored objects:

    $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/igcinfo/api/igc

Output:
     
     [{"id":"IGC_file_1"},]
     200 application/json
     
Get meta information about a object:

     $ curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/igcinfo/api/igc/IGC_file_1

Output:

    {"h_date":"2017-08-07 00:00:00 +0000 UTC","pilot":"Pascal GENIN","glider":"LS 6","glider_id":"D-5860","track_length":507}
    200 application/json
     
Get a spesified meta infromation field from a object:

     curl  --write-out "%{http_code} %{content_type}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/igcinfo/api/igc/IGC_file_1/glider_id

Output:

    D-5860
    200 text/plain


Delete all content 

    $ curl  --write-out "%{http_code}"  -X DELETE http://localhost:8080/igcinfo/api/drop_table
    
Output
    
     200



### code quality ###

code quality checking:

test 1:  

      $ go tool vet --all .
Result:

     Everything OK

test 2;

     $ golint .
     
Result:

     Everything OK 
     
test 3:

    $ go fmt . 

Result:

     Everything OK
     
test 4:

    $ go test .
    
Result:

    PASS
    ok  	github.com/heroku/go-getting-started	4.563s
     
test 5:
     
    $ go test -cover

Result:
    
    PASS
    coverage: 77.8% of statements
    ok  	github.com/heroku/go-getting-started	3.918s
      
### Who do I talk to? ###
Alexander Jakobsen, 16BITSEC, Studentnr: 473151, alexajak@stud.ntnu.no