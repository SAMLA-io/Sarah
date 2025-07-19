package supabase

import (
	"log"
	"os"
	supabaseTypes "sarah/types/supabase"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Campaign = supabaseTypes.Campaign

var SupabaseClient *supabase.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	var err error
	SupabaseClient, err = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

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
