package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

// PhoneNumber represents a VapiAI phone number configuration in the database.
// This structure stores information about phone numbers that can be used
// for outbound calls in campaigns.
type PhoneNumber struct {
	// Id is the unique MongoDB ObjectID for this phone number record
	Id bson.ObjectID `json:"id" bson:"_id"`

	// Name is the human-readable name for the phone number (e.g., "Main Office Line")
	Name string `json:"name" bson:"name"`

	// PhoneNumberId is the unique identifier for the phone number in VapiAI
	// This ID is used when making API calls to VapiAI for outbound calls
	PhoneNumberId string `json:"phone_number_id" bson:"phone_number_id"`

	// PhoneNumber is the actual phone number in E.164 format (e.g., "+1987654321")
	// This is the number that will be displayed to recipients during calls
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}
