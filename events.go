package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

func eventMonitor() {
	// Make the HTTP GET request
	response, err := http.Get(config.KillEventUrl)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read the response body: %s\n", err)
		return
	}

	// Use Gjson to parse and query the JSON response
	json := string(body)

	// Example query to get the first event's Killer's Name
	name := gjson.Get(json, "0.Killer.Name")
	fmt.Printf("Killer's Name: %s\n", name.String())

	// Example: Iterate over all events and print the Killer's Name
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		killerName := value.Get("Killer.Name").String()
		fmt.Printf("Event %s: Killer's Name: %s\n", key.String(), killerName)
		return true // keep iterating
	})

}
