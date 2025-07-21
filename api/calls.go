package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sarah/mongodb"
	"sort"

	api "github.com/VapiAI/server-sdk-go"
	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/joho/godotenv"
)

var VapiClient *vapiclient.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	VapiClient = createClient(os.Getenv("VAPI_API_KEY"))
}

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

	customerList := []*api.CreateCustomerDto{}
	for _, phoneNumber := range phoneNumbers {
		customerList = append(customerList, &api.CreateCustomerDto{
			Number: api.String(phoneNumber),
		})
	}

	if len(customerList) == 0 {
		http.Error(w, "No phone numbers provided", http.StatusBadRequest)
		return
	}

	resp, err := VapiClient.Calls.Create(context.Background(), &api.CreateCallDto{
		AssistantId:   api.String(assistantId),
		PhoneNumberId: api.String(assistantNumberId),
		Customers:     customerList,
	})

	if err != nil {
		log.Printf("Error creating call: %v", err)
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

	resp, err := VapiClient.Calls.Get(context.Background(), callId)

	if err != nil {
		log.Printf("Error getting call: %v", err)
		http.Error(w, "Failed to get call", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Call retrieved successfully: %+v\n", resp)
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

	calls, err := VapiClient.Calls.List(context.Background(), callListRequest)

	if err != nil {
		log.Printf("Error listing calls: %v", err)
		http.Error(w, "Failed to list calls", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}

// GetCallListByOrgId handles GET requests to retrieve all calls for an organization.
// This endpoint aggregates calls from all assistants belonging to the specified organization.
//
// HTTP Method: GET
// Endpoint: /calls/org
//
// Query Parameters:
//   - orgId: The organization ID to retrieve calls for (required)
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

	assistants := mongodb.GetOrganizationAssistants(orgID)

	calls := []*api.Call{}
	for _, assistant := range assistants {
		assistantCalls, err := VapiClient.Calls.List(context.Background(), &api.CallsListRequest{
			AssistantId: api.String(assistant.VapiAssistantId),
		})
		if err != nil {
			log.Printf("Error listing calls: %v", err)
			http.Error(w, "Failed to list calls", http.StatusInternalServerError)
			return
		}
		calls = append(calls, assistantCalls...)
	}

	sort.Slice(calls, func(i, j int) bool {
		return calls[i].CreatedAt.After(calls[j].CreatedAt)
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}
