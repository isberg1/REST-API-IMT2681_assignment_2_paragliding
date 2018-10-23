package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	FilePath  string = "./ClockTrigger/Config/clockTriggerConfig.txt"
	SlackURL  string = "https://hooks.slack.com/services/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"
	PostToURL string = "http://localhost:8080/paragliding/api/ticker/latest"
)

func main() {
	var lastTimestamp int64

	for {

		res, err := http.Get(PostToURL)
		if err != nil {
			postToSlack("unable to get content from website")
		} else {
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

		a := readFromFile()
		if a < 1 {
			a = 1
		}
		for i := 0; i < a; i++ {
			time.Sleep(time.Minute)
		}
	}
}
