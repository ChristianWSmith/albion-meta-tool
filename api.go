package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define a struct to represent the request body
type RequestBody struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// Define a struct to represent the response body
type ResponseBody struct {
	Message string `json:"message"`
}

// Handler function for the endpoint
func handler(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set default values for optional fields if they are nil
	defaultName := "Guest"
	defaultEmail := "guest@example.com"

	name := defaultName
	email := defaultEmail

	if reqBody.Name != nil {
		name = *reqBody.Name
	}

	if reqBody.Email != nil {
		email = *reqBody.Email
	}

	// Create a response
	response := ResponseBody{
		Message: "Hello, " + name + "! Your email is " + email,
	}

	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func startAPI() {
	http.HandleFunc("/hello", handler)
	log.Info("Server starting on port ", config.Port, "...")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
