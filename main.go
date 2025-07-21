package main

import (
	"log"
	"net/http"
	"sarah/api"
	"time"
)

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		startTime := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[API] Response: %s %s -> STATUS: %d completed in %v", r.Method, r.URL.Path, http.StatusOK, time.Since(startTime))
	}
}

func main() {
	// Health check endpoint
	http.HandleFunc("/test", test)

	// Call management endpoints
	http.HandleFunc("/calls/create", middleware(api.CreateCall))      // POST: Create a new call
	http.HandleFunc("/calls/list", middleware(api.ListCalls))         // GET: List all calls
	http.HandleFunc("/calls/call", middleware(api.GetCall))           // GET: Get specific call by ID
	http.HandleFunc("/calls/org", middleware(api.GetCallListByOrgId)) // GET: Get calls by organization ID

	// Campaign management endpoints
	http.HandleFunc("/campaigns/create", middleware(api.CreateCampaign))   // POST: Create a new campaign
	http.HandleFunc("/campaigns/org", middleware(api.GetCampaignViaOrgID)) // GET: Get campaigns by organization ID

	// Organization resource endpoints
	http.HandleFunc("/assistants/org", middleware(api.GetOrganizationAssistants))      // GET: Get assistants by organization ID
	http.HandleFunc("/contacts/org", middleware(api.GetOrganizationContacts))          // GET: Get contacts by organization ID
	http.HandleFunc("/phone_numbers/org", middleware(api.GetOrganizationPhoneNumbers)) // GET: Get phone numbers by organization ID

	// Start the server on port 8080
	log.Println("Starting Sarah AI Call assistant on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
