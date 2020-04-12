package main

import (
	"fmt"
	//"os"
	//"io/ioutil"
	//"bufio"
	//"strings"
	"net/http"
	"log"
	//"github.com/gorilla/mux"
)

func main() {
	// Initialize router
	router := NewRouter()

	// Set up the repository
	RepoInit()

	fmt.Println("Listening on port 8080...")
	
	log.Fatal(http.ListenAndServe(":8080", router))
}