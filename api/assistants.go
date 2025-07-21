package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"
)

// GetOrganizationAssistants handles GET requests to retrieve all assistants for an organization.
// This endpoint returns all VapiAI assistants that belong to the specified organization.
//
// HTTP Method: GET
// Endpoint: /assistants/org
//
// Query Parameters:
//   - orgId: The organization ID to retrieve assistants for (required)
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
