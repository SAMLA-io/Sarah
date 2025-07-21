package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"
)

// GetOrganizationPhoneNumbers handles GET requests to retrieve all phone numbers for an organization.
// This endpoint returns all VapiAI phone numbers that belong to the specified organization.
//
// HTTP Method: GET
// Endpoint: /phone_numbers/org
//
// Query Parameters:
//   - orgId: The organization ID to retrieve phone numbers for (required)
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
