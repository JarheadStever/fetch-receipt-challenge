package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)


var processedReceipts = make(map[uuid.UUID]int)

func main() {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/receipts/process", Process).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", RetrievePoints).Methods("GET")

	port := 3005
	log.Println("Starting on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), r))
}
