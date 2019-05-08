package main

import (
	"fmt"
	"net/http"
)

// Ping sends back an "OK".
func Ping(w http.ResponseWriter, r *http.Request) {

	// Get rid off warning
	_ = r

	// Send response
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

/*
curl http://localhost:8080/ping
*/
