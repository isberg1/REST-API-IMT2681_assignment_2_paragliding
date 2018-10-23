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

type SlackMessage struct {
	Text string `json:"text"`
}

func readFromFile() int {
	//open file
	file, err := os.Open(FilePath)
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

func postToSlack(str string) error {

	message, err := json.Marshal(SlackMessage{Text: str})
	if err != nil {
		return err
	}

	res, err1 := http.Post(SlackURL, "application/json", bytes.NewBuffer(message))
	if err1 != nil {
		return err1
	}

	var text string
	json.NewDecoder(res.Body).Decode(&text)

	fmt.Println(string(message))
	fmt.Println(text)

	return nil
}
