package main

import (
	"log"
	"net/http"
	"sarah/api"
)

func main() {
	http.HandleFunc("/calls/create", api.CreateCall)
	http.HandleFunc("/calls/list", api.ListCalls)
	http.HandleFunc("/calls/call/", api.GetCall)

	http.HandleFunc("/campaigns/create", api.CreateCampaign)
	http.HandleFunc("/campaigns/org/", api.GetCampaignViaOrgID)
	http.HandleFunc("/campaigns/campaign/", api.GetCampaignViaCampaignID)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
