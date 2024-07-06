package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

func generateReport() [][]string {

	// TODO: finish this
	// get all of the events
	// filter for those with 1 participant
	// get all the builds
	// get all the builds which we can calculate the cost of fully and do it
	// get all of the items from those builds
	// for each build, for each item in it, add the silver gained/lost, kills/deaths, cost
	// for each item stats, divide the cost by the kills+deaths to get an average
	// format items to csv

	response := [][]string{}
	events, err := queryAllEvents()
	if err != nil {
		log.Error("Failed to query all events: ", err)

	}

	for _, event := range events {
		response = append(response, []string{fmt.Sprintf("%v", event)})
	}

	return response
}

// Handler function for the endpoint
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response in CSV format
	response := generateReport()

	// Encode and send the CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=report.csv")
	writer := csv.NewWriter(w)
	for _, record := range response {
		if err := writer.Write(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func startAPI() {
	http.HandleFunc("/report", reportHandler)
	log.Info("Server starting on port ", config.Port, "...")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
