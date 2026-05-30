package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/neo/week/", getIDs)

	logger.Println("Server is starting...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal(err)
	}
}
