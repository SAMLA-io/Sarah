package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Campaign struct {
	Id            bson.ObjectID  `json:"id" bson:"_id"`
	Name          string         `json:"name" bson:"name"`
	AssistantId   string         `json:"assistant_id" bson:"assistant_id"`
	PhoneNumberId string         `json:"phone_number_id" bson:"phone_number_id"`
	SchedulePlan  *SchedulePlan  `json:"schedule_plan" bson:"schedule_plan"`
	Customers     []Customer     `json:"customers" bson:"customers"`
	Type          CampaignType   `json:"type" bson:"type"`
	Status        CampaignStatus `json:"status" bson:"status"`
	StartDate     *time.Time     `json:"start_date" bson:"start_date"`
	EndDate       *time.Time     `json:"end_date" bson:"end_date"`
	TimeZone      string         `json:"timezone" bson:"timezone"`
}

type CampaignCreateDto struct {
	Name          string         `json:"name" bson:"name"`
	AssistantId   string         `json:"assistant_id" bson:"assistant_id"`
	PhoneNumberId string         `json:"phone_number_id" bson:"phone_number_id"`
	SchedulePlan  *SchedulePlan  `json:"schedule_plan" bson:"schedule_plan"`
	Customers     []Customer     `json:"customers" bson:"customers"`
	Type          CampaignType   `json:"type" bson:"type"`
	Status        CampaignStatus `json:"status" bson:"status"`
	StartDate     *time.Time     `json:"start_date" bson:"start_date"`
	EndDate       *time.Time     `json:"end_date" bson:"end_date"`
	TimeZone      string         `json:"timezone" bson:"timezone"`
}

type SchedulePlan struct {
	BeforeDay  int   `json:"before_day" bson:"before_day"`
	AfterDay   int   `json:"after_day" bson:"after_day"`
	WeekDays   []int `json:"week_days" bson:"week_days"`     // 0=Sunday, 1=Monday, etc.
	MonthDays  []int `json:"month_days" bson:"month_days"`   // Days of month (1-31)
	YearMonths []int `json:"year_months" bson:"year_months"` // Months (1-12)
}

type Customer struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	DayNumber   int    `json:"day_number" bson:"day_number"`
	MonthNumber int    `json:"month_number" bson:"month_number"`
	WeekDay     int    `json:"week_day" bson:"week_day"`
	// Additional fields for more flexible scheduling
	CustomDate *time.Time `json:"custom_date" bson:"custom_date"` // For one-time specific dates
	ExpiryDate *time.Time `json:"expiry_date" bson:"expiry_date"` // For insurance/annual renewals
}

type CampaignType string

const (
	RECURRENT_WEEKLY  CampaignType = "recurrent_weekly"
	RECURRENT_MONTHLY CampaignType = "recurrent_monthly"
	RECURRENT_YEARLY  CampaignType = "recurrent_yearly"
	ONE_TIME          CampaignType = "one_time"
)

type CampaignStatus string

const (
	STATUS_ACTIVE    CampaignStatus = "active"
	STATUS_PAUSED    CampaignStatus = "paused"
	STATUS_COMPLETED CampaignStatus = "completed"
	STATUS_CANCELLED CampaignStatus = "cancelled"
)
