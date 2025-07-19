package supabase

import (
	"log"
	"os"
	supabaseTypes "sarah/types/supabase"
)

type Campaign = supabaseTypes.Campaign

func GetAllCampaigns() []supabaseTypes.Campaign {
	var campaigns []supabaseTypes.Campaign

	_, err := SupabaseClient.From(os.Getenv("SUPABASE_TABLE_CAMPAIGNS")).Select("*", "", false).ExecuteTo(&campaigns)
	if err != nil {
		log.Fatal(err)
	}

	return campaigns
}

func GetCampaignByOrgId(orgId string) []supabaseTypes.Campaign {
	var campaigns []supabaseTypes.Campaign

	_, err := SupabaseClient.From(os.Getenv("SUPABASE_TABLE_CAMPAIGNS")).Select("*", "exact", false).Eq("org_id", orgId).ExecuteTo(&campaigns)
	if err != nil {
		log.Fatal(err)
	}

	return campaigns
}
