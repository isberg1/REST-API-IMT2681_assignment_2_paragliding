package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

// calculates and returns the time the application has been running
func getUptime() string {
	// base values to be used in conditional statements and calculations
	var min = 60
	var hour = min * 60
	var day = hour * 24
	var week = day * 7
	var month = day * 30
	var year = day * 365

	// get the runtime duration in seconds
	ti := time.Since(StartUpTime).Seconds()
	elapsed := int(ti)
	// the string to be returned
	timeString := ""

	a := make([]interface{}, 0)
	a = append(a, "P", year, "Y", month, "M", week, "W", day, "D", "T", hour, "H", min, "M")

	for i := 0; i < len(a); i++ {
		// use type assertion to determine if "a" is "string" or "int"
		if str, ok := a[i].(string); ok {
			timeString += str
		} else if nr, ok := a[i].(int); ok {
			// if time elapsed since startup is more then a year/month/week/day/hour/min
			if elapsed >= nr {
				temp := elapsed / nr             // number of units of year/month/week/day/hour/min
				elapsed -= nr * temp             // subtract from total the number of units
				timeString += strconv.Itoa(temp) // add number of units to the return string
			} else { // if time elapsed since startup is less than 1, then skip 1 iteration
				i++
			}
		}
	}
	//convert remaining time in seconds to string
	timeString += strconv.Itoa(elapsed) + "S"

	return timeString
}

// gets the current version from herocu
func getVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		version = unavalabeVersinNr
	}
	return version
}

// responds to URL: "/paragliding/api/"
func api(w http.ResponseWriter, r *http.Request) {

	//set http header content-type
	http.Header.Add(w.Header(), "content-type", "application/json")
	//get relevant info
	info := GetIgcinfoAPI{Uptime: getUptime(), Info: infoSting, Version: getVersion()}
	//convert info to json and write back to the client
	err := json.NewEncoder(w).Encode(&info)
	if err != nil {
		http.Error(w, "serverside error(api)", http.StatusInternalServerError)
	}
}
