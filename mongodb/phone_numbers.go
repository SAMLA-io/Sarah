package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetPhoneNumberByOrgId retrieves all phone numbers for a specific organization from the database.
// This function queries the phone numbers collection in the organization's database
// and returns all phone number documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve phone numbers for
//
// Returns:
//   - []mongodb.PhoneNumber: Array of phone numbers for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_PHONE_NUMBERS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
//
// Example Usage:
//
//	phoneNumbers := GetPhoneNumberByOrgId("org_1234567890abcdef")
//	// Returns all phone numbers for the specified organization
func GetPhoneNumberByOrgId(orgId string) []mongodb.PhoneNumber {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_PHONE_NUMBERS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var phoneNumbers []mongodb.PhoneNumber
	if err := cursor.All(context.Background(), &phoneNumbers); err != nil {
		log.Fatal(err)
	}

	return phoneNumbers
}
