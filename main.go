package main

import (
	"log"
	"net/http"
	"sarah/api"
)

func main() {
	// Health check endpoint
	http.HandleFunc("/test", test)

	// Call management endpoints
	http.HandleFunc("/calls/create", api.CreateCall)      // POST: Create a new call
	http.HandleFunc("/calls/list", api.ListCalls)         // GET: List all calls
	http.HandleFunc("/calls/call", api.GetCall)           // GET: Get specific call by ID
	http.HandleFunc("/calls/org", api.GetCallListByOrgId) // GET: Get calls by organization ID

	// Campaign management endpoints
	http.HandleFunc("/campaigns/create", api.CreateCampaign)   // POST: Create a new campaign
	http.HandleFunc("/campaigns/org", api.GetCampaignViaOrgID) // GET: Get campaigns by organization ID

	// Organization resource endpoints
	http.HandleFunc("/assistants/org", api.GetOrganizationAssistants)      // GET: Get assistants by organization ID
	http.HandleFunc("/contacts/org", api.GetOrganizationContacts)          // GET: Get contacts by organization ID
	http.HandleFunc("/phone_numbers/org", api.GetOrganizationPhoneNumbers) // GET: Get phone numbers by organization ID

	// Start the server on port 8080
	log.Println("Starting Sarah Campaign Management API on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
