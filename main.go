package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	// Initialize router
	router := NewRouter()

	// Set up the repository
	RepoInit()

	fmt.Println("Listening on port 8080...")
	
	log.Fatal(http.ListenAndServe(":8080", router))
}