package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

// Assistant represents a VapiAI assistant configuration in the database.
// This structure stores information about AI assistants that can be used
// for automated calling campaigns.
type Assistant struct {
	// Id is the unique MongoDB ObjectID for this assistant
	Id bson.ObjectID `json:"id" bson:"_id,omitempty"`

	// Name is the human-readable name for the assistant (e.g., "Insurance Reminder Assistant")
	Name string `json:"name" bson:"name"`

	// VapiAssistantId is the unique identifier for the assistant in VapiAI
	// This ID is used when making API calls to VapiAI
	VapiAssistantId string `json:"vapi_assistant_id" bson:"vapi_assistant_id"`

	// Type describes the category or purpose of the assistant
	// Examples: "insurance", "appointment", "reminder", "support"
	Type string `json:"type" bson:"type"`
}
