package main

import (
	"log"
	"net/http"
	"sarah/api"
)

func main() {

	http.HandleFunc("/test", test)

	http.HandleFunc("/calls/create", api.CreateCall)
	http.HandleFunc("/calls/list", api.ListCalls)
	http.HandleFunc("/calls/call", api.GetCall)
	http.HandleFunc("/calls/org", api.GetCallListByOrgId)

	http.HandleFunc("/campaigns/create", api.CreateCampaign)
	http.HandleFunc("/campaigns/org", api.GetCampaignViaOrgID)

	http.HandleFunc("/assistants/org", api.GetOrganizationAssistants)

	http.HandleFunc("/contacts/org", api.GetOrganizationContacts)

	http.HandleFunc("/phone_numbers/org", api.GetOrganizationPhoneNumbers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
