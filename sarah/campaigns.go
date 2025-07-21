package sarah

import (
	"fmt"
	"log"
	"time"

	clerk "sarah/clerk"
	"sarah/mongodb"

	mongodbTypes "sarah/types/mongodb"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}
}

func CreateCampaign(campaignCreateDto mongodbTypes.Campaign, orgId string) *mongo.InsertOneResult {
	campaign := mongodb.CreateCampaign(orgId, mongodbTypes.Campaign{
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

	if campaign == nil {
		return nil
	}

	return campaign

}

// iterate voer all orgs in clerk
// for each org, get the campaings from mongodb
// for each campaign, check if it is time to send the call
// if it is, create the one-time campaign in Vapi with the phone nombers
// from the mongodb campaign

type CampaignScheduler struct {
}

func (c *CampaignScheduler) Start() {
	go func() {
		//c.run()
	}()
}

func (c *CampaignScheduler) Stop() {

}

func Run() {
	start := time.Now()
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	allOrgIDs, err := clerk.GetAllOrganizations()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Retrieved %d organizations:\n", len(allOrgIDs))

	for _, id := range allOrgIDs {
		campaigns := mongodb.GetCampaignByOrgId(id)
		fmt.Printf("Retrieved %d campaigns for organization %s:\n", len(campaigns), id)

		for _, campaign := range campaigns {
			fmt.Printf("Campaign: %s\n", campaign.Name)
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
	// for {
	// 	// Output the IDs

	// 	time.Sleep(1 * time.Second)
	// }
}

// check if the campaign is time to send the call
// if it is, create the one-time campaign in Vapi with the phone nombers
// from the mongodb campaign
func CheckCampaign(campaign mongodbTypes.Campaign) {
	campaignType := campaign.Type

	switch campaignType {
	case mongodbTypes.RECURRENT_WEEKLY:
	case mongodbTypes.RECURRENT_MONTHLY:
	case mongodbTypes.RECURRENT_YEARLY:
	case mongodbTypes.ONE_TIME:
	default:
		fmt.Printf("Campaign type %s not supported\n", campaignType)
		return
	}
}

func CheckRecurrentWeeklyCampaign(campaign mongodbTypes.Campaign) {

}

func CheckRecurrentMonthlyCampaign(campaign mongodbTypes.Campaign) {

}

func CheckRecurrentYearlyCampaign(campaign mongodbTypes.Campaign) {

}

func CheckOneTimeCampaign(campaign mongodbTypes.Campaign) {
}
