package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetContactByOrgId retrieves all contacts for a specific organization from the database.
// This function queries the contacts collection in the organization's database
// and returns all contact documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve contacts for
//
// Returns:
//   - []mongodb.Contact: Array of contacts for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_CONTACTS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
//
// Example Usage:
//
//	contacts := GetContactByOrgId("org_1234567890abcdef")
//	// Returns all contacts for the specified organization
func GetContactByOrgId(orgId string) []mongodb.Contact {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var contacts []mongodb.Contact
	if err := cursor.All(context.Background(), &contacts); err != nil {
		log.Fatal(err)
	}

	return contacts
}

// CreateContact creates a new contact in the database.
// This function inserts a new contact document into the contacts collection
// for a specific organization.
//
// Parameters:
//   - orgId: The organization ID to create the contact for
//   - contact: The contact to create
func CreateContact(orgId string, contact mongodb.Contact) *mongo.InsertOneResult {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	result, err := coll.InsertOne(context.Background(), contact)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
