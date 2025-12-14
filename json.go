package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("responding with 5xx error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal JSON to response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
