package main

import (
	"log"
	"net/http"
	"sarah/api"
)

func main() {
	http.HandleFunc("/calls", api.CreateCall)
	http.HandleFunc("/calls/{callId}", api.GetCall)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
