package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

// Handler function for the endpoint
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response in CSV format
	response := [][]string{
		{"col1", "col2", "col3"},
		{"val1", "val2", "val3"},
		{"val4", "val5", "val6"},
	}

	// Encode and send the CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=response.csv")
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
