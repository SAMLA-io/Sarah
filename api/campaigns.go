package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sarah/supabase"

	api "github.com/VapiAI/server-sdk-go"
)

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func GetCampaignViaCampaignID(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignId := ExtractCampaignId(r)

	campaign, err := VapiClient.Campaigns.CampaignControllerFindOne(context.Background(), campaignId)
	if err != nil {
		http.Error(w, "Failed to get campaign", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
}

func GetCampaignViaOrgID(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgId := ExtractOrgId(r)

	campaigns := supabase.GetCampaignByOrgId(orgId)

	var vapiCampaigns []api.Campaign
	for _, campaign := range campaigns {
		campaign, err := VapiClient.Campaigns.CampaignControllerFindOne(context.Background(), campaign.VapiCampaignId)
		if err != nil {
			http.Error(w, "Failed to get campaign", http.StatusInternalServerError)
			return
		}

		vapiCampaigns = append(vapiCampaigns, *campaign)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vapiCampaigns)
}
