package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "jstew",
		"ssn":  "509-22-3233",
	})

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://206.128.153.183:12970", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
