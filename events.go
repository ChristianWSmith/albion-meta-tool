package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

func getKillEventUrls() []string {
	var urls []string
	for offset := 0; offset <= 1000; offset += 50 {
		urls = append(urls, fmt.Sprintf("%s?limit=51&offset=%v", config.KillEventUrl, offset))
	}
	return urls
}

func eventMonitor() {
	// Make the HTTP GET request

	logInfo(fmt.Sprintf("Kill event urls: %v", getKillEventUrls()), nil)

	response, err := http.Get(config.KillEventUrl)
	if err != nil {
		logWarn("The HTTP request failed with error %s\n", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logError("Failed to read the response body: %s\n", err)
		return
	}

	// Use Gjson to parse and query the JSON response
	json := string(body)

	if !gjson.Valid(json) {
		logError("Invalid json", nil)
	}
	// Example: Iterate over all events and print the Killer's Name
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		killerName := value.Get("Killer.Name").String()
		logDebug(fmt.Sprintf("Event %s: Killer's Name: %s\n", key.String(), killerName), nil)
		return true // keep iterating
	})

}
