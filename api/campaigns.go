package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"

	mongodbTypes "sarah/types/mongodb"
)

type Campaign = mongodbTypes.Campaign

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignCreateDto := ExtractCampaignCreateDto(r)
	orgId := ExtractOrgId(r)

	// Adds the campaign to the database
	campaign := mongodb.CreateCampaign(orgId, Campaign{
		Name:          campaignCreateDto.Name,
		AssistantId:   campaignCreateDto.AssistantId,
		PhoneNumberId: campaignCreateDto.PhoneNumberId,
		SchedulePlan:  campaignCreateDto.SchedulePlan,
		Customers:     campaignCreateDto.Customers,
		Type:          campaignCreateDto.Type,
		Status:        campaignCreateDto.Status,
		StartDate:     campaignCreateDto.StartDate,
		EndDate:       campaignCreateDto.EndDate,
		TimeZone:      campaignCreateDto.TimeZone,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
}

func GetCampaignViaOrgID(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgId := ExtractOrgId(r)

	campaigns := mongodb.GetCampaignByOrgId(orgId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaigns)
}
