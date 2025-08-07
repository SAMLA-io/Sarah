package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Contact represents a customer contact in the database.
// This structure stores comprehensive information about customers
// that can be used in campaigns and call management.
type Contact struct {
	// Id is the unique MongoDB ObjectID for this contact
	Id bson.ObjectID `json:"id" bson:"_id,omitempty"`

	// Name is the full name of the contact (e.g., "John Doe")
	Name string `json:"name" bson:"name"`

	// Email is the contact's email address (e.g., "john.doe@example.com")
	Email string `json:"email" bson:"email"`

	// PhoneNumber is the contact's phone number in E.164 format (e.g., "+1234567890")
	PhoneNumber string `json:"phone_number" bson:"phone_number"`

	// Company is the name of the company the contact works for
	Company string `json:"company" bson:"company"`

	// Position is the contact's job title or position within their company
	Position string `json:"position" bson:"position"`

	// Address is the contact's physical address
	Address string `json:"address" bson:"address"`

	// Customer is the customer object that this contact belongs to
	Customer Customer `json:"customer" bson:"customer"`

	// Metadata is a flexible field for storing additional contact information
	// This can include custom fields, preferences, or any other relevant data
	Metadata map[string]interface{} `json:"metadata" bson:"metadata"`
}
