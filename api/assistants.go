package api

import (
	"encoding/json"
	"net/http"
	"sarah/supabase"
)

func GetOrganizationAssistants(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)
	assistants := supabase.GetAssistantByOrgId(orgID)

	json.NewEncoder(w).Encode(assistants)
}
