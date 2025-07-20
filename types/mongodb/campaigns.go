package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Campaign struct {
	Id     bson.ObjectID `json:"id" bson:"_id"`
	VapiId string        `json:"vapi_id" bson:"vapi_id"`
	Type   string        `json:"type" bson:"type"`
}

type CampaignCreateDto struct {
	Name          string        `json:"name" bson:"name"`
	AssistantId   string        `json:"assistant_id" bson:"assistant_id"`
	PhoneNumberId string        `json:"phone_number_id" bson:"phone_number_id"`
	SchedulePlan  *SchedulePlan `json:"schedule_plan" bson:"schedule_plan"`
	Customers     []Customer    `json:"customers" bson:"customers"`
}

type SchedulePlan struct {
	BeforeDay int `json:"before_day" bson:"before_day"`
	AfterDay  int `json:"after_day" bson:"after_day"`
}

type Customer struct {
	Number    string `json:"number" bson:"number"`
	DayNumber int    `json:"day_number" bson:"day_number"`
}
