package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sarah/mongodb"
	"sarah/sarah"
	mongodbTypes "sarah/types/mongodb"
	"sort"
)

// CreateCall handles POST requests to create a new call using VapiAI.
// This endpoint creates outbound calls to multiple phone numbers using a specified assistant.
//
// HTTP Method: POST
// Endpoint: /calls/create
//
// Query Parameters:
//   - assistantId: The VapiAI assistant ID to use for the call (required)
//   - assistantNumberId: The VapiAI phone number ID to use for outbound calls (required)
//
// Request Body:
//
//	{
//	  "phoneNumbers": ["+1234567890", "+1987654321"]
//	}
//
// Response:
//   - 201 Created: Call created successfully, returns the call details
//   - 400 Bad Request: If no phone numbers are provided
//   - 405 Method Not Allowed: If not using POST method
//   - 500 Internal Server Error: If VapiAI API call fails
//
// Example Response:
//
//	{
//	  "id": "call_abc123def456",
//	  "assistantId": "asst_1234567890abcdef",
//	  "phoneNumberId": "phone_0987654321fedcba",
//	  "status": "queued",
//	  "createdAt": "2024-01-01T12:00:00Z"
//	}
func CreateCall(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantId := ExtractAssistantId(r)
	phoneNumbers := ExtractPhoneNumbers(r)
	assistantNumberId := ExtractAssistantNumberId(r)

	customers := []mongodbTypes.Customer{}
	for _, phoneNumber := range phoneNumbers {
		customers = append(customers, mongodbTypes.Customer{
			PhoneNumber: phoneNumber,
		})
	}

	resp := sarah.CreateCall(assistantId, assistantNumberId, customers)

	if resp == nil {
		http.Error(w, "Failed to create call", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Call created successfully: %+v\n", resp)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetCall handles GET requests to retrieve a specific call by its ID.
// This endpoint fetches detailed information about a single call from VapiAI.
//
// HTTP Method: GET
// Endpoint: /calls/call
//
// Query Parameters:
//   - callId: The VapiAI call ID to retrieve (required)
//
// Response:
//   - 200 OK: Call retrieved successfully, returns the call details
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If VapiAI API call fails
//
// Example Response:
//
//	{
//	  "id": "call_abc123def456",
//	  "assistantId": "asst_1234567890abcdef",
//	  "phoneNumberId": "phone_0987654321fedcba",
//	  "status": "completed",
//	  "duration": 120,
//	  "createdAt": "2024-01-01T12:00:00Z",
//	  "endedAt": "2024-01-01T12:02:00Z"
//	}
func GetCall(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callId := ExtractCallId(r)

	resp := sarah.GetCall(callId)

	if resp == nil {
		http.Error(w, "Failed to get call", http.StatusInternalServerError)
		return
	}

	log.Printf("Call retrieved successfully: %+v\n", resp)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// ListCalls handles GET requests to list calls based on specified criteria.
// This endpoint retrieves a list of calls from VapiAI using optional filtering parameters.
//
// HTTP Method: GET
// Endpoint: /calls/list
//
// Request Body (optional):
//
//	{
//	  "callListRequest": {
//	    "assistantId": "asst_1234567890abcdef",
//	    "limit": 10,
//	    "offset": 0,
//	    "status": "completed"
//	  }
//	}
//
// Response:
//   - 200 OK: Calls retrieved successfully, returns an array of calls
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If VapiAI API call fails
//
// Example Response:
//
//	[
//	  {
//	    "id": "call_abc123def456",
//	    "assistantId": "asst_1234567890abcdef",
//	    "status": "completed",
//	    "createdAt": "2024-01-01T12:00:00Z"
//	  }
//	]
func ListCalls(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callListRequest := ExtractCallListRequest(r)

	calls := sarah.ListCalls(callListRequest)

	if calls == nil {
		http.Error(w, "Failed to list calls", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}

// GetCallListByOrgId handles GET requests to retrieve all calls for an organization.
// This endpoint aggregates calls from all assistants belonging to the organization from the auth bearer token.
//
// HTTP Method: GET
// Endpoint: /calls/org
//
// Response:
//   - 200 OK: Calls retrieved successfully, returns an array of calls sorted by creation date
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database or VapiAI API calls fail
//
// Example Response:
//
//	[
//	  {
//	    "id": "call_abc123def456",
//	    "assistantId": "asst_1234567890abcdef",
//	    "status": "completed",
//	    "createdAt": "2024-01-01T12:00:00Z"
//	  }
//	]
func GetCallListByOrgId(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)

	calls := sarah.GetOrganizationCalls(orgID)

	sort.Slice(calls, func(i, j int) bool {
		return calls[i].CreatedAt.After(calls[j].CreatedAt)
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}

// GetOrganizationAssistants handles GET requests to retrieve all assistants for an organization.
// This endpoint returns all VapiAI assistants that belong to the organization from the auth bearer token.
//
// HTTP Method: GET
// Endpoint: /assistants/org
//
// Response:
//   - 200 OK: Assistants retrieved successfully, returns an array of assistants
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	[
//	  {
//	    "id": "asst_1234567890abcdef",
//	    "name": "Insurance Reminder Assistant",
//	    "vapiAssistantId": "asst_1234567890abcdef",
//	    "organizationId": "org_1234567890abcdef"
//	  }
//	]
func GetOrganizationAssistants(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)
	assistants := mongodb.GetOrganizationAssistants(orgID)

	json.NewEncoder(w).Encode(assistants)
}

// CreateAssistant handles POST requests to create a new assistant for an organization.
//
// HTTP Method: POST
// Endpoint: /assistants/create
//
// Request Body:
//
//	The request body must be a JSON object representing the assistant to create. The expected structure is:
//
//	{
//	  "id": "foo",
//	  "orgId": "foo",
//	  "createdAt": "foo",
//	  "updatedAt": "foo",
//	  "transcriber": { ... },
//	  "model": { ... },
//	  "voice": { ... },
//	  "firstMessage": "Hello! How can I help you today?",
//	  "firstMessageInterruptionsEnabled": false,
//	  "firstMessageMode": "assistant-speaks-first",
//	  "voicemailDetection": { ... },
//	  "clientMessages": "conversation-update",
//	  "serverMessages": "conversation-update",
//	  "maxDurationSeconds": 600,
//	  "backgroundSound": "off",
//	  "modelOutputInMessagesEnabled": false,
//	  "transportConfigurations": [ ... ],
//	  "observabilityPlan": { ... },
//	  "credentials": [ ... ],
//	  "hooks": [ ... ],
//	  "name": "foo",
//	  "voicemailMessage": "foo",
//	  "endCallMessage": "foo",
//	  "endCallPhrases": [ "foo" ],
//	  "compliancePlan": { ... },
//	  "metadata": {},
//	  "backgroundSpeechDenoisingPlan": { ... },
//	  "analysisPlan": { ... },
//	  "artifactPlan": { ... },
//	  "messagePlan": { ... },
//	  "startSpeakingPlan": { ... },
//	  "stopSpeakingPlan": { ... },
//	  "monitorPlan": { ... },
//	  "credentialIds": [ "foo" ],
//	  "server": { ... },
//	  "keypadInputPlan": { ... },
//	  "backgroundDenoisingEnabled": false
//	}
//
//	(See API documentation for full schema details.)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Assistant created successfully, returns the created assistant object
//   - 405 Method Not Allowed: If not using POST method
//   - 400 Bad Request: If the request body is invalid
//
// Example Request:
//
//	POST /assistants/create
//	Content-Type: application/json
//	Authorization: Bearer <token>
//	{ ...assistant body as above... }
//
// Example Response:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//	{ ...assistant object... }
func CreateAssistant(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantCreateDto := ExtractAssistantCreateDto(r)
	orgId := ExtractOrgId(r)

	result := sarah.CreateAsisstant(orgId, *assistantCreateDto)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// CreateCampaign handles POST requests to create a new campaign.
// This endpoint accepts a campaign creation request and stores it in the database.
//
// HTTP Method: POST
// Endpoint: /campaigns/create
//
// Request Body:
//
//	{
//	  "campaignCreateRequest": {
//	    "name": "Weekly Insurance Reminders",
//	    "assistant_id": "asst_1234567890abcdef",
//	    "phone_number_id": "phone_0987654321fedcba",
//	    "schedule_plan": {
//	      "before_day": 3,
//	      "after_day": 0,
//	      "week_days": [1, 3, 5],
//	      "month_days": [],
//	      "year_months": []
//	    },
//	    "customers": [
//	      {
//	        "phone_number": "+1234567890",
//	        "day_number": 15,
//	        "month_number": 3,
//	        "week_day": 1,
//	        "custom_date": null,
//	        "expiry_date": "2024-12-31T23:59:59Z"
//	      }
//	    ],
//	    "type": "recurrent_weekly",
//	    "status": "active",
//	    "start_date": "2024-01-01T00:00:00Z",
//	    "end_date": "2024-12-31T23:59:59Z",
//	    "timezone": "America/New_York"
//	  }
//	}
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Campaign created successfully, returns the created campaign
//   - 405 Method Not Allowed: If not using POST method
//   - 500 Internal Server Error: If database operation fails
func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignCreateDto := ExtractCampaignCreateDto(r)
	orgId := ExtractOrgId(r)

	// Adds the campaign to the database
	campaign := sarah.CreateCampaign(*campaignCreateDto, orgId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
}

// GetCampaignViaOrgID handles GET requests to retrieve all campaigns for an organization.
// This endpoint returns all campaigns associated with the organization from the auth bearer token.
//
// HTTP Method: GET
// Endpoint: /campaigns/org
//
// Response:
//   - 200 OK: Returns an array of campaigns for the organization
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	[
//	  {
//	    "name": "Weekly Insurance Reminders",
//	    "assistant_id": "asst_1234567890abcdef",
//	    "phone_number_id": "phone_0987654321fedcba",
//	    "schedule_plan": { ... },
//	    "customers": [ ... ],
//	    "type": "recurrent_weekly",
//	    "status": "active",
//	    "start_date": "2024-01-01T00:00:00Z",
//	    "end_date": "2024-12-31T23:59:59Z",
//	    "timezone": "America/New_York"
//	  }
//	]
func GetCampaignViaOrgID(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgId := ExtractOrgId(r)

	campaigns := mongodb.GetCampaignByOrgId(orgId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaigns)
}

// GetOrganizationContacts handles GET requests to retrieve all contacts for an organization.
// This endpoint returns all customer contacts that belong to the organization from the auth bearer token.
//
// HTTP Method: GET
// Endpoint: /contacts/org
//
// Response:
//   - 200 OK: Contacts retrieved successfully, returns an array of contacts
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	[
//	  {
//	    "id": "contact_1234567890abcdef",
//	    "name": "John Doe",
//	    "phoneNumber": "+1234567890",
//	    "email": "john.doe@example.com",
//	    "organizationId": "org_1234567890abcdef"
//	  }
//	]
func GetOrganizationContacts(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)

	contacts := mongodb.GetContactByOrgId(orgID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)
}

// GetOrganizationPhoneNumbers handles GET requests to retrieve all phone numbers for an organization.
// This endpoint returns all VapiAI phone numbers that belong to the organization from the auth bearer token.
//
// HTTP Method: GET
// Endpoint: /phone_numbers/org
//
// Response:
//   - 200 OK: Phone numbers retrieved successfully, returns an array of phone numbers
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	[
//	  {
//	    "id": "phone_0987654321fedcba",
//	    "phoneNumber": "+1987654321",
//	    "vapiPhoneNumberId": "phone_0987654321fedcba",
//	    "organizationId": "org_1234567890abcdef"
//	  }
//	]
func GetOrganizationPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)

	phoneNumbers := mongodb.GetPhoneNumberByOrgId(orgID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phoneNumbers)
}
