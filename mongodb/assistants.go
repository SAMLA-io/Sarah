package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetOrganizationAssistants retrieves all assistants for a specific organization from the database.
// This function queries the assistants collection in the organization's database
// and returns all assistant documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve assistants for
//
// Returns:
//   - []mongodb.Assistant: Array of assistants for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_ASSISTANTS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
//
// Example Usage:
//
//	assistants := GetOrganizationAssistants("org_1234567890abcdef")
//	// Returns all assistants for the specified organization
func GetOrganizationAssistants(orgId string) []mongodb.Assistant {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_ASSISTANTS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var assistants []mongodb.Assistant
	if err := cursor.All(context.Background(), &assistants); err != nil {
		log.Fatal(err)
	}

	return assistants
}
