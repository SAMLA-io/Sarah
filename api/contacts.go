package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"
)

// GetOrganizationContacts handles GET requests to retrieve all contacts for an organization.
// This endpoint returns all customer contacts that belong to the specified organization.
//
// HTTP Method: GET
// Endpoint: /contacts/org
//
// Query Parameters:
//   - orgId: The organization ID to retrieve contacts for (required)
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
