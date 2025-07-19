package supabase

import (
	"log"
	"os"
	supabaseTypes "sarah/types/supabase"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Assistant = supabaseTypes.Assistant

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

func GetAllAssistants() []Assistant {
	var assistants []Assistant

	_, err := SupabaseClient.From(os.Getenv("SUPABASE_TABLE_ASSISTANTS")).Select("*", "", false).ExecuteTo(&assistants)
	if err != nil {
		log.Fatal(err)
	}

	return assistants
}

func GetAssistantByOrgId(orgId string) []Assistant {
	var assistants []Assistant

	_, err := SupabaseClient.From(os.Getenv("SUPABASE_TABLE_ASSISTANTS")).Select("*", "exact", false).Eq("org_id", orgId).ExecuteTo(&assistants)
	if err != nil {
		log.Fatal(err)
	}

	return assistants
}
