package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, struct{}{})
}
