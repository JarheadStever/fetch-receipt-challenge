package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var processedReceipts = make(map[uuid.UUID]int)

func main() {

	var portNumber int
	flag.IntVar(&portNumber, "port", 3005, "Port used by Receipt Service")
	flag.Parse()
	if 1 > portNumber || portNumber > 65535 {
		log.Fatal("Port must be between 1 and 65535")
	}

	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", Process).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", RetrievePoints).Methods("GET")

	log.Println("Starting Receipt Service on port", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", portNumber), r))
}
