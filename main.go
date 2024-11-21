package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/receipts/process", process).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	port := 3005
	log.Println("Starting on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), r))
}
