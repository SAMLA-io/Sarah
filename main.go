package main

import (
	"log"
	"net/http"
	"sarah/api"
)

func main() {
	http.HandleFunc("/calls", api.CreateCall)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
