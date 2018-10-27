package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// SlackMessage struct that is comparable with Slack webhooks
type SlackMessage struct {
	Text string `json:"text"`
}

// reads the time between checks from file
func readFromFile() int {
	//open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("unable to open config file", file)
		panic(err)
	} else {
		defer file.Close()
	}

	read, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		fmt.Println("unable to read config file", read)
		panic(err)
	}

	strRead := string(read)
	strRead = strings.TrimPrefix(strRead, "checkInterval =")
	strRead = strings.TrimSpace(strRead)
	checkInterval, err3 := strconv.Atoi(strRead)
	if err3 != nil {
		fmt.Println("invalid value in config file", checkInterval)
		panic(err)
	}

	return checkInterval
}

// posts messages to a Slack webhook
func postToSlack(str string) {

	message, err := json.Marshal(SlackMessage{Text: str})
	if err != nil {
		fmt.Println(err, "1")
		return
	}

	_, err1 := http.Post(slackURL, "application/json", bytes.NewBuffer(message))
	if err1 != nil {
		fmt.Println(err1, "2")
		return
	}

}
