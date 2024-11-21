package main

import (
	"encoding/json"
	"io"
	"net/http"
	// "github.com/gorilla/mux"
	"github.com/google/uuid"
)

type ProcessResponse struct {
	ID uuid.UUID `json:"id"`
}

type RetrievePointsResponse struct {
	Points int64 `json:"points"`
}


func Process(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil { 

		http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	receipt := Receipt{}
	json.Unmarshal(body, &receipt)

	if err := receipt.Validate(); err != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

    id := uuid.New()
    json, err := json.Marshal(ProcessResponse{ID: id})
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// TODO: add to "database"

	w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

func RetrievePoints(w http.ResponseWriter, r *http.Request) {

	response := map[string]int{"points": 23}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
