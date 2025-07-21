package api

import (
	"encoding/json"
	"net/http"
	"sarah/mongodb"

	mongodbTypes "sarah/types/mongodb"
)

type Campaign = mongodbTypes.Campaign

// CreateCampaign handles POST requests to create a new campaign.
// This endpoint accepts a campaign creation request and stores it in the database.
//
// HTTP Method: POST
// Endpoint: /campaigns/create
//
// Request Body:
//
//	{
//	  "campaignCreateRequest": {
//	    "name": "Weekly Insurance Reminders",
//	    "assistant_id": "asst_1234567890abcdef",
//	    "phone_number_id": "phone_0987654321fedcba",
//	    "schedule_plan": {
//	      "before_day": 3,
//	      "after_day": 0,
//	      "week_days": [1, 3, 5],
//	      "month_days": [],
//	      "year_months": []
//	    },
//	    "customers": [
//	      {
//	        "phone_number": "+1234567890",
//	        "day_number": 15,
//	        "month_number": 3,
//	        "week_day": 1,
//	        "custom_date": null,
//	        "expiry_date": "2024-12-31T23:59:59Z"
//	      }
//	    ],
//	    "type": "recurrent_weekly",
//	    "status": "active",
//	    "start_date": "2024-01-01T00:00:00Z",
//	    "end_date": "2024-12-31T23:59:59Z",
//	    "timezone": "America/New_York"
//	  }
//	}
//
// Query Parameters:
//   - orgId: The organization ID (required)
//
// Response:
//   - 200 OK: Campaign created successfully, returns the created campaign
//   - 405 Method Not Allowed: If not using POST method
//   - 500 Internal Server Error: If database operation fails
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

// GetCampaignViaOrgID handles GET requests to retrieve all campaigns for an organization.
// This endpoint returns all campaigns associated with the specified organization ID.
//
// HTTP Method: GET
// Endpoint: /campaigns/org
//
// Query Parameters:
//   - orgId: The organization ID (required)
//
// Response:
//   - 200 OK: Returns an array of campaigns for the organization
//   - 405 Method Not Allowed: If not using GET method
//   - 500 Internal Server Error: If database operation fails
//
// Example Response:
//
//	[
//	  {
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
//	]
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
