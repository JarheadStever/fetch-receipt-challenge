package main

import (
	"encoding/json"
	"net/http"
	// "github.com/gorilla/mux"
)

func process(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"output": "success"}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getPoints(w http.ResponseWriter, r *http.Request) {

	response := map[string]int{"points": 23}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
