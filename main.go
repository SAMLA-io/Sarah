package main

import (
	"log"
	"net/http"
	"sarah/api"
	"sarah/auth"
)

func main() {
	// Health check endpoint
	http.HandleFunc("/test", test)

	// Call management endpoints
	http.Handle("/calls/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateCall)))      // POST: Create a new call
	http.Handle("/calls/list", auth.VerifyingMiddleware(http.HandlerFunc(api.ListCalls)))         // GET: List all calls
	http.Handle("/calls/call", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCall)))           // GET: Get specific call by ID
	http.Handle("/calls/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCallListByOrgId))) // GET: Get calls by organization ID

	// Campaign management endpoints
	http.Handle("/campaigns/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateCampaign)))   // POST: Create a new campaign
	http.Handle("/campaigns/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCampaignViaOrgID))) // GET: Get campaigns by organization ID

	// Organization resource endpoints
	http.Handle("/assistants/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationAssistants)))      // GET: Get assistants by organization ID
	http.Handle("/contacts/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationContacts)))          // GET: Get contacts by organization ID
	http.Handle("/phone_numbers/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationPhoneNumbers))) // GET: Get phone numbers by organization ID

	// Start the server on port 8080
	log.Println("Starting Sarah AI Call assistant on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
