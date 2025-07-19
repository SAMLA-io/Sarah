package supabase

type Campaign struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	VapiCampaignId string `json:"vapi_campaign_id"`
	OrgId          string `json:"org_id"`
	Status         string `json:"status"`
}
