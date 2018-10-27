package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	filePath   string = "./Config/clockTriggerConfig.txt"
	slackURL   string = "https://hooks.slack.com/services/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"
	getFromURL string = "https://calm-mesa-59678.herokuapp.com/paragliding/api/ticker/latest"
)

// this app runs in openstack. it will check if there is a new track positng and if so it will send a massage to
// a predefined Slack webhook
func main() {
	var lastTimestamp int64

	for {

		res, err := http.Get(getFromURL)
		if err != nil {
			postToSlack("unable to get content from website")
		} else if res.StatusCode == http.StatusOK {

			read, err1 := ioutil.ReadAll(res.Body)
			if err1 != nil {
				postToSlack("unable to read content from website")
			}

			strRead := string(read)
			strRead = strings.TrimSpace(strRead)

			newestTimestamp, err2 := strconv.ParseInt(strRead, 10, 64)
			if err2 != nil {
				message := "invalide value, not int64" + strconv.FormatInt(newestTimestamp, 10)
				postToSlack(message)
			}

			if newestTimestamp > lastTimestamp {
				// call slack webhook
				postToSlack(" a new IGC track has been posted")
			}
			lastTimestamp = newestTimestamp
		}
		// the time between checks is determend by a value in a configfile that may by changed at runtime
		a := readFromFile()
		if a < 1 {
			a = 1
		}
		for i := 0; i < a; i++ {
			time.Sleep(time.Minute)
		}
	}
}
