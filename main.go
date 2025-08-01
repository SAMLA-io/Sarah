package main

import (
	"log"
	"net/http"
	"sarah/api"
	"sarah/auth"
	"sarah/sarah"
	"time"
)

func main() {
	campaignScheduler := sarah.CampaignScheduler{}
	campaignScheduler.Start()

	http.HandleFunc("/", welcome)

	// Call management endpoints
	http.Handle("/calls/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateCall)))      // POST: Create a new call
	http.Handle("/calls/list", auth.VerifyingMiddleware(http.HandlerFunc(api.ListCalls)))         // GET: List all calls
	http.Handle("/calls/call", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCall)))           // GET: Get specific call by ID
	http.Handle("/calls/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCallListByOrgId))) // GET: Get calls by organization ID

	// Campaign management endpoints
	http.Handle("/campaigns/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetCampaignViaOrgID))) // GET: Get campaigns by organization ID
	http.Handle("/campaigns/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateCampaign)))   // POST: Create a new campaign
	http.Handle("/campaigns/update", auth.VerifyingMiddleware(http.HandlerFunc(api.UpdateCampaign)))   // PATCH: Update an existing campaign
	http.Handle("/campaigns/delete", auth.VerifyingMiddleware(http.HandlerFunc(api.DeleteCampaign)))   // DELETE: Delete an existing campaign

	// Organization resource endpoints
	http.Handle("/assistants/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationAssistants))) // GET: Get assistants by organization ID
	http.Handle("/assistants/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateAssistant)))        // POST: Create a new assistant
	http.Handle("/assistants/update", auth.VerifyingMiddleware(http.HandlerFunc(api.UpdateAssistant)))        // PATCH: Update an assistant
	http.Handle("/assistants/delete", auth.VerifyingMiddleware(http.HandlerFunc(api.DeleteAssistant)))        // DELETE: Delete an assistant
	http.Handle("/assistants/register", auth.VerifyingMiddleware(http.HandlerFunc(api.RegisterAssistant)))    // POST: Register an existing assistant

	http.Handle("/contacts/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreateContact)))        // POST: Create a new contact
	http.Handle("/contacts/update", auth.VerifyingMiddleware(http.HandlerFunc(api.UpdateContact)))        // PATCH: Update an existing contact
	http.Handle("/contacts/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationContacts))) // GET: Get contacts by organization ID
	http.Handle("/contacts/delete", auth.VerifyingMiddleware(http.HandlerFunc(api.DeleteContact)))        // DELETE: Delete an existing contact

	http.Handle("/phone_numbers/org", auth.VerifyingMiddleware(http.HandlerFunc(api.GetOrganizationPhoneNumbers))) // GET: Get phone numbers by organization ID
	http.Handle("/phone_numbers/create", auth.VerifyingMiddleware(http.HandlerFunc(api.CreatePhoneNumber)))        // POST: Create a new phone number
	http.Handle("/phone_numbers/delete", auth.VerifyingMiddleware(http.HandlerFunc(api.DeletePhoneNumber)))        // DELETE: Delete an existing phone number

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      http.DefaultServeMux,
	}

	log.Println("Starting Sarah AI Call assistant on port 8080...")
	log.Fatal(server.ListenAndServe())

}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Sarah AI Call Assistant!"))
}
