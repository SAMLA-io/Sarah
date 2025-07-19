package supabase

type Assistant struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	VapiAssistantId string `json:"vapi_assistant_id"`
	OrgId           string `json:"org_id"`
	Type            string `json:"type"`
}
