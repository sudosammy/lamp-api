package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

/*
@sudosammy
Credit to ChatGPT for writing the first draft of this.
*/

var apiKey string                 // Set the API key global
var lampState = make(map[int]int) // lampState is a map that stores the current state of each lamp

func main() {
	if os.Getenv("LAMP_APIKEY") == "" {
		fmt.Println("LAMP_APIKEY environment variable must be set.")
		os.Exit(1)
	}
	apiKey = os.Getenv("LAMP_APIKEY")

	// Initialize the lamp states
	lampState[0] = 0
	lampState[1] = 0

	// Set up the http server
	mux := http.NewServeMux()
	mux.HandleFunc("/lamp/", noCacheMiddleware(apiKeyMiddleware(lampHandler)))

	// Start the server
	fmt.Println("Listening on port 8000...")
	http.ListenAndServe(":8000", mux)
}

// middleware that sets the Cache-Control header on the response
func noCacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Cache-Control header
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// middleware that checks for the API key in the request
func apiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the API key is provided in the request
		if r.Header.Get("x-api-key") != apiKey {
			http.Error(w, "Invalid API key", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func lampHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Authenticated Request: " + r.URL.String())

	// Get the lamp ID from the request
	id, err := strconv.Atoi(r.URL.Path[len("/lamp/"):])
	if err != nil {
		http.Error(w, "Invalid lamp ID", http.StatusBadRequest)
		return
	}

	// Check if the lamp ID is valid (0 or 1)
	if id != 0 && id != 1 {
		http.Error(w, "Invalid lamp ID", http.StatusBadRequest)
		return
	}

	// Handle GET requests
	if r.Method == http.MethodGet {
		// Get the current state of the lamp
		state := lampState[id]
		// Write the state to the response
		w.Write([]byte(strconv.Itoa(state)))
		return
	}

	// Handle POST requests
	if r.Method == http.MethodPost {
		// Toggle the state of the lamp
		if lampState[id] == 0 {
			lampState[id] = 1
		} else {
			lampState[id] = 0
		}

		// Write a message to the response
		w.Write([]byte("Lamp state toggled"))
	}
}
