package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"
)

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
