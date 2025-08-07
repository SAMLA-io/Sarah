// Package mongodb contains data structures and types used for MongoDB operations
// in the Sarah Campaign Management API. This package defines the core data models
// for campaigns, customers, and scheduling.
package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Campaign represents a complete campaign configuration in the database.
// A campaign defines automated calling strategies with specific scheduling,
// customer lists, and integration with VapiAI assistants.
type Campaign struct {
	// Id is the unique MongoDB ObjectID for this campaign
	Id bson.ObjectID `json:"id" bson:"_id,omitempty"`

	// Name is the human-readable identifier for the campaign
	Name string `json:"name" bson:"name"`

	// AssistantId is the VapiAI assistant ID that will handle the calls
	AssistantId string `json:"assistant_id" bson:"assistant_id"`

	// PhoneNumberId is the VapiAI phone number ID to use for outbound calls
	PhoneNumberId string `json:"phone_number_id" bson:"phone_number_id"`

	// SchedulePlan defines when and how often the campaign should run
	SchedulePlan *SchedulePlan `json:"schedule_plan" bson:"schedule_plan"`

	// DynamicCustomers indicates if the campaign should use dynamic customers
	DynamicCustomers bool `json:"dynamic_customers" bson:"dynamic_customers"`

	// Customers is the list of customers to contact in this campaign
	Customers []Customer `json:"customers" bson:"customers"`

	// Type determines the recurrence pattern of the campaign
	Type CampaignType `json:"type" bson:"type"`

	// Status indicates the current state of the campaign
	Status CampaignStatus `json:"status" bson:"status"`

	// StartDate is when the campaign should begin execution
	StartDate *time.Time `json:"start_date" bson:"start_date"`

	// EndDate is when the campaign should stop execution
	EndDate *time.Time `json:"end_date" bson:"end_date"`

	// TimeZone is the timezone for all date/time calculations (e.g., "America/New_York")
	TimeZone string `json:"timezone" bson:"timezone"`
}

// SchedulePlan defines the scheduling strategy for a campaign.
// This structure allows for flexible scheduling with multiple recurrence patterns.
type SchedulePlan struct {
	// BeforeDay specifies how many days before a customer's relevant date to make the call
	// For example, if BeforeDay=3 and a customer has an expiry on day 15, calls will be made on day 12
	BeforeDay int `json:"before_day" bson:"before_day"`

	// AfterDay specifies how many days after a customer's relevant date to make the call
	// For example, if AfterDay=1 and a customer has an expiry on day 15, calls will be made on day 16
	AfterDay int `json:"after_day" bson:"after_day"`
}

// Customer represents an individual customer in a campaign with their contact information
// and scheduling preferences. Each customer can have different scheduling rules.
type Customer struct {
	// PhoneNumber is the customer's contact phone number in E.164 format (e.g., "+1234567890")
	PhoneNumber string `json:"phone_number" bson:"phone_number"`

	// DayNumber is the day of the month when this customer's calls should be scheduled
	// This is typically used for monthly or yearly campaigns
	DayNumber int `json:"day_number" bson:"day_number"`

	// MonthNumber is the month when this customer's calls should be scheduled (1-12)
	// This is typically used for yearly campaigns
	MonthNumber int `json:"month_number" bson:"month_number"`

	// YearNumber is the year when this customer's calls should be scheduled
	// This is typically used for yearly campaigns
	YearNumber int `json:"year_number" bson:"year_number"`
}

// CampaignType defines the different types of campaign recurrence patterns.
type CampaignType string

const (
	// RECURRENT_WEEKLY runs the campaign on a weekly basis
	// Example: Every Monday, Wednesday, Friday
	RECURRENT_WEEKLY CampaignType = "recurrent_weekly"

	// RECURRENT_MONTHLY runs the campaign on a monthly basis
	// Example: Every 15th of the month
	RECURRENT_MONTHLY CampaignType = "recurrent_monthly"

	// RECURRENT_YEARLY runs the campaign on a yearly basis
	// Example: Every March 15th
	RECURRENT_YEARLY CampaignType = "recurrent_yearly"

	// ONE_TIME runs the campaign only once on a specific date
	// Example: A single reminder call on a specific date
	ONE_TIME CampaignType = "one_time"
)

// CampaignStatus defines the possible states of a campaign.
type CampaignStatus string

const (
	// STATUS_ACTIVE indicates the campaign is currently running and executing calls
	STATUS_ACTIVE CampaignStatus = "active"

	// STATUS_PAUSED indicates the campaign is temporarily stopped but can be resumed
	STATUS_PAUSED CampaignStatus = "paused"

	// STATUS_COMPLETED indicates the campaign has finished all scheduled calls
	STATUS_COMPLETED CampaignStatus = "completed"

	// STATUS_CANCELLED indicates the campaign has been permanently stopped
	STATUS_CANCELLED CampaignStatus = "cancelled"
)
