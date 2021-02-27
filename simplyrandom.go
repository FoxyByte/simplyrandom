package main

import (
	"log"
	"net/http"
)

func handleRequests() {
	http.HandleFunc("/random/mean", getRandoms)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
