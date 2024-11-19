package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseJson(w, 200, struct {
		Message string `json:"message"`
	}{Message: "Hello world!"})
}
