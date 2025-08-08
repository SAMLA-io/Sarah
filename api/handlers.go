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

	resp, err := sarah.CreateCall(assistantId, assistantNumberId, customers)

	if resp == nil {
		http.Error(w, "Failed to create call", http.StatusInternalServerError)
		return
	} else if err != nil {
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

	resp, err := sarah.GetCall(callId)

	if resp == nil {
		http.Error(w, "Failed to get call", http.StatusInternalServerError)
		return
	} else if err != nil {
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
// HTTP Method: POST
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
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callListRequest := ExtractCallListRequest(r)

	calls, err := sarah.ListCalls(callListRequest)

	if calls == nil {
		http.Error(w, "Failed to list calls", http.StatusInternalServerError)
		return
	} else if err != nil {
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

	calls, err := sarah.GetOrganizationCalls(orgID)

	sort.Slice(calls, func(i, j int) bool {
		return calls[i].CreatedAt.After(calls[j].CreatedAt)
	})
	if err != nil {
		http.Error(w, "Failed to get organization calls", http.StatusInternalServerError)
		return
	}

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
	assistants, err := mongodb.GetOrganizationAssistants(orgID)

	json.NewEncoder(w).Encode(assistants)
	if err != nil {
		http.Error(w, "Failed to get organization assistants", http.StatusInternalServerError)
		return
	}
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
//	(See API documentation for full schema details. https://docs.vapi.ai/api-reference/assistants/create)
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
//	{ ...mongodb insert one result object... }
func CreateAssistant(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantCreateDto := ExtractAssistantCreateDto(r)
	orgId := ExtractOrgId(r)

	if assistantCreateDto == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := sarah.CreateAsisstant(orgId, *assistantCreateDto)

	if err != nil {
		http.Error(w, "Failed to create assistant", http.StatusInternalServerError)
		return
	}

	if result == nil {
		http.Error(w, "Failed to create assistant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterAssistant handles POST requests to register an already existing assistant.
// This endpoint accepts an assistant registration request and registers the assistant in the database.
// This is useful when an assistant is already created in VapiAI and needs to be registered in the database manually.
// This endpoint will not create the assistant in VapiAI, it will only register the assistant in the database IF it exists in VapiAI.

// HTTP Method: POST
// Endpoint: /assistants/register
//
// Request Body:
//
//	{
//	  "assistant": { ... }
//	}
//
// The organization ID is obtained from the auth bearer token.
func RegisterAssistant(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistant := ExtractAssistant(r)
	orgId := ExtractOrgId(r)

	log.Printf("RegisterAssistant: Checking existence of assistant with ID: %s", assistant.VapiAssistantId)

	if !sarah.ExistsAssistant(assistant.VapiAssistantId) {
		http.Error(w, "Assistant does not exist in VapiAI", http.StatusBadRequest)
		return
	}

	result, err := mongodb.CreateAssistant(orgId, *assistant)

	if result == nil {
		http.Error(w, "Failed to create assistant", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to create assistant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdateAssistant handles PUT requests to update an existing assistant.
// This endpoint accepts an assistant update request and updates the assistant in the database.
//
// HTTP Method: PATCH
// Endpoint: /assistants/update
//
// Request Body:
//
// The request body must be a JSON object representing the assistant to update. The expected structure is:
//
//	{
//		  "id": "foo",
//		  "orgId": "foo",
//		  "createdAt": "foo",
//		  "updatedAt": "foo",
//		  "transcriber": { ... },
//		  "model": { ... },
//		  "voice": { ... },
//		  "firstMessage": "Hello! How can I help you today?",
//		  "firstMessageInterruptionsEnabled": false,
//		  "firstMessageMode": "assistant-speaks-first",
//		  "voicemailDetection": { ... },
//		  "clientMessages": "conversation-update",
//		  "serverMessages": "conversation-update",
//		  "maxDurationSeconds": 600,
//		  "backgroundSound": "off",
//		  "modelOutputInMessagesEnabled": false,
//		  "transportConfigurations": [ ... ],
//		  "observabilityPlan": { ... },
//		  "credentials": [ ... ],
//		  "hooks": [ ... ],
//		  "name": "foo",
//		  "voicemailMessage": "foo",
//		  "endCallMessage": "foo",
//		  "endCallPhrases": [ "foo" ],
//		  "compliancePlan": { ... },
//		  "metadata": {},
//		  "backgroundSpeechDenoisingPlan": { ... },
//		  "analysisPlan": { ... },
//		  "artifactPlan": { ... },
//		  "messagePlan": { ... },
//		  "startSpeakingPlan": { ... },
//		  "stopSpeakingPlan": { ... },
//		  "monitorPlan": { ... },
//		  "credentialIds": [ "foo" ],
//		  "server": { ... },
//		  "keypadInputPlan": { ... },
//		  "backgroundDenoisingEnabled": false
//		}
//
//		(See API documentation for full schema details. https://docs.vapi.ai/api-reference/assistants/update)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Assistant updated successfully, returns the updated assistant object
//   - 405 Method Not Allowed: If not using PUT method
//   - 400 Bad Request: If the request body is invalid
//
// Example Request:
//
//	PATCH /assistants/update
//	Content-Type: application/json
//	Authorization: Bearer <token>
//	{ ...assistant body as above... }
//
// Example Response:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//	{ ... vapi update assistant object... }
func UpdateAssistant(w http.ResponseWriter, r *http.Request) {

	if !VerifyMethod(r, []string{"PATCH"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantUpdateDto := ExtractAssistantUpdateDto(r)
	assistantId := ExtractAssistantId(r)

	result, err := sarah.UpdateAssistant(assistantId, *assistantUpdateDto)

	if result == nil {
		http.Error(w, "Failed to update assistant", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to update assistant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteAssistant handles DELETE requests to delete an existing assistant.
// This endpoint accepts an assistant deletion request and deletes the assistant from the database.
//
// HTTP Method: DELETE
// Endpoint: /assistants/delete
//
// Query Parameters:
//   - assistantId: The VapiAI assistant ID to delete (required)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Assistant deleted successfully, returns the deleted assistant object
//   - 405 Method Not Allowed: If not using DELETE method
//   - 500 Internal Server Error: If database operation fails
//
// Example Request:
//
//	DELETE /assistants/delete
//	Content-Type: application/json
//	Authorization: Bearer <token>
//	{ ...assistantId... }
//
// Example Response:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//	{ ... mongodb delete one result object... }
func DeleteAssistant(w http.ResponseWriter, r *http.Request) {

	if !VerifyMethod(r, []string{"DELETE"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantId := ExtractAssistantId(r)
	orgId := ExtractOrgId(r)

	result, err := sarah.DeleteAssistant(orgId, assistantId)

	if result == nil {
		http.Error(w, "Failed to delete assistant", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete assistant", http.StatusInternalServerError)
		return
	}

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
	campaign, err := sarah.CreateCampaign(*campaignCreateDto, orgId)

	if campaign == nil {
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
}

// UpdateCampaign handles PATCH requests to update an existing campaign.
// This endpoint accepts a campaign update request and updates the campaign in the database.
//
// HTTP Method: PATCH
// Endpoint: /campaigns/update
//
// Request Body:
//
//	{
//	  "campaignUpdateRequest": {
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
//   - 200 OK: Campaign updated successfully, returns the updated campaign
//   - 405 Method Not Allowed: If not using PATCH method
//   - 500 Internal Server Error: If database operation fails

// Example Response:
//
//	{
//	  "MatchedCount": 1,
//	  "ModifiedCount": 1,
//	  "UpsertedCount": 0,
//	  "UpsertedID": nil,
//	  "Acknowledged": true
//	}
func UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"PATCH"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignUpdateDto := ExtractCampaignUpdateDto(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.UpdateCampaign(orgId, *campaignUpdateDto)

	if result == nil {
		http.Error(w, "Failed to update campaign", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to update campaign", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteCampaign handles DELETE requests to delete an existing campaign.
// This endpoint accepts a campaign deletion request and deletes the campaign from the database.
//
// HTTP Method: DELETE
// Endpoint: /campaigns/delete
//
// Query Parameters:
//   - campaignId: The campaign ID to delete (required)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Campaign deleted successfully, returns the deleted campaign object
//   - 405 Method Not Allowed: If not using DELETE method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  "DeletedCount": 1,
//	  "Acknowledged": true
//	}
func DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"DELETE"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignId := ExtractCampaignId(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.DeleteCampaign(orgId, campaignId)

	if result == nil {
		http.Error(w, "Failed to delete campaign", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete campaign", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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

	campaigns, err := mongodb.GetCampaignByOrgId(orgId)

	if err != nil {
		http.Error(w, "Failed to get campaigns", http.StatusInternalServerError)
		return
	}

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

	contacts, err := mongodb.GetContactByOrgId(orgID)

	if err != nil {
		http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)
}

// CreateContact handles POST requests to create a new contact.
// This endpoint accepts a contact creation request and stores it in the database.
//
// HTTP Method: POST
// Endpoint: /contacts/create
//
// Request Body:
//
//	{
//	  "contact": {
//	    "name": "John Doe",
//	    "phoneNumber": "+1234567890",
//	    "email": "john.doe@example.com",
//	    "company": "Example Inc.",
//	    "position": "Software Engineer",
//	    "address": "123 Main St, Anytown, USA"
//	    "metadata": { ... }
//	  }
//	}
//
// Response:
//   - 200 OK: Contact created successfully, returns the created contact
//   - 405 Method Not Allowed: If not using POST method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  InsertedID: "507f1f77bcf86cd799439011",
//	  Acknowledged: true
//	}
//
// The organization ID is obtained from the auth bearer token.
func CreateContact(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contact := ExtractContact(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.CreateContact(orgId, *contact)

	if result == nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdateContact handles PATCH requests to update an existing contact.
// This endpoint accepts a contact update request and updates the contact in the database.
//
// HTTP Method: PATCH
// Endpoint: /contacts/update
//
// Request Body:
//
//	{
//	  "contact": {
//	    "name": "John Doe",
//	    "phoneNumber": "+1234567890",
//	    "email": "john.doe@example.com",
//	    "company": "Example Inc.",
//	    "position": "Software Engineer",
//	    "address": "123 Main St, Anytown, USA",
//	    "metadata": { ... }
//	  }
//	}
//
// Response:
//   - 200 OK: Contact updated successfully, returns the updated contact
//   - 405 Method Not Allowed: If not using PATCH method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  MatchedCount: 1,
//	  ModifiedCount: 1,
//	  UpsertedCount: 0,
//	  UpsertedID: nil,
//	  Acknowledged: true
//	}
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"PATCH"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contact := ExtractContact(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.UpdateContact(orgId, *contact)

	if result == nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteContact handles DELETE requests to delete an existing contact.
// This endpoint accepts a contact deletion request and deletes the contact from the database.
//
// HTTP Method: DELETE
// Endpoint: /contacts/delete
//
// Query Parameters:
//   - contactId: The contact ID to delete (required)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Contact deleted successfully, returns the deleted contact object
//   - 405 Method Not Allowed: If not using DELETE method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  "DeletedCount": 1,
//	  "Acknowledged": true
//	}
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"DELETE"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contactId := ExtractContactId(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.DeleteContact(orgId, contactId)

	if result == nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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

	phoneNumbers, err := mongodb.GetPhoneNumberByOrgId(orgID)

	if err != nil {
		http.Error(w, "Failed to get phone numbers", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phoneNumbers)
}

// CreatePhoneNumber handles POST requests to create a new phone number.
// This endpoint accepts a phone number creation request and stores it in the database.
//
// HTTP Method: POST
// Endpoint: /phone_numbers/create
//
// Request Body:
//
//	{
//	  "phoneNumber": {
//	    "name": "Main Office Line",
//	    "phoneNumber": "+1987654321",
//	    "vapiPhoneNumberId": "phone_0987654321fedcba"
//	  }
//	}
//
// Response:
//   - 200 OK: Phone number created successfully, returns the created phone number
//   - 405 Method Not Allowed: If not using POST method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  InsertedID: "507f1f77bcf86cd799439011",
//	  Acknowledged: true
//	}
//
// The organization ID is obtained from the auth bearer token.
func CreatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	phoneNumber := ExtractPhoneNumber(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.CreatePhoneNumber(orgId, *phoneNumber)

	if result == nil {
		http.Error(w, "Failed to create phone number", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to create phone number", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeletePhoneNumber handles DELETE requests to delete an existing phone number.
// This endpoint accepts a phone number deletion request and deletes the phone number from the database.
//
// HTTP Method: DELETE
// Endpoint: /phone_numbers/delete
//
// Query Parameters:
//   - phoneNumberId: The phone number ID to delete (required)
//
// The organization ID is obtained from the auth bearer token.
//
// Response:
//   - 200 OK: Phone number deleted successfully, returns the deleted phone number object
//   - 405 Method Not Allowed: If not using DELETE method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	{
//	  "DeletedCount": 1,
//	  "Acknowledged": true
//	}
//
// The organization ID is obtained from the auth bearer token.
func DeletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"DELETE"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	phoneNumberId := ExtractPhoneNumberId(r)
	orgId := ExtractOrgId(r)

	result, err := mongodb.DeletePhoneNumber(orgId, phoneNumberId)

	if result == nil {
		http.Error(w, "Failed to delete phone number", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete phone number", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
