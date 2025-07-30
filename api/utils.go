package api

import (
	"encoding/json"
	"io"
	"net/http"
	"sarah/auth"
	mongodbTypes "sarah/types/mongodb"
	"strings"

	vapiApi "github.com/VapiAI/server-sdk-go"
)

// ExtractAuthHeader extracts the Bearer token from the Authorization header.
// This function removes the "Bearer " prefix from the Authorization header value.
//
// Parameters:
//   - r: HTTP request containing the Authorization header
//
// Returns:
//   - string: The extracted Bearer token without the "Bearer " prefix
//
// Example:
//
//	Authorization: Bearer abc123def456
//	Returns: abc123def456
func ExtractAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// ExtractPhoneNumbers extracts phone numbers from the request body.
// The function expects a JSON body with a "phoneNumbers" array field.
//
// Parameters:
//   - r: HTTP request containing the phone numbers in the request body
//
// Returns:
//   - []string: Array of phone numbers with whitespace trimmed
//
// Request Body Format:
//
//	{
//	  "phoneNumbers": ["+1234567890", "+1987654321"]
//	}
func ExtractPhoneNumbers(r *http.Request) []string {
	var phoneNumbers []string

	type requestBody struct {
		PhoneNumbers []string `json:"phoneNumbers"`
	}

	var body requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return []string{}
	}

	for _, phone := range body.PhoneNumbers {
		phoneNumbers = append(phoneNumbers, strings.TrimSpace(phone))
	}

	return phoneNumbers
}

// ExtractAssistantId extracts the assistant ID from the request query parameters.
// The function looks for the "assistantId" query parameter.
//
// Parameters:
//   - r: HTTP request containing the assistantId query parameter
//
// Returns:
//   - string: The assistant ID with whitespace trimmed
//
// Example URL: /calls/create?assistantId=asst_1234567890abcdef
func ExtractAssistantId(r *http.Request) string {
	assistantId := r.URL.Query().Get("assistantId")
	return strings.TrimSpace(assistantId)
}

// ExtractAssistantNumberId extracts the assistant number ID from the request query parameters.
// The function looks for the "assistantNumberId" query parameter.
//
// Parameters:
//   - r: HTTP request containing the assistantNumberId query parameter
//
// Returns:
//   - string: The assistant number ID with whitespace trimmed
//
// Example URL: /calls/create?assistantNumberId=phone_0987654321fedcba
func ExtractAssistantNumberId(r *http.Request) string {
	assistantNumberId := r.URL.Query().Get("assistantNumberId")
	return strings.TrimSpace(assistantNumberId)
}

// VerifyMethod checks if the HTTP request method is in the list of allowed methods.
// This function is used to ensure endpoints only accept the correct HTTP methods.
//
// Parameters:
//   - r: HTTP request to verify
//   - allowedMethods: Array of allowed HTTP methods (e.g., ["GET", "POST"])
//
// Returns:
//   - bool: True if the request method is allowed, false otherwise
//
// Example:
//
//	VerifyMethod(r, []string{"POST"}) // Only allows POST
//	VerifyMethod(r, []string{"GET", "POST"}) // Allows both GET and POST
func VerifyMethod(r *http.Request, allowedMethods []string) bool {
	for _, method := range allowedMethods {
		if r.Method == strings.ToUpper(method) {
			return true
		}
	}
	return false
}

// ExtractCallId extracts the call ID from the request query parameters.
// The function looks for the "callId" query parameter.
//
// Parameters:
//   - r: HTTP request containing the callId query parameter
//
// Returns:
//   - string: The call ID with whitespace trimmed
//
// Example URL: /calls/call?callId=call_abc123def456
func ExtractCallId(r *http.Request) string {
	callId := r.URL.Query().Get("callId")
	return strings.TrimSpace(callId)
}

// ExtractCallListRequest extracts a call list request from the request body.
// The function expects a JSON body with a "callListRequest" object field.
//
// Parameters:
//   - r: HTTP request containing the call list request in the request body
//
// Returns:
//   - *api.CallsListRequest: The extracted call list request, or nil if extraction fails
//
// Request Body Format:
//
//	{
//	  "callListRequest": {
//	    "assistantId": "asst_1234567890abcdef",
//	    "limit": 10,
//	    "offset": 0
//	  }
//	}
func ExtractCallListRequest(r *http.Request) *vapiApi.CallsListRequest {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		return nil
	}

	raw, ok := bodyMap["callListRequest"]
	if !ok || raw == nil {
		return nil
	}

	var callListRequest vapiApi.CallsListRequest
	if err := json.Unmarshal(raw, &callListRequest); err != nil {
		return nil
	}
	return &callListRequest
}

// ExtractCampaignId extracts the campaign ID from the request URL path.
// The function removes the "/campaigns/campaign/" prefix from the URL path.
//
// Parameters:
//   - r: HTTP request containing the campaign ID in the URL path
//
// Returns:
//   - string: The campaign ID with whitespace trimmed
//
// Example URL: /campaigns/campaign/camp_1234567890abcdef
func ExtractCampaignId(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/campaigns/campaign/")
	campaignId := strings.TrimSpace(path)
	return strings.TrimSpace(campaignId)
}

// ExtractOrgId extracts the organization ID from the request query parameters.
// The function looks for the "orgId" query parameter.
//
// Parameters:
//   - r: HTTP request containing the orgId query parameter
//
// Returns:
//   - string: The organization ID with whitespace trimmed
//
// Example URL: /campaigns/org?orgId=org_1234567890abcdef
func ExtractOrgId(r *http.Request) string {
	orgId, ok := auth.GetOrganizationID(r)
	if !ok {
		return ""
	}
	return strings.TrimSpace(orgId)
}

// ExtractCampaignCreateDto extracts a campaign creation DTO from the request body.
// The function expects a JSON body with a "campaignCreateRequest" object field.
// This matches the structure shown in sample_campaigns.json.
//
// Parameters:
//   - r: HTTP request containing the campaign creation request in the request body
//
// Returns:
//   - *mongodb.CampaignCreateDto: The extracted campaign creation DTO, or nil if extraction fails
//
// Request Body Format:
//
//	{
//	  "campaignCreateRequest": {
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
//	}
func ExtractCampaignCreateDto(r *http.Request) *mongodbTypes.Campaign {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		CampaignCreateRequest mongodbTypes.Campaign `json:"campaignCreateRequest"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.CampaignCreateRequest
}

func ExtractAssistantCreateDto(r *http.Request) *vapiApi.CreateAssistantDto {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		AssistantCreateRequest vapiApi.CreateAssistantDto `json:"assistantCreateRequest"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.AssistantCreateRequest
}

func ExtractAssistantUpdateDto(r *http.Request) *vapiApi.UpdateAssistantDto {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		AssistantUpdateRequest vapiApi.UpdateAssistantDto `json:"assistantUpdateRequest"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.AssistantUpdateRequest
}

func ExtractAssistant(r *http.Request) *mongodbTypes.Assistant {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		Assistant mongodbTypes.Assistant `json:"assistant"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.Assistant
}

func ExtractContact(r *http.Request) *mongodbTypes.Contact {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		Contact mongodbTypes.Contact `json:"contact"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.Contact
}
