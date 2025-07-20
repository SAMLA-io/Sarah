package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sarah/mongodb"

	mongodbTypes "sarah/types/mongodb"

	api "github.com/VapiAI/server-sdk-go"
)

type Campaign = mongodbTypes.Campaign

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaignCreateDto := ExtractCampaignCreateDto(r)

	customers := []*api.CreateCustomerDto{}
	for _, customer := range campaignCreateDto.Customers {
		customers = append(customers, &api.CreateCustomerDto{
			Number: api.String(customer.Number),
		})
	}

	campaign, err := VapiClient.Campaigns.CampaignControllerCreate(context.Background(), &api.CreateCampaignDto{
		Name:          campaignCreateDto.Name,
		AssistantId:   &campaignCreateDto.AssistantId,
		PhoneNumberId: campaignCreateDto.PhoneNumberId,
		Customers:     customers,
	})

	orgId := ExtractOrgId(r)

	mongodb.CreateCampaign(orgId, Campaign{
		VapiId: campaign.Id,
		Type:   string(campaign.Status),
	})

	if err != nil {
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
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

	campaigns := mongodb.GetCampaignByOrgId(orgId)

	var vapiCampaigns []api.Campaign
	for _, campaign := range campaigns {
		campaign, err := VapiClient.Campaigns.CampaignControllerFindOne(context.Background(), campaign.VapiId)
		if err != nil {
			http.Error(w, "Failed to get campaign", http.StatusInternalServerError)
			return
		}

		vapiCampaigns = append(vapiCampaigns, *campaign)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vapiCampaigns)
}
