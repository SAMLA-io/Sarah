package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetCampaignByOrgId retrieves all campaigns for a specific organization from the database.
// This function queries the campaigns collection in the organization's database
// and returns all campaign documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve campaigns for
//
// Returns:
//   - []mongodb.Campaign: Array of campaigns for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_CAMPAIGNS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
func GetCampaignByOrgId(orgId string) []mongodb.Campaign {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CAMPAIGNS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var campaigns []mongodb.Campaign
	if err := cursor.All(context.Background(), &campaigns); err != nil {
		log.Fatal(err)
	}

	return campaigns
}

// CreateCampaign creates a new campaign in the database for the specified organization.
// This function inserts a campaign document into the campaigns collection
// and returns the result of the insertion operation.
//
// Parameters:
//   - orgId: The organization ID to create the campaign for
//   - campaign: The campaign data to insert into the database
//
// Returns:
//   - *mongo.InsertOneResult: The result of the insertion operation
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_CAMPAIGNS environment variable
//   - Operation: Inserts a single campaign document
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
func CreateCampaign(orgId string, campaign mongodb.Campaign) *mongo.InsertOneResult {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CAMPAIGNS"))

	result, err := coll.InsertOne(context.Background(), campaign)
	if err != nil {
		return nil
	}

	return result
}
