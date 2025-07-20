package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"
)

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
