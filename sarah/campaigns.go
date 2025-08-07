package sarah

import (
	"fmt"
	"log"
	"os"
	"time"

	clerk "sarah/clerk"
	"sarah/mongodb"

	mongodbTypes "sarah/types/mongodb"

	api "github.com/VapiAI/server-sdk-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}
}

/* API Methods */

func CreateCampaign(campaignCreateDto mongodbTypes.Campaign, orgId string) (*mongo.InsertOneResult, error) {
	campaign, err := mongodb.CreateCampaign(orgId, mongodbTypes.Campaign{
		Name:             campaignCreateDto.Name,
		AssistantId:      campaignCreateDto.AssistantId,
		PhoneNumberId:    campaignCreateDto.PhoneNumberId,
		SchedulePlan:     campaignCreateDto.SchedulePlan,
		Customers:        campaignCreateDto.Customers,
		Type:             campaignCreateDto.Type,
		Status:           campaignCreateDto.Status,
		StartDate:        campaignCreateDto.StartDate,
		EndDate:          campaignCreateDto.EndDate,
		TimeZone:         campaignCreateDto.TimeZone,
		DynamicCustomers: campaignCreateDto.DynamicCustomers,
	})

	if campaign == nil {
		return nil, err
	}

	return campaign, nil

}

// iterate voer all orgs in clerk
// for each org, get the campaings from mongodb
// for each campaign, check if it is time to send the call
// if it is, create the one-time campaign in Vapi with the phone nombers
// from the mongodb campaign

type CampaignScheduler struct{}

func (c *CampaignScheduler) Start() {
	go func() {
		c.run()
	}()
}

func (c *CampaignScheduler) Stop() {
	os.Exit(0)
}

func (c *CampaignScheduler) run() {
	for {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found, using system environment variables")
		}

		allOrgIDs, err := clerk.GetAllOrganizations()
		if err != nil {
			panic(err)
		}

		log.Printf("[CampaignScheduler] Retrieved %d organizations", len(allOrgIDs))

		for _, id := range allOrgIDs {
			campaigns, err := mongodb.GetCampaignByOrgId(id)
			if err != nil {
				log.Printf("Error getting campaigns for organization %s: %v", id, err)
				continue
			}
			log.Printf("[CampaignScheduler] Retrieved %d campaigns for organization %s", len(campaigns), id)

			for _, campaign := range campaigns {
				if campaign.Status == mongodbTypes.STATUS_ACTIVE {
					log.Printf("[CampaignScheduler] Campaign: %s", campaign.Name)
					err := CheckCampaign(id, campaign)
					if err != nil {
						log.Printf("[CampaignScheduler] Error checking campaign: %v", err)
					}
				}
			}

		}

		log.Printf("[CampaignScheduler] --------------------------------")
		time.Sleep(1 * time.Minute)
	}
}

// check if the campaign is time to send the call
// if it is, create the one-time campaign in Vapi with the phone nombers
// from the mongodb campaign
func CheckCampaign(orgId string, campaign mongodbTypes.Campaign) error {
	campaignType := campaign.Type

	switch campaignType {
	case mongodbTypes.RECURRENT_WEEKLY:
		return CheckRecurrentWeeklyCampaign(orgId, campaign)
	case mongodbTypes.RECURRENT_MONTHLY:
		return CheckRecurrentMonthlyCampaign(orgId, campaign)
	case mongodbTypes.RECURRENT_YEARLY:
		return CheckRecurrentYearlyCampaign(orgId, campaign)
	case mongodbTypes.ONE_TIME:
		return CheckOneTimeCampaign(orgId, campaign)
	default:
		log.Printf("[CampaignScheduler] Campaign type %s not supported", campaignType)
		return fmt.Errorf("campaign type %s not supported", campaignType)
	}
}

// Helper function to get timezone location
func getTimezoneLocation(timezone string) *time.Location {
	if timezone == "" {
		return time.UTC
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Warning: Invalid timezone %s, using UTC: %v", timezone, err)
		return time.UTC
	}
	return loc
}

// Helper function to check if a customer should be called based on BeforeDay/AfterDay logic
func shouldCallCustomer(customer mongodbTypes.Customer, now time.Time, schedulePlan *mongodbTypes.SchedulePlan, campaignType mongodbTypes.CampaignType, timezone string) bool {
	if schedulePlan == nil {
		return false
	}

	// Get customer's target date based on campaign type
	var customerTargetDate time.Time
	loc := getTimezoneLocation(timezone)

	switch campaignType {
	case mongodbTypes.RECURRENT_WEEKLY:
		// For weekly campaigns, use the customer's DayNumber as day of week (0-6)
		if customer.DayNumber == -1 {
			return false
		}
		// Calculate the next occurrence of this weekday
		daysUntilTarget := (customer.DayNumber - int(now.Weekday()) + 7) % 7
		if daysUntilTarget == 0 {
			customerTargetDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		} else {
			customerTargetDate = now.AddDate(0, 0, daysUntilTarget)
		}

	case mongodbTypes.RECURRENT_MONTHLY:
		// For monthly campaigns, use the customer's DayNumber as day of month (1-31)
		if customer.DayNumber == -1 {
			return false
		}
		// Calculate the next occurrence of this day in the month
		year, month, _ := now.Date()
		customerTargetDate = time.Date(year, month, customer.DayNumber, 0, 0, 0, 0, loc)

		// Compare only date components (ignore time)
		if customerTargetDate.Year() < now.Year() ||
			(customerTargetDate.Year() == now.Year() && customerTargetDate.Month() < now.Month()) ||
			(customerTargetDate.Year() == now.Year() && customerTargetDate.Month() == now.Month() && customerTargetDate.Day() < now.Day()) {

			// Move to next month
			customerTargetDate = customerTargetDate.AddDate(0, 1, 0)
		}

	case mongodbTypes.RECURRENT_YEARLY:
		// For yearly campaigns, use MonthNumber and DayNumber
		if customer.MonthNumber == -1 || customer.DayNumber == -1 {
			return false
		}
		year, _, _ := now.Date()
		customerTargetDate = time.Date(year, time.Month(customer.MonthNumber), customer.DayNumber, 0, 0, 0, 0, loc)
		// Compare only date components (ignore time)
		if customerTargetDate.Year() < now.Year() ||
			(customerTargetDate.Year() == now.Year() && customerTargetDate.Month() < now.Month()) ||
			(customerTargetDate.Year() == now.Year() && customerTargetDate.Month() == now.Month() && customerTargetDate.Day() < now.Day()) {
			// Move to next year
			customerTargetDate = customerTargetDate.AddDate(1, 0, 0)
		}

	case mongodbTypes.ONE_TIME:
		// For one-time campaigns, use YearNumber, MonthNumber, and DayNumber
		if customer.YearNumber == -1 || customer.MonthNumber == -1 || customer.DayNumber == -1 {
			return false
		}
		customerTargetDate = time.Date(customer.YearNumber, time.Month(customer.MonthNumber), customer.DayNumber, 0, 0, 0, 0, loc)
	default:
		return false
	}

	// Check BeforeDay logic
	if schedulePlan.BeforeDay != -1 {
		beforeDate := customerTargetDate.AddDate(0, 0, -schedulePlan.BeforeDay)
		if now.Year() == beforeDate.Year() && now.Month() == beforeDate.Month() && now.Day() == beforeDate.Day() {
			return true
		}
	}

	// Check AfterDay logic
	if schedulePlan.AfterDay != -1 {
		afterDate := customerTargetDate.AddDate(0, 0, schedulePlan.AfterDay)
		if now.Year() == afterDate.Year() && now.Month() == afterDate.Month() && now.Day() == afterDate.Day() {
			return true
		}
	}

	// If neither BeforeDay nor AfterDay is specified, call on the exact target date
	if schedulePlan.BeforeDay == -1 && schedulePlan.AfterDay == -1 {
		if now.Year() == customerTargetDate.Year() && now.Month() == customerTargetDate.Month() && now.Day() == customerTargetDate.Day() {
			return true
		}
	}

	return false
}

func CheckRecurrentWeeklyCampaign(orgId string, campaign mongodbTypes.Campaign) error {
	log.Printf("[CampaignScheduler] Checking recurrent weekly campaign: %s", campaign.Name)

	customers, err := getEligibleCustomers(orgId, campaign)
	if err != nil {
		log.Printf("[CampaignScheduler] Error getting eligible customers: %v", err)
		return err
	}

	if len(customers) == 0 {
		return nil
	}

	resp, err := executeCampaign(campaign.AssistantId, campaign.PhoneNumberId, customers)

	if err != nil {
		log.Printf("Error creating campaign: %v", err)
		return err
	}

	if resp == nil {
		return fmt.Errorf("campaign not created")
	}

	return nil
}

func CheckRecurrentMonthlyCampaign(orgId string, campaign mongodbTypes.Campaign) error {
	log.Printf("[CampaignScheduler] Checking recurrent monthly campaign: %s", campaign.Name)

	customers, err := getEligibleCustomers(orgId, campaign)
	if err != nil {
		log.Printf("[CampaignScheduler] Error getting eligible customers: %v", err)
		return err
	}

	if len(customers) == 0 {
		return nil
	}

	resp, err := executeCampaign(campaign.AssistantId, campaign.PhoneNumberId, customers)

	if err != nil {
		log.Printf("Error creating campaign: %v", err)
		return err
	}

	if resp == nil {
		return fmt.Errorf("campaign not created: %v", err)
	}

	return nil
}

func CheckRecurrentYearlyCampaign(orgId string, campaign mongodbTypes.Campaign) error {
	log.Printf("[CampaignScheduler] Checking recurrent yearly campaign: %s", campaign.Name)

	customers, err := getEligibleCustomers(orgId, campaign)
	if err != nil {
		log.Printf("[CampaignScheduler] Error getting eligible customers: %v", err)
		return err
	}

	if len(customers) == 0 {
		return nil
	}

	resp, err := executeCampaign(campaign.AssistantId, campaign.PhoneNumberId, customers)

	if err != nil {
		log.Printf("Error creating campaign: %v", err)
		return err
	}

	if resp == nil {
		return fmt.Errorf("campaign not created: %v", err)
	}

	return nil
}

func CheckOneTimeCampaign(orgId string, campaign mongodbTypes.Campaign) error {
	log.Printf("[CampaignScheduler] Checking one-time campaign: %s", campaign.Name)

	customers, err := getEligibleCustomers(orgId, campaign)
	if err != nil {
		log.Printf("[CampaignScheduler] Error getting eligible customers: %v", err)
		return err
	}

	if len(customers) == 0 {
		return nil
	}

	resp, err := executeCampaign(campaign.AssistantId, campaign.PhoneNumberId, customers)

	if err != nil {
		log.Printf("[CampaignScheduler] Error creating campaign: %v", err)
		return err
	}

	if resp == nil {
		log.Printf("[CampaignScheduler] Campaign not created, executeCampaign returned nil response")
		return fmt.Errorf("campaign not created, executeCampaign returned nil response")
	}

	campaign.Status = mongodbTypes.STATUS_COMPLETED

	res, err := mongodb.UpdateCampaign(orgId, campaign)

	if err != nil {
		log.Printf("[CampaignScheduler] Error updating campaign: %v", err)
		return err
	}

	if res.MatchedCount == 0 {
		log.Printf("[CampaignScheduler] Campaign not updated, updateCampaign returned nil response")
		return fmt.Errorf("campaign not updated, updateCampaign returned matched count 0")
	}

	log.Printf("[CampaignScheduler] Campaign executed: %+v", res)
	return nil
}

// Creates an immediate campaign in Vapi
func executeCampaign(assistantId string, phoneNumberId string, customers []mongodbTypes.Customer) (*api.CallsCreateResponse, error) {
	resp, err := CreateCall(assistantId, phoneNumberId, customers)

	if err != nil {
		log.Printf("[CampaignScheduler] Error creating call: %v", err)
		return nil, err
	}

	log.Printf("[CampaignScheduler] Campaign created: %+v", resp)
	return resp, nil
}

func getDynamicCustomers(orgId string) ([]mongodbTypes.Customer, error) {
	contacts, err := mongodb.GetContactByOrgId(orgId)
	if err != nil {
		log.Printf("[CampaignScheduler] Error getting contacts: %v", err)
		return nil, err
	}

	customers := []mongodbTypes.Customer{}

	for _, contact := range contacts {
		customers = append(customers, contact.Customer)
	}

	return customers, nil
}

func getEligibleCustomers(orgId string, campaign mongodbTypes.Campaign) ([]mongodbTypes.Customer, error) {
	log.Printf("[CampaignScheduler] Getting eligible customers for campaign: %s", campaign.Name)

	now := time.Now()
	loc := getTimezoneLocation(campaign.TimeZone)
	now = now.In(loc)

	customers := []mongodbTypes.Customer{}

	if campaign.DynamicCustomers {
		allCustomers, err := getDynamicCustomers(orgId)
		if err != nil {
			log.Printf("[CampaignScheduler] Error getting dynamic customers: %v", err)
			return nil, err
		}

		for _, customer := range allCustomers {
			if shouldCallCustomer(customer, now, campaign.SchedulePlan, campaign.Type, campaign.TimeZone) {
				customers = append(customers, customer)
			}
		}
	} else {
		for _, customer := range campaign.Customers {
			if shouldCallCustomer(customer, now, campaign.SchedulePlan, campaign.Type, campaign.TimeZone) {
				customers = append(customers, customer)
			}
		}
	}

	return customers, nil
}
