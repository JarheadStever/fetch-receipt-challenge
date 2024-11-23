package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProcessResponse struct {
	ID uuid.UUID `json:"id"`
}

type RetrievePointsResponse struct {
	Points int `json:"points"`
}

/*
Validate an input receipt, calculate its points, generate an ID, store it's score, and return the ID
*/
func Process(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	receipt := Receipt{}
	if err := json.Unmarshal(body, &receipt); err != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

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

	processedReceipts[id] = receipt.CountPoints()

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(json); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
Retrieve the points value of a receipt given its UUID in the URL (param: "id")
*/
func RetrievePoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points, exists := processedReceipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	json, err := json.Marshal(RetrievePointsResponse{Points: points})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(json); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
