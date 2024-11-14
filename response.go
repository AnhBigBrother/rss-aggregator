package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func errResponse(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Server error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"` // reflect tag for json.Marshal func
	}
	jsonResponse(w, code, errResponse{Error: msg})
}
